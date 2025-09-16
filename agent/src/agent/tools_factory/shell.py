
from pydantic import Field
from src.agent.tools_factory.tools_registry import BaseToolOutput, register_tool
import subprocess
from src.utils.logger import logger


@register_tool(title="系统命令")
def shell(query: str = Field(description="The command to execute")) -> BaseToolOutput:
    """Use Shell to execute system shell commands."""
    try:
        logger.info(f"执行shell命令: {query}")
        result = subprocess.run(query, shell=True, capture_output=True, text=True, timeout=30)
        output = result.stdout if result.returncode == 0 else result.stderr
        logger.info(f"命令执行结果: {output}")
        
        tool_output = BaseToolOutput(data=output)
        logger.info(f"BaseToolOutput.__str__() 结果: {str(tool_output)}")
        
        return tool_output
    except Exception as e:
        error_msg = str(e)
        logger.error(f"命令执行出错: {error_msg}")
        return BaseToolOutput(data=error_msg)