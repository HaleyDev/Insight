import os.path

from langchain_core.tools import StructuredTool
from pydantic import BaseModel, Field


class ReadFileInput(BaseModel):
    file_path: str = Field(description="The path to the code file to read")

def read_code_file_func(file_path: str) -> str:
    """Read the content of a code file."""
    try:
        if not os.path.exists(file_path):
            return f"Error: File '{file_path}' does not exist."

        if not os.path.isfile(file_path):
            return f"Error: '{file_path}' is not a file."

        file_size = os.path.getsize(file_path)
        if file_size > 1024 * 1024:
            return f"Error: File is too large ({file_size} bytes). Maximum allowed size is 1MB."

        # 检查文件扩展名，确保是代码文件
        code_extensions = {'.py', '.js', '.java', '.cpp', '.c', '.h', '.html',
                          '.css', '.ts', '.go', '.rs', '.rb', '.php', '.swift'}
        file_ext = os.path.splitext(file_path)[1].lower()

        if file_ext not in code_extensions:
            return f"Warning: File extension '{file_ext}' is not a common code file extension. Proceeding anyway..."

        with open(file_path, "r", encoding="utf-8") as f:
            context = f.read()

        return f"Content of code file '{file_path}' :\n\n{context}"

    except Exception as e:
        return f"Error reading file '{file_path}': {str(e)}"

# 创建工具
read_code_tool = StructuredTool.from_function(
    func=read_code_file_func,
    name="read_code_file",
    description="Read the content of a code file. Supports Python, JavaScript, Java, C++, HTML, CSS and other common code files.",
    args_schema=ReadFileInput,
    return_direct=False,
)