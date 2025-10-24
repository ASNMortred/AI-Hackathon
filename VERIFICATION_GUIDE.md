# ✅ Docker部署验证指南

本指南帮助您验证Docker化部署方案的完整性和正确性。

## 📋 验证清单

### 第一步：文件完整性检查

#### 1.1 核心文件检查

```bash
# 检查根目录文件
ls -la docker-compose.yml deploy.sh Makefile .env.example

# 检查后端文件
ls -la server/Dockerfile server/.dockerignore

# 检查前端文件
ls -la web/Dockerfile web/.dockerignore web/nginx.conf

# 检查文档文件
ls -la *.md
```

**预期输出**：所有文件都应该存在

#### 1.2 脚本权限检查

```bash
# 查看deploy.sh权限
ls -l deploy.sh

# 如果没有执行权限，添加权限
chmod +x deploy.sh
```

**预期**: `-rwxr-xr-x` (可执行)

### 第二步：环境检查

#### 2.1 Docker环境检查

```bash
# 检查Docker版本
docker --version

# 检查Docker Compose版本
docker-compose --version

# 检查Docker服务状态
docker info
```

**最低要求**:
- Docker: ≥ 20.10
- Docker Compose: ≥ 2.0

#### 2.2 端口可用性检查

```bash
# 检查80端口
lsof -i :80 || echo "端口80可用"

# 检查8080端口
lsof -i :8080 || echo "端口8080可用"
```

**预期**: 两个端口都应该可用

#### 2.3 磁盘空间检查

```bash
# 检查可用磁盘空间
df -h .
```

**最低要求**: 至少2GB可用空间

### 第三步：配置验证

#### 3.1 Dockerfile语法检查

```bash
# 验证后端Dockerfile
docker build -f server/Dockerfile --no-cache -t test-server server/ --target builder

# 验证前端Dockerfile
docker build -f web/Dockerfile --no-cache -t test-web web/ --target builder
```

**预期**: 构建应该成功

#### 3.2 Docker Compose配置检查

```bash
# 验证docker-compose.yml语法
docker-compose config

# 验证服务定义
docker-compose config --services
```

**预期输出**:
```
server
web
```

#### 3.3 Nginx配置检查

```bash
# 检查nginx.conf语法
docker run --rm -v $(pwd)/web/nginx.conf:/etc/nginx/conf.d/default.conf nginx:alpine nginx -t
```

**预期**: `nginx: configuration file ... test is successful`

### 第四步：构建验证

#### 4.1 镜像构建测试

```bash
# 方式1: 使用make
make build

# 方式2: 使用部署脚本
./deploy.sh build

# 方式3: 使用docker-compose
docker-compose build
```

**预期**: 两个镜像构建成功

#### 4.2 验证镜像

```bash
# 查看构建的镜像
docker images | grep ai-hackathon

# 查看镜像详情
docker inspect ai-hackathon-server:latest
docker inspect ai-hackathon-web:latest
```

**预期输出示例**:
```
ai-hackathon-server   latest   xxx   10-20MB
ai-hackathon-web      latest   xxx   20-30MB
```

### 第五步：部署验证

#### 5.1 启动服务

```bash
# 方式1: 使用make（推荐）
make deploy

# 方式2: 使用部署脚本
./deploy.sh deploy

# 方式3: 使用docker-compose
docker-compose up -d
```

#### 5.2 检查容器状态

```bash
# 查看运行中的容器
docker-compose ps

# 或使用make
make status
```

**预期输出**:
```
NAME                    STATUS              PORTS
ai-hackathon-server     Up (healthy)        0.0.0.0:8080->8080/tcp
ai-hackathon-web        Up (healthy)        0.0.0.0:80->80/tcp
```

#### 5.3 检查容器健康状态

```bash
# 等待几秒让健康检查完成
sleep 10

# 查看健康状态
docker ps --filter name=ai-hackathon --format "table {{.Names}}\t{{.Status}}"
```

**预期**: 两个容器都应该显示 `(healthy)`

#### 5.4 检查容器日志

