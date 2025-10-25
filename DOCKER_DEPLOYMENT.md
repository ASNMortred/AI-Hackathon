# Docker 部署指南

## 📋 项目概述

本项目采用Docker Compose进行容器化部署，包含以下服务：
- **server**: Golang后端服务（基于Gin框架）
- **web**: React前端应用（基于Vite + Ant Design）

## 🚀 快速开始

### 前置要求

确保您的系统已安装：
- Docker (≥ 20.10)
- Docker Compose (≥ 2.0)

检查安装：
```bash
docker --version
docker-compose --version
```

### 一键启动

在项目根目录执行：

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 访问应用

- **前端应用**: http://localhost
- **后端API**: http://localhost:8080/api/v1
- **健康检查**: http://localhost:8080/api/v1/health

## 🏗️ 构建说明

### 后端服务 (server)

**Dockerfile特性**：
- ✅ 多阶段构建（golang:1.24.5-alpine + alpine:latest）
- ✅ 优化依赖缓存（先复制go.mod/go.sum）
- ✅ 静态编译（CGO_ENABLED=0）
- ✅ 精简镜像（仅包含二进制文件和CA证书）
- ✅ 健康检查（/api/v1/health）

**手动构建**：
```bash
cd server
docker build -t ai-hackathon-server:latest .
```

### 前端服务 (web)

**Dockerfile特性**：
- ✅ 多阶段构建（node:18-alpine + nginx:alpine）
- ✅ 优化依赖缓存（先复制package.json）
- ✅ 生产构建（npm ci + npm run build）
- ✅ Nginx反向代理（自动转发API请求到后端）
- ✅ SPA路由支持
- ✅ Gzip压缩和静态资源缓存

**手动构建**：
```bash
cd web
docker build -t ai-hackathon-web:latest .
```

## 🔧 配置说明

### 环境变量（可选）

创建 `.env` 文件在项目根目录：

```env
# 服务端口配置
SERVER_PORT=8080
WEB_PORT=80

# 时区设置
TZ=Asia/Shanghai
```

然后在 `docker-compose.yml` 中引用：
```yaml
services:
  server:
    env_file:
      - .env
```

### 配置文件修改

**后端配置**：修改 `server/configs/config.yaml`

```yaml
server:
  port: "8080"

upload:
  max_size: 524288000  # 500MB
  upload_dir: "uploads"
```

**前端API代理**：修改 `web/nginx.conf`

```nginx
location /api/ {
    proxy_pass http://server:8080;  # 指向后端服务
}
```

## 📦 数据持久化

### 上传文件持久化

上传的文件默认存储在主机的 `./server/uploads` 目录：

```yaml
volumes:
  - ./server/uploads:/app/uploads
```

### 自定义挂载点

修改 `docker-compose.yml`：

```yaml
volumes:
  - /your/custom/path:/app/uploads
```

## 🛠️ 常用命令

### 服务管理

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 停止并删除所有数据（包括volumes）
docker-compose down -v
```

### 日志查看

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f server
docker-compose logs -f web

# 查看最近100行日志
docker-compose logs --tail=100 server
```

### 服务操作

```bash
# 进入容器shell
docker-compose exec server sh
docker-compose exec web sh

# 查看容器资源使用
docker stats

# 重新构建并启动
docker-compose up -d --build

# 仅重建特定服务
docker-compose up -d --build server
```

### 健康检查

```bash
# 检查后端健康状态
curl http://localhost:8080/api/v1/health

# 查看容器健康状态
docker-compose ps
```

## 🔍 故障排查

### 检查服务状态

```bash
# 查看所有容器状态
docker-compose ps

# 查看详细信息
docker-compose logs server
docker-compose logs web
```

### 常见问题

**1. 端口冲突**

错误信息：`port is already allocated`

解决方法：修改 `docker-compose.yml` 中的端口映射
```yaml
ports:
  - "8081:8080"  # 改用其他端口
```

**2. 构建失败**

```bash
# 清理Docker缓存
docker system prune -a

# 重新构建
docker-compose build --no-cache
```

**3. 网络连接问题**

```bash
# 检查网络
docker network ls
docker network inspect ai-hackathon_ai-hackathon-network

# 重建网络
docker-compose down
docker network prune
docker-compose up -d
```

**4. 权限问题**

```bash
# 检查uploads目录权限
ls -la server/uploads

# 修改权限
chmod 777 server/uploads
```

## 🎯 生产部署建议

### 1. 使用环境变量管理敏感配置

```bash
# 创建 .env 文件
cat > .env << EOF
SERVER_PORT=8080
DATABASE_URL=postgresql://...
SECRET_KEY=your-secret-key
EOF
```

### 2. 配置SSL证书（Nginx）

修改 `web/nginx.conf`：
```nginx
server {
    listen 443 ssl;
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
}
```

挂载证书：
```yaml
web:
  volumes:
    - ./ssl:/etc/nginx/ssl:ro
```

### 3. 资源限制

在 `docker-compose.yml` 中添加：
```yaml
services:
  server:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 4. 日志管理

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

### 5. 备份策略

```bash
# 备份上传文件
tar -czf uploads-backup-$(date +%Y%m%d).tar.gz server/uploads

# 备份配置
tar -czf config-backup-$(date +%Y%m%d).tar.gz server/configs
```

## 📊 性能优化

### 1. 构建优化

```dockerfile
# 使用构建缓存
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
```

### 2. 镜像大小优化

```bash
# 查看镜像大小
docker images | grep ai-hackathon

# 清理未使用的镜像
docker image prune -a
```

### 3. Nginx缓存优化

已在 `web/nginx.conf` 中配置：
- Gzip压缩
- 静态资源缓存（1年）
- 代理缓存

## 🔐 安全建议

1. **不要在镜像中硬编码敏感信息**
2. **定期更新基础镜像**
3. **使用非root用户运行容器**
4. **限制容器资源使用**
5. **启用Docker内容信任**
6. **定期扫描镜像漏洞**

```bash
# 扫描镜像漏洞
docker scan ai-hackathon-server:latest
```

## 📝 开发模式

### 开发环境配置

创建 `docker-compose.dev.yml`：

```yaml
version: '3.8'

services:
  server:
    volumes:
      - ./server:/app
    command: go run cmd/server/main.go
    
  web:
    volumes:
      - ./web:/app
    command: npm run dev
```

启动开发环境：
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## 🆘 获取帮助

- 查看容器日志：`docker-compose logs -f`
- 检查健康状态：`docker-compose ps`
- 进入容器调试：`docker-compose exec server sh`

## 📄 文件清单

```
.
├── docker-compose.yml          # Docker Compose编排文件
├── server/
│   ├── Dockerfile             # 后端Dockerfile
│   └── .dockerignore          # 后端忽略文件
└── web/
    ├── Dockerfile             # 前端Dockerfile
    ├── nginx.conf             # Nginx配置
    └── .dockerignore          # 前端忽略文件
```

## ⚖️ 许可证

MIT License
