# AI-Hackathon Web Service

基于Gin框架的Web服务项目，提供文件上传、视频播放、聊天功能以及用户认证（登录/注册）的API接口。

## 项目结构

```
.
├── cmd/
│   └── server/          # 应用程序入口
│       └── main.go
├── internal/
│   ├── config/          # 配置管理
│   │   └── config.go
│   ├── handlers/        # API处理器
│   │   ├── upload.go
│   │   ├── play.go
│   │   ├── chat.go
│   │   └── user.go
│   ├── logger/          # 日志系统
│   │   └── logger.go
│   ├── middleware/      # 中间件
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── database/        # 数据库连接
│   │   └── database.go
│   ├── models/          # 数据模型
│   │   └── user.go
│   ├── dao/             # 数据访问对象
│   │   └── user.go
│   └── services/        # 业务逻辑层
│       └── user.go
├── configs/             # 配置文件
│   └── config.yaml
├── uploads/             # 上传文件存储目录
└── bin/                 # 编译后的二进制文件
```

## 技术栈

- **Web框架**: Gin
- **配置管理**: Viper + Pflag
- **日志系统**: Zap
- **数据库**: MySQL
- **数据库驱动**: go-sql-driver/mysql
- **语言版本**: Go 1.24.5

## 功能特性

### 1. 文件上传 (`POST /api/v1/upload`)
- 支持multipart/form-data格式上传
- 文件大小限制：500MB
- 支持格式：视频(.mp4, .avi, .mov, .mkv)、图片(.jpg, .jpeg, .png, .gif)、音频(.mp3, .wav)、压缩包(.zip, .rar, .7z)
- 文件类型和大小校验
- 自动创建上传目录

### 2. 视频播放 (`GET /api/v1/play/:videoID`)
- 模拟播放流程
- 控制台输出播放信息
- 返回JSON响应

### 3. 聊天接口 (`POST /api/v1/chat`)
- 接收用户消息
- 控制台输出接收信息
- 返回确认响应

### 4. 用户注册 (`POST /api/v1/register`)
- 用户名和密码注册
- 用户名唯一性验证
- 数据库存储用户信息

### 5. 用户登录 (`POST /api/v1/login`)
- 用户名密码验证
- 返回认证token
- 基于数据库的用户验证

## 数据库表结构

```sql
CREATE TABLE `users` (
  `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识(主键)',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '加密后的密码',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

## 配置说明

配置文件位于 `configs/config.yaml`：

```yaml
server:
  port: "8080"                    # 服务端口

qiniu:
  access_key: "your_access_key"   # 七牛云AccessKey
  secret_key: "your_secret_key"   # 七牛云SecretKey
  bucket: "your_bucket_name"      # 存储桶名称

upload:
  max_size: 524288000             # 最大文件大小（字节，默认500MB）
  allowed_types:                  # 允许的文件类型
    - ".mp4"
    - ".avi"
    # ...
  upload_dir: "uploads"           # 上传目录

database:
  host: "localhost"               # 数据库主机
  port: "3306"                    # 数据库端口
  user: "root"                    # 数据库用户名
  password: "password"            # 数据库密码
  name: "ai_hackathon"            # 数据库名称
```

## 构建和运行

### 1. 安装依赖

```bash
go mod download
```

### 2. 数据库准备

确保MySQL服务正在运行，并创建数据库：

```sql
CREATE DATABASE ai_hackathon;
USE ai_hackathon;

CREATE TABLE `users` (
  `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识(主键)',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '加密后的密码',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

### 3. 构建项目

```bash
go build -o bin/server ./cmd/server
```

### 4. 运行服务

使用默认配置文件：
```bash
./bin/server
```

指定配置文件：
```bash
./bin/server --config /path/to/config.yaml
```

### 5. 开发模式运行

```bash
go run cmd/server/main.go
```

## API使用示例

### 上传文件

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@/path/to/your/file.mp4"
```

响应：
```json
{
  "message": "File uploaded successfully",
  "filename": "file.mp4",
  "size": 1048576,
  "path": "uploads/file.mp4"
}
```

### 播放视频

```bash
curl http://localhost:8080/api/v1/play/video123
```

响应：
```json
{
  "message": "Playing video",
  "video_id": "video123"
}
```

### 发送聊天消息

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，世界！"}'
```

响应：
```json
{
  "message": "Message received",
  "echo": "你好，世界！"
}
```

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpass123"}'
```

响应：
```json
{
  "message": "User registered successfully",
  "username": "testuser"
}
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpass123"}'
```

响应：
```json
{
  "message": "Login successful",
  "username": "testuser",
  "token": "token_testuser"
}
```

## 日志系统

项目使用Zap进行结构化日志记录：
- 请求日志：记录每个HTTP请求的详细信息（方法、路径、状态码、延迟等）
- 错误日志：记录应用程序错误和异常
- 业务日志：记录业务操作（文件上传、消息接收等）

## 错误处理

- 全局panic恢复中间件
- 统一的错误响应格式
- 详细的错误日志记录

## 开发说明

1. 代码结构遵循标准Go项目布局
2. 使用依赖注入模式
3. 配置与代码分离
4. 完善的错误处理机制
5. 结构化日志记录
6. 数据库访问层分离

## 后续扩展

当前版本为基础架构，后续可扩展：
- 集成七牛云存储
- 实现真实的视频播放服务
- 集成MCP服务进行智能对话
- 添加用户认证和授权
- 数据库集成
- 完善的测试覆盖

## License

MIT License