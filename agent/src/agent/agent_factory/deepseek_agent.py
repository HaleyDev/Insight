from __future__ import annotations

import json
import logging
import re
from functools import partial
from operator import itemgetter
from typing import Any, List, Sequence, Tuple, Union

from langchain.agents.agent import AgentExecutor, RunnableAgent
from .enhanced_executor import EnhancedAgentExecutor
from langchain.agents.structured_chat.output_parser import StructuredChatOutputParser
from langchain.prompts.chat import BaseChatPromptTemplate
from langchain.schema import (
    AgentAction,
    AgentFinish,
    AIMessage,
    HumanMessage,
    OutputParserException,
    SystemMessage,
)
from langchain.schema.language_model import BaseLanguageModel
from langchain.tools.base import BaseTool
from langchain_core.callbacks import Callbacks
from langchain_core.runnables import Runnable, RunnablePassthrough
from src.utils.logger import logger
from .prompt import get_prompt_template, get_react_prompt_template

# from chatchat.server.utils import get_prompt_template  # Not available in this project


# DeepSeek models work well with streaming for tool calls
# but we keep the non-streaming option for consistency
def _plan_without_stream(
    self: RunnableAgent,
    intermediate_steps: List[Tuple[AgentAction, str]],
    callbacks: Callbacks = None,
    **kwargs: Any,
) -> Union[AgentAction, AgentFinish]:
    inputs = {**kwargs, "intermediate_steps": intermediate_steps}
    return self.runnable.invoke(inputs, config={"callbacks": callbacks})


async def _aplan_without_stream(
    self: RunnableAgent,
    intermediate_steps: List[Tuple[AgentAction, str]],
    callbacks: Callbacks = None,
    **kwargs: Any,
) -> Union[AgentAction, AgentFinish]:
    inputs = {**kwargs, "intermediate_steps": intermediate_steps}
    return await self.runnable.ainvoke(inputs, config={"callbacks": callbacks})


class DeepSeekChatAgentPromptTemplate(BaseChatPromptTemplate):
    # The template to use
    template: str
    # The list of tools available
    tools: List[BaseTool]

    def format_messages(self, **kwargs) -> str:
        # Get the intermediate steps (AgentAction, Observation tuples)
        # Format them in a particular way
        intermediate_steps = kwargs.pop("intermediate_steps", [])
        thoughts = ""
        for action, observation in intermediate_steps:
            thoughts += action.log
            thoughts += f"\nObservation: {observation}\nThought: "
        # Set the agent_scratchpad variable to that value
        if thoughts:
            kwargs[
                "agent_scratchpad"
            ] = f"These were previous tasks you completed:\n{thoughts}\n\n"
        else:
            kwargs["agent_scratchpad"] = ""
        
        # Create a tools variable from the list of tools provided
        # DeepSeek models prefer more structured tool descriptions
        tools = []
        for t in self.tools:
            desc = re.sub(r"\n+", " ", t.description)
            # DeepSeek models work better with cleaner tool descriptions
            text = (
                f"{t.name}: {desc.strip()}"
                f" Parameters: {t.args}"
            )
            tools.append(text)
        kwargs["tools"] = "\n\n".join(tools)
        
        # Create a list of tool names for the tools provided
        kwargs["tool_names"] = ", ".join([tool.name for tool in self.tools])
        formatted = self.template.format(**kwargs)
        
        # Debug: log the formatted template
        logger.info(f"Formatted template:\n{formatted}")
        
        return [HumanMessage(content=formatted)]


def validate_json(json_data: str) -> bool:
    """Validate if a string is valid JSON"""
    try:
        json.loads(json_data)
        return True
    except (ValueError, json.JSONDecodeError):
        return False


