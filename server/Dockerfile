# 第一阶段：构建阶段
FROM golang:1.20-alpine AS builder

# 配置Alpine国内镜像源（阿里云）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /build

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 配置Go模块代理（使用阿里云和七牛云镜像）
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct \
    GOSUMDB=sum.golang.google.cn

# 优先复制依赖文件，利用Docker缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用程序
# CGO_ENABLED=0 生成静态二进制文件
# -ldflags="-w -s" 减小二进制文件大小
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/bin/server \
    ./cmd/server/main.go

# 第二阶段：运行阶段
FROM alpine:latest

# 配置Alpine国内镜像源（阿里云）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /app

# 安装 wget 用于健康检查和 CA 证书
RUN apk add --no-cache wget ca-certificates tzdata

# 从构建器复制时区信息
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 设置时区
ENV TZ=Asia/Shanghai

# 复制编译好的二进制文件
COPY --from=builder /build/bin/server /app/server

# 复制配置文件
COPY configs/config.yaml /app/configs/config.yaml

# 创建上传目录
RUN mkdir -p /app/uploads

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查 - 由于没有 /health 端点,暂时禁用
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# 运行应用
CMD ["/app/server"]
