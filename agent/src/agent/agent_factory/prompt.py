"""
Prompt templates for different AI models and agent types.
"""


def get_prompt_template(model_type: str, prompt_name: str) -> str:
    """
    Get prompt template for different model types and prompt names.
    
    Args:
        model_type: Type of model (e.g., "action_model")
        prompt_name: Name of the prompt template (e.g., "qwen", "deepseek", "structured-chat-agent")
    
    Returns:
        str: The formatted prompt template
    """
    if model_type == "action_model":
        if prompt_name == "qwen":
            return """工具: {tools}

格式:
问题: {input}
Action: [{tool_names}]之一
Action Input: JSON格式  
Observation: 结果
Final Answer: 最终答案

{agent_scratchpad}

问题: {input}
Action:"""

        elif prompt_name == "deepseek":
            return """工具: {tools}

格式:
Question: {input}
Action: [{tool_names}]之一  
Action Input: JSON格式
Observation: 结果
Final Answer: 最终答案

注意事项:
1. 执行工具后，必须分析结果
2. 若结果满足需求，立即提供最终答案
3. 不要重复执行相同的命令
4. 格式为"Final Answer: 你的答案"

{agent_scratchpad}

Question: {input}
Action:"""

        elif prompt_name == "structured-chat-agent":
            return """You are an AI assistant that can use tools to help answer questions.

Respond to the human as helpfully and accurately as possible. You have access to the following tools:

{tools}

Use a json blob to specify a tool by providing an action key (tool name) and an action_input key (tool input).

Valid "action" values: "Final Answer" or {tool_names}

Provide only ONE action per $JSON_BLOB, as shown:

```
{{
  "action": $TOOL_NAME,
  "action_input": $INPUT
}}
```

Follow this format:

Question: input question to answer
Thought: consider previous and subsequent steps
Action:
```
$JSON_BLOB
```
Observation: action result
... (repeat Thought/Action/Observation N times)
Thought: I know what to respond
Action:
```
{{
  "action": "Final Answer",
  "action_input": "Final response to human"
}}
```

Begin! Reminder to ALWAYS respond with a valid json blob of a single action. Use tools if necessary. Respond directly if appropriate. Format is Action:```$JSON_BLOB```then Observation:.

{agent_scratchpad}

Question: {input}
Thought:"""
    
    # Default fallback template
    return """You are a helpful AI assistant.

Question: {input}
Answer:"""


def get_react_prompt_template() -> str:
    """
    Get ReAct style prompt template.
    
    Returns:
        str: The ReAct prompt template
    """
    return """You are an AI assistant that can use tools to help answer questions. You have access to the following tools:

{tools}

Use the following format:

Question: the input question you must answer
Thought: you should always think about what to do
Action: the action to take, should be one of [{tool_names}]
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I now know the final answer
Final Answer: the final answer to the original input question

Begin!

Question: {input}
Thought: {agent_scratchpad}
"""


def get_custom_prompt_template(model_name: str, task_type: str = "general") -> str:
    """
    Get custom prompt template for specific models and tasks.
    
    Args:
        model_name: Name of the model (e.g., "qwen", "deepseek", "chatglm")
        task_type: Type of task (e.g., "general", "code", "analysis")
    
    Returns:
        str: The custom prompt template
    """
    if model_name.lower() == "qwen":
        if task_type == "code":
            return """你是一个专业的代码助手，可以帮助用户解决编程问题。

问题: {input}
请提供详细的代码解答:"""
        elif task_type == "analysis":
            return """你是一个数据分析专家，可以帮助用户分析和理解数据。

问题: {input}
请提供详细的分析:"""
    
    elif model_name.lower() == "deepseek":
        if task_type == "code":
            return """You are a professional coding assistant that can help users solve programming problems.

Question: {input}
Please provide a detailed code solution:"""
        elif task_type == "analysis":
            return """You are a data analysis expert that can help users analyze and understand data.

Question: {input}
Please provide a detailed analysis:"""
    
    # Default general template
    return """You are a helpful AI assistant.

Question: {input}
Answer:"""
