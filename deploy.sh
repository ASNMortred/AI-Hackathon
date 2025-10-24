#!/bin/bash

# AI-Hackathon Docker 部署脚本
# 用途：简化Docker Compose操作

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印信息
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker和Docker Compose
check_requirements() {
    info "检查系统要求..."
    
    if ! command -v docker &> /dev/null; then
        error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    info "✓ Docker版本: $(docker --version)"
    info "✓ Docker Compose版本: $(docker-compose --version)"
}

# 创建必要的目录
create_directories() {
    info "创建必要的目录..."
    mkdir -p server/uploads
    chmod 777 server/uploads
    info "✓ 目录创建完成"
}

# 构建镜像
build() {
    info "开始构建Docker镜像..."
    docker-compose build --no-cache
    info "✓ 镜像构建完成"
}

# 启动服务
start() {
    info "启动服务..."
    docker-compose up -d
    info "✓ 服务启动完成"
    
    # 等待服务就绪
    info "等待服务就绪..."
    sleep 5
    
    # 检查健康状态
    check_health
}

# 停止服务
stop() {
    info "停止服务..."
    docker-compose down
    info "✓ 服务已停止"
}

# 重启服务
restart() {
    info "重启服务..."
    docker-compose restart
    info "✓ 服务已重启"
}

# 查看日志
logs() {
    docker-compose logs -f
}

# 查看状态
status() {
    info "服务状态:"
    docker-compose ps
}

# 健康检查
check_health() {
    info "检查服务健康状态..."
    
    # 检查后端
    if curl -sf http://localhost:8080/api/v1/health > /dev/null; then
        info "✓ 后端服务健康"
    else
        warn "✗ 后端服务未就绪"
    fi
    
    # 检查前端
    if curl -sf http://localhost > /dev/null; then
        info "✓ 前端服务健康"
    else
        warn "✗ 前端服务未就绪"
    fi
}

# 清理所有
clean() {
    warn "警告：此操作将删除所有容器、镜像和数据卷"
    read -p "确认继续？(y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        info "清理中..."
        docker-compose down -v
        docker rmi ai-hackathon-server:latest ai-hackathon-web:latest 2>/dev/null || true
        info "✓ 清理完成"
    else
        info "已取消"
    fi
}

# 备份数据
backup() {
    info "开始备份..."
    BACKUP_DIR="backups"
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    
    mkdir -p $BACKUP_DIR
    
    # 备份上传文件
    if [ -d "server/uploads" ]; then
        tar -czf "$BACKUP_DIR/uploads_$TIMESTAMP.tar.gz" server/uploads
        info "✓ 上传文件已备份到 $BACKUP_DIR/uploads_$TIMESTAMP.tar.gz"
    fi
    
    # 备份配置文件
    tar -czf "$BACKUP_DIR/configs_$TIMESTAMP.tar.gz" server/configs
    info "✓ 配置文件已备份到 $BACKUP_DIR/configs_$TIMESTAMP.tar.gz"
}

# 显示帮助信息
show_help() {
    cat << EOF
AI-Hackathon Docker 部署脚本

用法: $0 [命令]

命令:
    init        初始化环境（检查依赖、创建目录）
    build       构建Docker镜像
    start       启动所有服务
    stop        停止所有服务
    restart     重启所有服务
    logs        查看服务日志
    status      查看服务状态
    health      检查服务健康状态
    clean       清理所有容器和镜像
    backup      备份数据和配置
    deploy      完整部署（init + build + start）
    help        显示此帮助信息

示例:
    $0 deploy     # 完整部署
    $0 logs       # 查看日志
    $0 health     # 健康检查

EOF
}

# 初始化环境
init() {
    check_requirements
    create_directories
    
    # 检查.env文件
    if [ ! -f ".env" ]; then
        warn ".env文件不存在"
        if [ -f ".env.example" ]; then
            info "从.env.example创建.env文件"
            cp .env.example .env
            warn "请编辑.env文件配置必要的环境变量"
        fi
    fi
    
    info "✓ 初始化完成"
}

# 完整部署
deploy() {
    info "开始完整部署..."
    init
    build
    start
    info "========================================="
    info "部署完成！"
    info "前端地址: http://localhost"
    info "后端API: http://localhost:8080/api/v1"
    info "========================================="
}

# 主逻辑
main() {
    case "${1:-help}" in
        init)
            init
            ;;
        build)
            build
            ;;
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        logs)
            logs
            ;;
        status)
            status
            ;;
        health)
            check_health
            ;;
        clean)
            clean
            ;;
        backup)
            backup
            ;;
        deploy)
            deploy
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
