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
        
        # 如果没有任何中间步骤，返回默认消息
        if not intermediate_steps:
            return {"output": "Agent stopped due to iteration limit or time limit."}
        
        # 提取所有观察结果
        observations = [observation for _, observation in intermediate_steps]
        
        # 如果我们有观察结果，返回最后一个观察结果作为最终答案
        if observations:
            last_observation = observations[-1]
            logger.info(f"使用最后一个观察结果生成最终答案: {last_observation}")
            
            # 将最后的观察结果格式化为用户友好的答案
            if "pwd" in kwargs.get("input", "").lower():
                return {
                    "output": f"当前目录是: {last_observation}"
                }
            else:
                return {
                    "output": f"结果: {last_observation}"
                }
                
        # 如果没有有效的观察结果，返回默认消息
        return {"output": "Agent stopped due to iteration limit or time limit."}
