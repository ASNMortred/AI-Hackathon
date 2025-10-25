# 🚀 Docker 快速开始指南

## 🇨🇳 国内用户优先配置（推荐）

为了加速构建，建议国内用户先配置镜像源：

```bash
# 一键配置所有国内镜像源
chmod +x setup-china-mirrors.sh
sudo ./setup-china-mirrors.sh

# 或者仅配置Docker镜像源
sudo ./setup-china-mirrors.sh --docker
```

详细说明请查看：[CHINA_MIRROR_CONFIG.md](CHINA_MIRROR_CONFIG.md)

## 30秒快速部署

```bash
# 方法1: 使用部署脚本
./deploy.sh deploy

# 方法2: 使用Makefile
make deploy

# 方法3: 直接使用Docker Compose
docker-compose up -d
```

## 访问应用

部署完成后，打开浏览器访问：

- **前端应用**: http://localhost
- **后端API**: http://localhost:8080/api/v1
- **健康检查**: http://localhost:8080/api/v1/health

## 常用命令速查

### 使用部署脚本 (deploy.sh)

```bash
./deploy.sh init      # 初始化环境
./deploy.sh build     # 构建镜像
./deploy.sh start     # 启动服务
./deploy.sh stop      # 停止服务
./deploy.sh logs      # 查看日志
./deploy.sh health    # 健康检查
./deploy.sh backup    # 备份数据
./deploy.sh clean     # 清理所有
```

### 使用Makefile

```bash
make init            # 初始化环境
make build           # 构建镜像
make start           # 启动服务
make stop            # 停止服务
make logs            # 查看日志
make status          # 查看状态
make health          # 健康检查
make backup          # 备份数据
make clean           # 清理所有
```

### 使用Docker Compose

```bash
docker-compose up -d              # 启动服务
docker-compose down               # 停止服务
docker-compose logs -f            # 查看日志
docker-compose ps                 # 查看状态
docker-compose restart            # 重启服务
docker-compose up -d --build      # 重新构建并启动
```

## 测试API

### 健康检查

```bash
curl http://localhost:8080/api/v1/health
```

### 上传文件

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@test.mp4"
```

### 播放视频

```bash
curl http://localhost:8080/api/v1/play/video123
```

### 发送消息

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}'
```

## 故障排查

### 查看日志

```bash
# 查看所有日志
make logs

# 只看后端日志
make logs-server

# 只看前端日志
make logs-web
```

### 检查服务状态

```bash
# 查看容器状态
make status

# 健康检查
make health
```

### 重启服务

```bash
# 重启所有服务
make restart

# 重新构建并启动
make rebuild
```

## 停止和清理

```bash
# 停止服务
make stop

# 停止并删除所有（包括数据）
make clean
```

## 需要帮助？

查看完整文档：[DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)

## 文件结构

```
.
├── docker-compose.yml       # Docker编排配置
├── deploy.sh               # 部署脚本
├── Makefile                # Make命令
├── .env.example            # 环境变量示例
├── server/
│   ├── Dockerfile          # 后端Dockerfile
│   └── .dockerignore       # 后端忽略文件
└── web/
    ├── Dockerfile          # 前端Dockerfile
    ├── nginx.conf          # Nginx配置
    └── .dockerignore       # 前端忽略文件
```

## 推荐工作流

### 开发环境

```bash
# 1. 初始化
make init

# 2. 启动服务（实时日志）
make dev
```

### 生产环境

```bash
# 1. 一键部署
make deploy

# 2. 定期备份
make backup

# 3. 监控日志
make logs
```

## 环境要求

- Docker ≥ 20.10
- Docker Compose ≥ 2.0

检查版本：
```bash
docker --version
docker-compose --version
```