class DeepSeekChatAgentOutputParserCustom(StructuredChatOutputParser):
    """Output parser with retries for the structured chat agent with custom DeepSeek prompt."""

    def parse(self, text: str) -> Union[AgentAction, AgentFinish]:
        # 记录原始输出内容，用于调试
        logger.info(f"解析DeepSeek输出:\n{text}")
        
        # 尝试匹配常见的Action/Action Input模式
        patterns = [
            r"\nAction:\s*(.+)\nAction\s+Input:\s*(.+)",  # 标准格式
            r"Action:\s*(.+)\nAction\s+Input:\s*(.+)",    # 无前导换行
            r"Action:\s*(.+)[\s\n]+Action\s+Input:\s*(.+)",  # 间隔多个空格或换行
            r"Action:\s*([^\n]+).*?(?:Input|参数):\s*(.+)",  # 中英文混合格式
        ]
        
        for pattern in patterns:
            if s := re.findall(pattern, text, flags=re.DOTALL):
                s = s[-1]  # 取最后一个匹配结果
                action_name = s[0].strip()
                json_string: str = s[1].strip()
                
                logger.info(f"匹配到Action: {action_name}")
                logger.info(f"匹配到Action Input: {json_string}")
                
                json_input = None
                try:
                    json_input = json.loads(json_string)
                    logger.info(f"成功解析JSON: {json_input}")
                except (ValueError, json.JSONDecodeError) as e:
                    logger.warning(f"解析JSON失败: {e}, 尝试修复")
                    # DeepSeek sometimes outputs with minor formatting issues
                    try:
                        # Try to fix common issues
                        fixed_json_string = json_string
                        
                        # Fix single quotes to double quotes
                        if "'" in fixed_json_string and '"' not in fixed_json_string:
                            fixed_json_string = fixed_json_string.replace("'", '"')
                        
                        # Ensure proper JSON ending
                        if not fixed_json_string.endswith('}') and not fixed_json_string.endswith(']'):
                            if '{' in fixed_json_string:
                                fixed_json_string += '}'
                            elif '[' in fixed_json_string:
                                fixed_json_string += ']'
                        
                        # Remove any trailing commas
                        fixed_json_string = re.sub(r',(\s*[}\]])', r'\1', fixed_json_string)
                        
                        # 处理可能的中文冒号
                        fixed_json_string = fixed_json_string.replace("：", ":")
                        
                        logger.info(f"修复后的JSON字符串: {fixed_json_string}")
                        
                        if validate_json(fixed_json_string):
                            json_input = json.loads(fixed_json_string)
                            logger.info(f"修复JSON成功: {json_input}")
                        else:
                            # If all fixes fail, create a simple string input
                            json_input = {"query": json_string}
                            logger.warning(f"无法修复JSON，使用简单输入: {json_input}")
                            
                    except Exception as fix_error:
                        logger.error(f"修复JSON失败: {fix_error}")
                        json_input = {"query": json_string}

                return AgentAction(tool=action_name, tool_input=json_input, log=text)
        
        # 检查是否为最终答案
        final_patterns = [
            r"\nFinal\s+Answer:\s*(.+)",  # 英文格式
            r"\n最终\s*答案:\s*(.+)",      # 中文格式
            r"Final\s+Answer:\s*(.+)",    # 无前导换行
            r"最终\s*答案:\s*(.+)"         # 中文无前导换行
        ]
        
        for pattern in final_patterns:
            if s := re.findall(pattern, text, flags=re.DOTALL):
                s = s[-1].strip()
                logger.info(f"匹配到最终答案: {s}")
                return AgentFinish({"output": s}, log=text)
        
        # 如果未找到结构化格式，作为最终答案处理
        logger.warning(f"未找到结构化格式，将整个输出作为最终答案处理")
        return AgentFinish({"output": text.strip()}, log=text)

    @property
    def _type(self) -> str:
        return "StructuredDeepSeekChatOutputParserCustom"


class DeepSeekChatAgentOutputParserLC(StructuredChatOutputParser):
    """Output parser with retries for the structured chat agent with standard LangChain prompt."""

    def parse(self, text: str) -> Union[AgentAction, AgentFinish]:
        if s := re.findall(r"\nAction:\s*```(.+)```", text, flags=re.DOTALL):
            try:
                action = json.loads(s[0])
                tool = action.get("action")
                if tool == "Final Answer":
                    return AgentFinish({"output": action.get("action_input", "")}, log=text)
                else:
                    return AgentAction(
                        tool=tool, tool_input=action.get("action_input", {}), log=text
                    )
            except (ValueError, json.JSONDecodeError) as e:
                logger.error(f"Failed to parse DeepSeek LangChain format JSON: {e}")
                raise OutputParserException(f"Could not parse LLM output: {text}")
        else:
            raise OutputParserException(f"Could not parse LLM output: {text}")

    @property
    def _type(self) -> str:
        return "StructuredDeepSeekChatOutputParserLC"


