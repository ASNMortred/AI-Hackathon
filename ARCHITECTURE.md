# 🏗️ Docker化架构设计

## 系统架构图

```mermaid
graph TB
    User[用户浏览器] --> Web[Web容器<br/>Nginx:Alpine<br/>Port 80]
    Web --> Server[Server容器<br/>Go 1.24.5<br/>Port 8080]
    Server --> Uploads[(上传文件<br/>Volume挂载)]
    
    Web -.健康检查.-> HealthWeb[/api/v1/health]
    Server -.健康检查.-> HealthServer[/api/v1/health]
    
    Compose[Docker Compose] --> Web
    Compose --> Server
    Compose --> Network[ai-hackathon-network]
    
    Web --> Network
    Server --> Network
```

## 容器架构

### 前端容器（web）

```mermaid
graph LR
    A[node:18-alpine<br/>构建阶段] -->|npm ci| B[安装依赖]
    B -->|npm run build| C[构建产物]
    C -->|复制到| D[nginx:alpine<br/>运行阶段]
    D --> E[静态文件服务<br/>Port 80]
    D --> F[API反向代理<br/>→ server:8080]
```

**特性**:
- 📦 多阶段构建，最终镜像仅约20MB
- 🚀 Nginx高性能静态文件服务
- 🔄 自动API反向代理
- 💾 Gzip压缩 + 静态资源缓存
- ♻️ SPA路由支持

### 后端容器（server）

```mermaid
graph LR
    A[golang:1.24.5-alpine<br/>构建阶段] -->|go mod download| B[下载依赖]
    B -->|CGO_ENABLED=0| C[静态编译]
    C -->|复制到| D[alpine:latest<br/>运行阶段]
    D --> E[二进制文件<br/>Port 8080]
    D --> F[健康检查<br/>/api/v1/health]
```

**特性**:
- 🔒 静态编译，无外部依赖
- 📦 最终镜像仅约15MB
- ⚡ Alpine Linux极致精简
- 🏥 内置健康检查
- 📁 持久化上传文件

## 网络架构

```mermaid
graph TB
    Internet[互联网] -->|80| HostNginx[主机:80端口]
    Internet -->|8080| HostServer[主机:8080端口]
    
    subgraph Docker网络
        HostNginx --> WebContainer[web容器:80]
        HostServer --> ServerContainer[server容器:8080]
        
        WebContainer -->|内部网络| ServerContainer
    end
    
    ServerContainer --> Volume[持久化卷<br/>./server/uploads]
```

## 构建流程

### 前端构建流程

```mermaid
graph TD
    Start[开始构建] --> Cache1{package.json<br/>缓存命中?}
    Cache1 -->|是| Skip1[跳过npm ci]
    Cache1 -->|否| Install[npm ci安装依赖]
    
    Install --> Cache2{源代码<br/>缓存命中?}
    Skip1 --> Cache2
    Cache2 -->|是| Skip2[跳过构建]
    Cache2 -->|否| Build[npm run build]
    
    Build --> Copy[复制到nginx镜像]
    Skip2 --> Copy
    Copy --> Config[配置nginx.conf]
    Config --> End[构建完成]
```

### 后端构建流程

```mermaid
graph TD
    Start[开始构建] --> Cache1{go.mod/go.sum<br/>缓存命中?}
    Cache1 -->|是| Skip1[跳过下载]
    Cache1 -->|否| Download[go mod download]
    
    Download --> Cache2{源代码<br/>缓存命中?}
    Skip1 --> Cache2
    Cache2 -->|是| Skip2[跳过编译]
    Cache2 -->|否| Compile[静态编译]
    
    Compile --> Copy[复制到alpine镜像]
    Skip2 --> Copy
    Copy --> Config[复制配置文件]
    Config --> End[构建完成]
```

## 部署流程

```mermaid
graph TD
    Start[开始部署] --> Init[初始化环境]
    Init --> CheckDocker{检查Docker}
    CheckDocker -->|未安装| Error1[安装Docker]
    CheckDocker -->|已安装| CreateDir[创建目录]
    
    CreateDir --> Build[构建镜像]
    Build --> BuildServer[构建server镜像]
    Build --> BuildWeb[构建web镜像]
    
    BuildServer --> Start[启动服务]
    BuildWeb --> Start
    
    Start --> StartServer[启动server容器]
    StartServer --> HealthServer{server健康检查}
    HealthServer -->|失败| Retry1[重试]
    HealthServer -->|成功| StartWeb[启动web容器]
    
    StartWeb --> HealthWeb{web健康检查}
    HealthWeb -->|失败| Retry2[重试]
    HealthWeb -->|成功| Complete[部署完成]
    
    Retry1 --> HealthServer
    Retry2 --> HealthWeb
```

## 数据流

