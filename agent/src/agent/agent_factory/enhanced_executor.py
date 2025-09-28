"""
Enhanced AgentExecutor that properly handles iteration limits.
"""

from typing import Any, Dict, List, Optional, Sequence, Tuple, Union
from langchain.agents import AgentExecutor
from langchain.schema import AgentAction, AgentFinish
from src.utils.logger import logger

class EnhancedAgentExecutor(AgentExecutor):
    """
    Enhanced version of AgentExecutor that provides better handling of iteration limits.
    """
    
    def _return_stopped_response(
        self,
        early_stopping_method: str,
        intermediate_steps: List[Tuple[AgentAction, str]],
        **kwargs: Any,
    ) -> Dict[str, Any]:
        """Return response when agent has been stopped due to max iterations or time limit."""
        logger.info(f"代理达到迭代限制，使用增强的响应处理...")
        logger.info(f"中间步骤数量: {len(intermediate_steps)}")
        logger.info(f"输入内容: {kwargs.get('input', 'N/A')}")
        
        if not intermediate_steps:
            logger.warning("没有中间步骤，返回默认消息")
            return {"output": "Agent stopped due to iteration limit or time limit."}
        
        for i, (action, observation) in enumerate(intermediate_steps):
            logger.info(f"步骤 {i+1}: 动作={action.tool}, 输入={action.tool_input}, 观察={observation}")

        observations = [observation for _, observation in intermediate_steps]
        
        if observations:
            last_observation = observations[-1].strip()
            logger.info(f"使用最后一个观察结果生成最终答案: {last_observation}")
            
            user_input = kwargs.get("input", "").lower()
            if "pwd" in user_input or "目录" in user_input:
                return {
                    "output": f"当前目录是: {last_observation}"
                }
            elif "shell" in user_input or "命令" in user_input:
                return {
                    "output": f"命令执行结果: {last_observation}"
                }
            else:
                return {
                    "output": f"执行结果: {last_observation}"
                }

        logger.warning("没有有效的观察结果")
        return {"output": "Agent stopped due to iteration limit or time limit."}
