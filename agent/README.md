# CI-Vision Agent

基于 LangChain 的智能代理系统，支持多种大语言模型（DeepSeek、Qwen）和工具集成的 AI 代理平台。

## 功能特性

- 🤖 **多模型支持**: 集成 DeepSeek 和 Qwen 大语言模型
- 🔧 **工具集成**: 支持 Shell 命令执行、文件操作等工具
- 🚀 **增强执行器**: 优化的代理执行器，更好地处理迭代限制和错误恢复
- 📊 **完善日志**: 详细的执行日志和调试信息
- 🌐 **gRPC 支持**: 提供 gRPC 接口用于服务间通信
- ⚙️ **灵活配置**: 支持 YAML 配置文件和环境变量

## 项目结构

```
agent/
├── .env                     # 环境变量配置
├── .venv/                   # Python虚拟环境
├── main.py                  # 程序入口文件
├── requirements.txt         # Python依赖包列表
├── config/                  # 配置管理模块
│   ├── __init__.py
│   ├── config.py           # 配置类实现
│   ├── hhw_config.yaml     # YAML配置文件
│   └── path.py             # 路径配置
├── src/                     # 核心源代码
│   ├── agent/              # AI智能代理模块
│   │   ├── agent_factory/  # 代理工厂模块
│   │   │   ├── agents_registry.py    # 代理注册器
│   │   │   ├── deepseek_agent.py     # DeepSeek代理实现
│   │   │   ├── qwen_agent.py         # Qwen代理实现
│   │   │   ├── enhanced_executor.py  # 增强执行器
│   │   │   └── prompt.py             # 提示词模板
│   │   ├── llms/           # 语言模型模块
│   │   │   └── llms.py     # 模型获取和配置
│   │   └── tools_factory/  # 工具工厂模块
│   │       ├── tools_registry.py     # 工具注册器
│   │       ├── shell.py              # Shell命令工具
│   │       └── file_operations.py   # 文件操作工具
│   ├── grpc/               # gRPC通信模块
│   │   ├── client.py       # gRPC客户端
│   │   ├── service.py      # gRPC服务端
│   │   └── proto/          # Protocol Buffers定义
│   ├── models/             # 数据模型定义
│   ├── services/           # 业务服务层
│   ├── settings/           # 系统设置模块
│   └── utils/              # 工具函数
│       ├── logger.py       # 日志工具
│       └── pydantic_v1.py  # Pydantic v1兼容层
├── docs/                   # 项目文档
├── logs/                   # 运行日志目录
├── test/                   # 测试示例
└── tests/                  # 单元测试
    ├── test_agent/         # 代理模块测试
    └── test_grpc/          # gRPC模块测试
```

## 快速开始

### 环境要求

- Python 3.8+
- pip 或 conda

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd agent
```

2. **创建虚拟环境**
```bash
python -m venv .venv
source .venv/bin/activate  # Linux/macOS
# 或
.venv\Scripts\activate     # Windows
```

3. **安装依赖**
```bash
pip install -r requirements.txt
```

4. **配置环境**
```bash
# 复制并编辑环境配置文件
cp .env.example .env
# 编辑 .env 文件，设置必要的环境变量
```

5. **运行程序**
```bash
python main.py
```

## 配置说明

### 环境变量

在 `.env` 文件中配置：

- `ENV`: 环境标识（如 `hhw`）
- `DEEPSEEK_API_KEY`: DeepSeek API 密钥
- `QWEN_API_KEY`: Qwen API 密钥

### YAML 配置

在 `config/hhw_config.yaml` 中配置详细参数：

- 模型配置
- API 端点
- 工具设置
- 日志级别

## 核心组件

### 代理系统

- **代理注册器**: 根据模型类型自动选择合适的代理实现
- **DeepSeek 代理**: 针对 DeepSeek 模型优化的代理
- **Qwen 代理**: 针对 Qwen 模型优化的代理
- **增强执行器**: 改进的执行逻辑，更好地处理迭代限制

### 工具系统

- **Shell 工具**: 执行系统命令
- **文件操作工具**: 文件读写、创建、删除等操作
- **工具注册器**: 动态工具加载和管理

### gRPC 服务

- **客户端**: 用于与其他服务通信
- **服务端**: 提供代理功能的 gRPC 接口
- **Protocol Buffers**: 定义服务接口和数据结构

## 使用示例

```python
from src.agent.llms.llms import get_llm
from src.agent.agent_factory.agents_registry import agents_registry
from src.agent.tools_factory.shell import shell

# 获取语言模型
llm = get_llm("deepseek")

# 创建代理
agent = agents_registry(llm=llm, tools=[shell], verbose=True)

# 执行任务
response = agent.invoke({
    "input": "使用shell工具执行pwd命令，并告诉我当前目录"
})

print(response['output'])
```

## 开发指南

### 添加新模型

1. 在 `src/agent/agent_factory/` 中创建新的代理实现
2. 在 `agents_registry.py` 中注册新模型
3. 更新配置文件支持新模型

### 添加新工具

1. 在 `src/agent/tools_factory/` 中实现工具类
2. 在 `tools_registry.py` 中注册新工具
3. 更新代理配置以包含新工具

## 故障排除

### 常见问题

1. **代理停止执行**: 检查迭代限制设置和增强执行器日志
2. **模型 API 错误**: 验证 API 密钥和网络连接
3. **工具执行失败**: 检查工具权限和参数格式

### 日志分析

查看 `logs/` 目录中的日志文件，包含：
- 代理执行步骤
- 工具调用记录
- 错误堆栈信息

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/new-feature`)
3. 提交更改 (`git commit -am 'Add new feature'`)
4. 推送到分支 (`git push origin feature/new-feature`)
5. 创建 Pull Request

## 许可证

[添加许可证信息]

## 联系方式

[添加联系信息]