### 用户请求流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant W as Web容器<br/>(Nginx)
    participant S as Server容器<br/>(Gin)
    participant V as 文件存储<br/>(Volume)
    
    U->>W: 访问页面
    W->>U: 返回HTML/JS/CSS
    
    U->>W: API请求 /api/v1/*
    W->>S: 反向代理到 :8080
    S->>S: 处理业务逻辑
    
    alt 文件上传
        S->>V: 保存文件
        V->>S: 返回路径
    end
    
    S->>W: 返回响应
    W->>U: 返回数据
```

### 健康检查流程

```mermaid
sequenceDiagram
    participant DC as Docker Compose
    participant S as Server容器
    participant W as Web容器
    
    loop 每30秒
        DC->>S: wget /api/v1/health
        S->>DC: 200 OK
        DC->>DC: 标记server为healthy
        
        DC->>W: wget /
        W->>DC: 200 OK
        DC->>DC: 标记web为healthy
    end
    
    Note over DC: depends_on条件满足<br/>web等待server健康
```

## 文件结构映射

### 构建时文件映射

```
主机                                容器内部
────────────────────────          ────────────────────────
server/
├── go.mod                  →    /build/go.mod
├── go.sum                  →    /build/go.sum
├── cmd/                    →    /build/cmd/
├── internal/               →    /build/internal/
└── configs/                →    /build/configs/

↓ 编译后 ↓

server/                          /app/
└── [构建产物]              →    ├── server (二进制)
                                 └── configs/config.yaml
```

### 运行时文件映射

```
主机                                容器内部
────────────────────────          ────────────────────────
server/
├── uploads/                ⇄    /app/uploads/
└── configs/config.yaml     →    /app/configs/config.yaml

web/
└── dist/                   →    /usr/share/nginx/html/
```

## 资源分配

### 推荐配置

| 容器 | CPU | 内存 | 磁盘 |
|------|-----|------|------|
| server | 0.5-1核 | 256-512MB | 100MB |
| web | 0.25-0.5核 | 128-256MB | 50MB |

### 可配置限制

```yaml
deploy:
  resources:
    limits:
      cpus: '1'
      memory: 512M
    reservations:
      cpus: '0.5'
      memory: 256M
```

## 安全架构

```mermaid
graph TB
    subgraph 容器隔离
        Server[Server容器<br/>独立网络命名空间]
        Web[Web容器<br/>独立网络命名空间]
    end
    
    subgraph 安全特性
        Static[静态编译<br/>无动态库依赖]
        Minimal[最小镜像<br/>减少攻击面]
        Readonly[只读根文件系统<br/>可选]
        NonRoot[非Root用户<br/>可配置]
    end
    
    Server --> Static
    Server --> Minimal
    Web --> Minimal
    
    subgraph 网络安全
        Firewall[防火墙规则]
        SSL[SSL/TLS<br/>可配置]
    end
```

## 扩展架构

### 水平扩展

```mermaid
graph TB
    LB[负载均衡器<br/>Nginx/HAProxy] --> W1[Web实例1]
    LB --> W2[Web实例2]
    LB --> W3[Web实例3]
    
    W1 --> S1[Server实例1]
    W2 --> S2[Server实例2]
    W3 --> S3[Server实例3]
    
    S1 --> DB[(共享数据库)]
    S2 --> DB
    S3 --> DB
    
    S1 --> Cache[(Redis缓存)]
    S2 --> Cache
    S3 --> Cache
```

### 微服务扩展

```mermaid
graph TB
    Gateway[API Gateway] --> Auth[认证服务]
    Gateway --> Upload[上传服务]
    Gateway --> Play[播放服务]
    Gateway --> Chat[聊天服务]
    
    Upload --> Storage[(对象存储)]
    Play --> CDN[CDN]
    Chat --> MQ[消息队列]
```

## 监控架构

```mermaid
graph TB
    subgraph 应用层
        Server[Server容器]
        Web[Web容器]
    end
    
    subgraph 监控层
        Server --> Logs[日志收集<br/>JSON格式]
        Web --> Logs
        
        Server --> Metrics[指标收集<br/>Prometheus]
        Web --> Metrics
        
        Server --> Health[健康检查<br/>Docker内置]
        Web --> Health
    end
    
    subgraph 告警层
        Logs --> Alert[告警系统]
        Metrics --> Alert
        Health --> Alert
    end
```

## 备份架构

```mermaid
graph LR
    App[应用运行] --> Upload[上传文件]
    App --> Config[配置文件]
    
    Upload --> Volume[Docker Volume]
    Config --> HostPath[主机路径]
    
    Volume --> Backup1[定期备份]
    HostPath --> Backup2[定期备份]
    
    Backup1 --> Storage[(备份存储)]
    Backup2 --> Storage
```

## 总结

本Docker化架构具有以下优势：

### ✅ 性能优势
- 多阶段构建，镜像小巧（总计<50MB）
- 构建缓存优化，加速开发迭代
- Nginx高性能静态文件服务
- Alpine Linux极致精简

### ✅ 可靠性优势
- 健康检查确保服务可用
- 自动重启策略
- 依赖隔离，版本锁定
- 容器化环境一致性

### ✅ 安全性优势
- 静态编译，无外部依赖
- 最小基础镜像
- 容器隔离
- 可配置安全加固

### ✅ 运维优势
- 一键部署
- 标准化配置
- 日志集中管理
- 易于扩展

---

**版本**: 1.0.0  
**最后更新**: 2025-10-24
