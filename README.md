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

**详细文档**：请查看 [server/README.md](server/README.md)

### Web端（前端应用）

前端应用目录，用于存放基于现代前端框架（React/Vue/Next.js）开发的用户界面。

**目录位置**：`web/`

## 快速开始

### 启动Server端

```bash
# 进入server目录
cd server

# 安装依赖
go mod download

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
├── server/              # Golang后端服务
│   ├── cmd/            # 应用程序入口
│   ├── internal/       # 内部代码
│   ├── configs/        # 配置文件
│   ├── go.mod          # Go模块依赖
│   └── README.md       # Server端详细文档
├── web/                # 前端应用
│   └── .gitignore      # 前端忽略规则
└── README.md           # 项目总体说明（本文件）
```

## 开发说明

- **Server端开发**：所有后端相关的开发工作在 `server/` 目录下进行
- **Web端开发**：所有前端相关的开发工作在 `web/` 目录下进行
- 两个模块可以独立开发、构建和部署

## License

MIT License