```bash
# 查看后端日志
docker-compose logs server

# 查看前端日志
docker-compose logs web

# 实时查看所有日志
docker-compose logs -f
```

**预期**: 
- Server日志应该显示 "Starting server on :8080"
- Web日志应该显示Nginx启动信息
- 无错误信息

### 第六步：功能验证

#### 6.1 健康检查端点

```bash
# 后端健康检查
curl -v http://localhost:8080/api/v1/health

# 前端访问测试
curl -I http://localhost/
```

**预期输出**:
- 后端: `{"status":"healthy"}` (HTTP 200)
- 前端: HTTP 200 OK

#### 6.2 API功能测试

##### 测试1: 聊天接口

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "测试消息"}' \
  -v
```

**预期**: 
```json
{
  "message": "Message received",
  "echo": "测试消息"
}
```

##### 测试2: 视频播放接口

```bash
curl http://localhost:8080/api/v1/play/test123 -v
```

**预期**:
```json
{
  "message": "Playing video",
  "video_id": "test123"
}
```

##### 测试3: 文件上传接口

```bash
# 创建测试文件
echo "test content" > test.txt

# 上传文件
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@test.txt" \
  -v

# 清理测试文件
rm test.txt
```

**预期**: 返回上传成功信息

#### 6.3 前端访问测试

```bash
# 测试主页
curl http://localhost/ | head -20

# 测试API代理（通过前端）
curl http://localhost/api/v1/health
```

**预期**: 
- 主页返回HTML内容
- API代理正常工作

#### 6.4 静态资源测试

```bash
# 测试Gzip压缩
curl -H "Accept-Encoding: gzip" -I http://localhost/

# 测试缓存头
curl -I http://localhost/assets/index.js 2>/dev/null || echo "静态资源可能不存在"
```

### 第七步：网络验证

#### 7.1 容器网络检查

```bash
# 查看网络
docker network ls | grep ai-hackathon

# 检查网络详情
docker network inspect ai-hackathon_ai-hackathon-network
```

#### 7.2 容器间通信测试

```bash
# 从web容器访问server容器
docker-compose exec web wget -O- http://server:8080/api/v1/health

# 从server容器测试
docker-compose exec server wget -O- http://web/
```

**预期**: 容器间可以正常通信

### 第八步：数据持久化验证

#### 8.1 验证上传目录

```bash
# 检查上传目录是否存在
ls -la server/uploads/

# 创建测试文件
echo "test" > server/uploads/test-persistence.txt

# 重启容器
docker-compose restart server

# 检查文件是否还在
cat server/uploads/test-persistence.txt

# 清理测试文件
rm server/uploads/test-persistence.txt
```

**预期**: 重启后文件依然存在

### 第九步：工具验证

#### 9.1 Makefile命令测试

```bash
# 测试各个命令
make help
make status
make health
make logs-server
make logs-web
```

#### 9.2 部署脚本测试

```bash
# 测试各个命令
./deploy.sh help
./deploy.sh status
./deploy.sh health
```

### 第十步：性能验证

#### 10.1 镜像大小检查

```bash
docker images | grep ai-hackathon
```

**预期**:
- server镜像: 10-20MB
- web镜像: 20-30MB

#### 10.2 资源使用检查

```bash
# 查看资源使用
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

**预期**: 
- Server: CPU < 5%, Memory < 100MB（空闲时）
- Web: CPU < 1%, Memory < 50MB

#### 10.3 响应时间测试

```bash
# 测试响应时间
time curl -s http://localhost/api/v1/health > /dev/null
time curl -s http://localhost/ > /dev/null
```

**预期**: 响应时间 < 100ms

### 第十一步：清理验证

#### 11.1 测试停止功能

```bash
# 停止服务
make stop

# 验证容器已停止
docker-compose ps
```

**预期**: 所有容器状态为 Exit

#### 11.2 测试重启功能

```bash
# 重启服务
make start

# 等待服务就绪
sleep 5

# 验证服务恢复
make health
```

#### 11.3 测试清理功能（可选）

