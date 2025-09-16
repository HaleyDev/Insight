from typing import List, Sequence
from venv import logger
from langchain_core.language_models import BaseLanguageModel
from langchain_core.callbacks import BaseCallbackHandler
from langchain_core.tools import BaseTool
from src.utils.logger import logger
from .qwen_agent import create_structured_qwen_chat_agent
from langchain.agents import AgentExecutor, create_structured_chat_agent

from src.agent.agent_factory.deepseek_agent import create_structured_deepseek_chat_agent

def agents_registry(
        llm: BaseLanguageModel,
        tools: Sequence[BaseTool] = [],
        callbacks: List[BaseCallbackHandler] = [],
        prompt: str = None,
        verbose: bool = False,
):
    logger.info(f"创建代理，模型: {llm.model_name}, 工具数量: {len(tools)}, verbose: {verbose}")

    if "deepseek" in llm.model_name.lower():
        logger.info("Using DeepSeek model for agent.")
        return create_structured_deepseek_chat_agent(
            llm=llm, 
            tools=tools, 
            callbacks=callbacks,
            use_custom_prompt=True,
            enable_streaming=False
        )
    elif "qwen" in llm.model_name.lower():
        logger.info("Using Qwen model for agent.")
        llm.stream = False  # qwen model does not support streaming with tool calls
        return create_structured_qwen_chat_agent(llm=llm, tools=tools, callbacks=callbacks)
    else:
        raise ValueError(f"Unsupported model for agent: {llm.model_name}")
