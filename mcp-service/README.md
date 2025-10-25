# MCP Chat Service

基于 FastAPI 的 MCP (Model Context Protocol) 聊天服务。

## 功能特性

- ✅ FastAPI 框架，支持完整的类型注解
- ✅ MCP 协议规范的工具注册
- ✅ OpenAI API 集成（可选，支持自定义 base_url）
- ✅ 会话管理，维护多轮对话上下文
- ✅ 结构化日志记录
- ✅ 健康检查端点
- ✅ Docker 化部署

## API 接口

### 1. 聊天接口

**POST** `/api/chat`

请求体:
```json
{
  "message": "用户输入",
  "session_id": "可选会话ID",
  "temperature": 0.7
}
```

响应:
```json
{
  "success": true,
  "data": {
    "response": "AI回复",
    "session_id": "会话ID"
  },
  "error": null
}
```

### 2. 工具列表

**GET** `/mcp/tools`

返回注册的 MCP 工具列表。

### 3. 健康检查

**GET** `/health`

返回服务健康状态。

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `MCP_SERVICE_PORT` | 服务端口 | 8000 |
| `OPENAI_API_KEY` | OpenAI API 密钥 | - |
| `OPENAI_BASE_URL` | OpenAI API 基础URL | https://api.openai.com/v1 |
| `OPENAI_MODEL` | 使用的模型 | gpt-3.5-turbo |

## 本地开发

```bash
# 安装依赖
pip install -r requirements.txt

# 运行服务
python main.py
```

## Docker 部署

服务已集成到 `docker-compose.yml` 中，使用以下命令启动：

```bash
docker-compose up -d mcp-service
```

## MCP 工具

当前已注册的工具：

1. **get_current_time**: 获取当前时间
2. **search_knowledge**: 搜索知识库（模拟实现）

## 日志

日志文件存储在 `/app/logs/mcp-service.log`，包含：

- 请求/响应日志
- 错误链路追踪
- 服务状态信息

## 架构说明

本服务作为 AI 聊天的专用层，与 Golang API 网关配合使用：

```
客户端 -> Golang API (8080) -> MCP Service (8000) -> OpenAI API
```

Golang 层负责：
- 请求转发
- 认证
- 限流

Python MCP 层负责：
- AI 聊天逻辑
- 工具调用
- 会话管理
