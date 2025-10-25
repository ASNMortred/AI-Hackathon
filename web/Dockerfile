# 第一阶段：构建阶段
FROM node:18-alpine AS builder

# 配置Alpine国内镜像源（阿里云）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /build

# 配置npm国内镜像源（淘宝镜像）
RUN npm config set registry https://registry.npmmirror.com

# 优先复制依赖文件，利用Docker缓存
COPY package.json package-lock.json ./

# 安装依赖（使用npm ci以获得更可靠的构建）
RUN npm ci --only=production=false

# 复制源代码
COPY . .

# 构建生产版本
RUN npm run build

# 第二阶段：运行阶段
FROM nginx:alpine

# 配置Alpine国内镜像源（阿里云）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 复制自定义nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 从构建阶段复制构建产物
COPY --from=builder /build/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
