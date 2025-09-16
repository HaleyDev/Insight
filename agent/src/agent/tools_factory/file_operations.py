from src.utils.pydantic_v1 import Field

from src.agent.tools_factory.tools_registry import BaseToolOutput, register_tool


@register_tool(title="读取代码文件源码")
def read_code_file(file_path: str = Field(description="The path to the code file")) -> str:
    """Read the content of a code file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    return BaseToolOutput(content, format="str")