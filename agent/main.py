#!/usr/bin/env python3

import sys
import os
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from src.agent.llms.llms import get_llm
from src.agent.agent_factory.agents_registry import agents_registry
from src.agent.tools_factory.shell import shell
from src.utils.logger import logger

def main():
    llm = get_llm("deepseek")
    
    # 更详细的日志记录
    logger.info("创建代理...")
    agent = agents_registry(llm=llm, tools=[shell], verbose=True)
    logger.info("代理创建完成")
    
    # 添加更多调试信息
    logger.info("开始执行代理...")
    try:
        # 使用更简单的提问方式
        user_input = "使用shell工具执行pwd命令，并告诉我当前目录"
        logger.info(f"用户输入: {user_input}")
        
        response = agent.invoke(
            {"input": user_input}, 
            config={"max_iterations": 5, "timeout": 60}  # 增加迭代次数和超时时间
        )
        logger.info(f"代理执行完成，响应类型: {type(response)}")
        logger.info(f"完整响应内容: {response}")
        
        # 提取并显示结果
        if isinstance(response, dict) and 'output' in response:
            output = response['output']
            logger.info(f"提取的输出: {output}")
            print("\n=== 执行结果 ===")
            print(output)
        else:
            logger.warning(f"意外的响应格式: {response}")
            print(f"响应: {response}")
            
    except Exception as e:
        logger.error(f"代理执行出错: {e}")
        import traceback
        logger.error(traceback.format_exc())
        print(f"执行错误: {e}")

if __name__ == "__main__":
    main()