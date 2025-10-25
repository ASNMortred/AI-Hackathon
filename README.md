# AI-Hackathon

一个包含后端Server和前端Web应用的全栈项目。

## 项目概述

本项目采用前后端分离架构，由以下两个主要模块组成：

### Server端（后端服务）

基于Golang + Gin框架开发的Web服务，提供文件上传、视频播放、聊天功能以及用户认证（登录/注册）的RESTful API接口。

**技术栈**：
- 语言：Go 1.20
- Web框架：Gin
- 配置管理：Viper + Pflag
- 日志系统：Zap
- 数据库：MySQL 8.0
- 对象存储：MinIO
- 数据库驱动：go-sql-driver/mysql

**详细文档**：请查看 [server/README.md](server/README.md)

### Web端（前端应用）

前端应用目录，用于存放基于现代前端框架（React/Vue/Next.js）开发的用户界面。

**目录位置**：`web/`

## 快速开始

### 🐳 Docker部署（推荐）

使用Docker Compose一键部署所有服务（包括MySQL和MinIO）：

```bash
# 启动所有服务
docker-compose up -d
```

访问地址：
- 前端应用: http://localhost
- 后端API: http://localhost:8080/api/v1
- MySQL数据库: localhost:3306
- MinIO控制台: http://localhost:9001
- MinIO API: http://localhost:9000

默认MinIO凭证（可通过环境变量修改）：
- Access Key: `minioadmin`
- Secret Key: `minioadmin`

### 本地开发模式

#### 启动Server端

```bash
# 进入server目录
cd server

# 安装依赖
go mod download

# 配置环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=hackathon
export DB_PASSWORD=hackathon123
export DB_NAME=ai_hackathon
export MINIO_ENDPOINT=localhost:9000
export MINIO_ACCESS_KEY=minioadmin
export MINIO_SECRET_KEY=minioadmin
export MINIO_BUCKET=uploads

# 运行服务
go run cmd/server/main.go
```

Server默认运行在 `http://localhost:8080`

### 启动Web端

```bash
# 进入web目录
cd web

# 安装依赖（以npm为例）
npm install

# 运行开发服务器
npm run dev
```

## 项目结构

```
.
├── docker-compose.yml      # Docker编排配置
├── init.sql                # 数据库初始化脚本
├── server/              # Golang后端服务
│   ├── cmd/            # 应用程序入口
│   ├── internal/       # 内部代码
│   │   ├── handlers/   # HTTP处理器
│   │   ├── storage/    # MinIO存储服务
│   │   ├── dao/        # 数据访问层
│   │   └── config/     # 配置管理
│   ├── configs/        # 配置文件
│   ├── Dockerfile      # 后端Docker构建文件
│   ├── go.mod          # Go模块依赖
│   └── README.md       # Server端详细文档
├── web/                # 前端应用
│   ├── src/            # 源代码
│   ├── Dockerfile      # 前端Docker构建文件
│   └── package.json    # NPM依赖
└── README.md           # 项目总体说明（本文件）
```

## 开发说明

- **Server端开发**：所有后端相关的开发工作在 `server/` 目录下进行
- **Web端开发**：所有前端相关的开发工作在 `web/` 目录下进行
- 两个模块可以独立开发、构建和部署

## 部署架构

项目使用Docker Compose编排以下服务：

1. **MySQL** - 用户数据和文件元数据存储
2. **MinIO** - 对象存储服务，用于文件上传
3. **Server** - Go后端服务，依赖MySQL和MinIO
4. **Web** - React前端应用

所有服务通过Docker网络互联，支持健康检查和自动重启。

## 环境变量配置

可通过创建 `.env` 文件自定义配置：

```bash
# 数据库配置
DB_ROOT_PASSWORD=rootpassword
DB_NAME=ai_hackathon
DB_USER=hackathon
DB_PASSWORD=hackathon123

# MinIO配置
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=uploads
```

## 常用命令

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

## License

MIT License