```bash
# 注意：这会删除所有容器和镜像
make clean

# 验证清理结果
docker-compose ps
docker images | grep ai-hackathon
```

## 📊 完整验证脚本

创建自动化验证脚本：

```bash
#!/bin/bash
# save as: verify-deployment.sh

set -e

echo "🔍 开始验证Docker部署方案..."

echo "✓ 步骤1: 检查文件完整性..."
test -f docker-compose.yml && echo "  ✓ docker-compose.yml"
test -f deploy.sh && echo "  ✓ deploy.sh"
test -f Makefile && echo "  ✓ Makefile"
test -f server/Dockerfile && echo "  ✓ server/Dockerfile"
test -f web/Dockerfile && echo "  ✓ web/Dockerfile"
test -f web/nginx.conf && echo "  ✓ web/nginx.conf"

echo "✓ 步骤2: 检查Docker环境..."
docker --version
docker-compose --version

echo "✓ 步骤3: 验证配置..."
docker-compose config > /dev/null && echo "  ✓ docker-compose.yml语法正确"

echo "✓ 步骤4: 构建镜像..."
docker-compose build

echo "✓ 步骤5: 启动服务..."
docker-compose up -d

echo "✓ 步骤6: 等待服务就绪..."
sleep 10

echo "✓ 步骤7: 健康检查..."
curl -f http://localhost:8080/api/v1/health
curl -f http://localhost/ > /dev/null

echo "✓ 步骤8: API测试..."
curl -f -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "test"}'

echo ""
echo "🎉 所有验证通过！"
echo "前端地址: http://localhost"
echo "后端地址: http://localhost:8080"
```

使用方法：
```bash
chmod +x verify-deployment.sh
./verify-deployment.sh
```

## 🐛 常见问题排查

### 问题1: 健康检查失败

**症状**: 容器启动后一直显示 `starting` 状态

**排查**:
```bash
# 查看详细日志
docker-compose logs server
docker-compose logs web

# 手动测试健康检查端点
curl http://localhost:8080/api/v1/health
```

### 问题2: 端口访问失败

**症状**: 无法访问 http://localhost

**排查**:
```bash
# 检查端口绑定
docker-compose ps

# 检查防火墙
sudo netstat -tlnp | grep 80
sudo netstat -tlnp | grep 8080
```

### 问题3: API代理不工作

**症状**: 前端无法访问后端API

**排查**:
```bash
# 进入web容器测试
docker-compose exec web sh
wget -O- http://server:8080/api/v1/health

# 检查nginx配置
docker-compose exec web cat /etc/nginx/conf.d/default.conf
```

### 问题4: 上传文件丢失

**症状**: 重启后上传的文件消失

**排查**:
```bash
# 检查volume映射
docker-compose config | grep volumes -A 5

# 检查主机目录
ls -la server/uploads/
```

## ✅ 验证通过标准

所有以下条件都满足才算验证通过：

- ✅ 所有文件都已创建
- ✅ Docker和Docker Compose版本符合要求
- ✅ 镜像构建成功，大小合理
- ✅ 容器启动成功，状态为healthy
- ✅ 健康检查端点返回200
- ✅ 所有API接口测试通过
- ✅ 前端页面可以正常访问
- ✅ 容器间网络通信正常
- ✅ 数据持久化工作正常
- ✅ 日志输出正常，无错误

## 📝 验证报告模板

```
Docker部署验证报告
==================

验证时间: ____________________
验证人员: ____________________

环境信息:
- Docker版本: ________________
- Docker Compose版本: ________
- 操作系统: __________________

验证结果:
[ ] 文件完整性检查
[ ] 环境检查
[ ] 配置验证
[ ] 构建验证
[ ] 部署验证
[ ] 功能验证
[ ] 网络验证
[ ] 数据持久化验证
[ ] 工具验证
[ ] 性能验证

问题记录:
_____________________________
_____________________________

总体评价: [ ] 通过  [ ] 未通过

备注:
_____________________________
_____________________________
```

---

**版本**: 1.0.0  
**最后更新**: 2025-10-24