def create_structured_deepseek_chat_agent(
    llm: BaseLanguageModel,
    tools: Sequence[BaseTool],
    callbacks: Sequence[Callbacks],
    use_custom_prompt: bool = True,
    enable_streaming: bool = False,
) -> AgentExecutor:
    """
    Create a structured chat agent optimized for DeepSeek models.
    
    Args:
        llm: The DeepSeek language model
        tools: Sequence of tools available to the agent
        callbacks: Callbacks for the agent
        use_custom_prompt: Whether to use custom DeepSeek prompt format
        enable_streaming: Whether to enable streaming (DeepSeek handles this better than Qwen)
    
    Returns:
        AgentExecutor: The configured agent executor
    """
    logger.info(f"创建DeepSeek代理，tools数量: {len(tools)}")
    
    if use_custom_prompt:
        prompt = "deepseek"
        output_parser = DeepSeekChatAgentOutputParserCustom()
        logger.info("使用自定义DeepSeek提示词模板和解析器")
    else:
        prompt = "structured-chat-agent"
        output_parser = DeepSeekChatAgentOutputParserLC()
        logger.info("使用标准LangChain提示词模板和解析器")

    tools = [t.copy(update={"callbacks": callbacks}) for t in tools]
    logger.info(f"工具名称: {[t.name for t in tools]}")
    
    template = get_prompt_template("action_model", prompt)
    logger.info(f"使用模板类型: action_model/{prompt}")
    
    prompt = DeepSeekChatAgentPromptTemplate(
        input_variables=["input", "intermediate_steps"], 
        template=template, 
        tools=tools
    )

    # DeepSeek models use different stop tokens
    stop_tokens = ["<|end_of_text|>", "<|eot_id|>", "\nObservation:", "Observation:"]
    logger.info(f"设置停止词: {stop_tokens}")
    
    agent = (
        RunnablePassthrough.assign(agent_scratchpad=itemgetter("intermediate_steps"))
        | prompt
        | llm.bind(stop=stop_tokens)
        | output_parser
    )
    
    # 配置执行器参数
    executor_kwargs = {
        "agent": agent,
        "tools": tools,
        "callbacks": callbacks,
        "verbose": True,
        "max_iterations": 3,  # 限制最大迭代次数
        # 移除不支持的early_stopping_method参数
    }
    logger.info(f"创建执行器，参数: {executor_kwargs}")
    
    # 使用增强的执行器来处理迭代限制
    executor = EnhancedAgentExecutor(**executor_kwargs)
    
    # Apply streaming fix only if streaming is disabled
    if not enable_streaming:
        logger.info("应用非流式修复")
        executor.agent.__dict__["plan"] = partial(_plan_without_stream, executor.agent)
        executor.agent.__dict__["aplan"] = partial(_aplan_without_stream, executor.agent)

    return executor


def create_deepseek_react_agent(
    llm: BaseLanguageModel,
    tools: Sequence[BaseTool],
    callbacks: Sequence[Callbacks],
) -> AgentExecutor:
    """
    Create a ReAct style agent optimized for DeepSeek models.
    
    Args:
        llm: The DeepSeek language model
        tools: Sequence of tools available to the agent
        callbacks: Callbacks for the agent
    
    Returns:
        AgentExecutor: The configured ReAct agent executor
    """
    from langchain.agents import create_react_agent
    from langchain.prompts import PromptTemplate
    
    # Use the react prompt template from prompt.py
    react_prompt_template = get_react_prompt_template()
    react_prompt = PromptTemplate.from_template(react_prompt_template)
    
    agent = create_react_agent(llm, tools, react_prompt)
    executor = AgentExecutor(agent=agent, tools=tools, callbacks=callbacks, verbose=True)
    
    return executor