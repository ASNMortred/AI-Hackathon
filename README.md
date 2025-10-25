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
- 数据库驱动：go-sql-driver/mysql

**详细文档**：请查看 [server/README.md](server/README.md)

### Web端（前端应用）

前端应用目录，用于存放基于现代前端框架（React/Vue/Next.js）开发的用户界面。

**目录位置**：`web/`

## 快速开始

### 🐳 Docker部署（推荐）

使用Docker Compose一键部署前后端服务：

```bash
# 方法1: 使用部署脚本
./deploy.sh deploy

# 方法2: 使用Makefile
make deploy

# 方法3: 直接使用Docker Compose
docker-compose up -d
```

访问地址：
- 前端应用: http://localhost
- 后端API: http://localhost:8080/api/v1
- MySQL数据库: localhost:3306

详细说明请查看：
- 📖 [Docker部署完整指南](DOCKER_DEPLOYMENT.md)
- 🚀 [快速开始指南](QUICKSTART.md)

### 本地开发模式

#### 启动Server端

```bash
# 进入server目录
cd server

# 安装依赖
go mod download

# 运行服务
go run cmd/server/main.go
```

Server默认运行在 `http://localhost:8080`

#### 启动Web端

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
├── deploy.sh               # 一键部署脚本
├── Makefile                # Make命令快捷方式
├── .env.example            # 环境变量示例
├── DOCKER_DEPLOYMENT.md    # Docker部署指南
├── QUICKSTART.md           # 快速开始指南
├── server/              # Golang后端服务
│   ├── cmd/            # 应用程序入口
│   ├── internal/       # 内部代码
│   ├── configs/        # 配置文件
│   ├── Dockerfile      # 后端Docker构建文件
│   ├── .dockerignore   # Docker忽略文件
│   ├── go.mod          # Go模块依赖
│   └── README.md       # Server端详细文档
├── web/                # 前端应用
│   ├── src/            # 源代码
│   ├── Dockerfile      # 前端Docker构建文件
│   ├── nginx.conf      # Nginx配置
│   ├── .dockerignore   # Docker忽略文件
│   └── package.json    # NPM依赖
└── README.md           # 项目总体说明（本文件）
```

## 开发说明

- **Server端开发**：所有后端相关的开发工作在 `server/` 目录下进行
- **Web端开发**：所有前端相关的开发工作在 `web/` 目录下进行
- 两个模块可以独立开发、构建和部署

## 部署方式

### Docker部署（生产环境推荐）

使用Docker Compose实现一键部署，包含以下特性：

- ✅ 多阶段构建，优化镜像大小
- ✅ 依赖缓存，加速构建过程
- ✅ 健康检查，确保服务可用
- ✅ 自动重启，提高可靠性
- ✅ 日志管理，限制日志大小
- ✅ 数据持久化，保护上传文件

详情查看 [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)

### 传统部署

分别构建和部署前后端服务，详见各模块的README文档。

## 常用命令

```bash
# 使用Makefile（推荐）
make deploy      # 一键部署
make start       # 启动服务
make stop        # 停止服务
make logs        # 查看日志
make status      # 查看状态
make health      # 健康检查
make backup      # 备份数据

# 使用部署脚本
./deploy.sh deploy
./deploy.sh logs
./deploy.sh health

# 直接使用Docker Compose
docker-compose up -d
docker-compose logs -f
docker-compose down
```

## License

MIT License