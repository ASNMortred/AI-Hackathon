#!/bin/bash

# Docker国内镜像源配置脚本
# 用途：自动配置Docker守护进程使用国内镜像加速器

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检测操作系统
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    else
        OS="unknown"
    fi
    info "检测到操作系统: $OS"
}

# 备份现有配置
backup_config() {
    local config_file=$1
    if [ -f "$config_file" ]; then
        local backup_file="${config_file}.backup.$(date +%Y%m%d_%H%M%S)"
        cp "$config_file" "$backup_file"
        info "已备份现有配置到: $backup_file"
    fi
}

# 配置Linux系统的Docker镜像源
configure_linux() {
    step "配置Linux系统Docker镜像源..."
    
    local config_file="/etc/docker/daemon.json"
    
    # 检查是否有sudo权限
    if [ "$EUID" -ne 0 ]; then
        error "需要root权限，请使用sudo运行此脚本"
        exit 1
    fi
    
    # 创建配置目录
    mkdir -p /etc/docker
    
    # 备份现有配置
    backup_config "$config_file"
    
    # 写入新配置
    cat > "$config_file" << 'EOF'
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.ccs.tencentyun.com",
    "https://registry.docker-cn.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF
    
    info "✓ 配置文件已创建: $config_file"
    
    # 重启Docker服务
    step "重启Docker服务..."
    systemctl restart docker
    
    # 等待Docker启动
    sleep 3
    
    # 验证配置
    if docker info | grep -q "Registry Mirrors"; then
        info "✓ Docker镜像源配置成功"
        docker info | grep -A 10 "Registry Mirrors"
    else
        warn "配置可能未生效，请手动检查"
    fi
}

# 配置macOS系统的Docker镜像源
configure_macos() {
    step "配置macOS系统Docker镜像源..."
    
    local config_file="$HOME/.docker/daemon.json"
    
    # 创建配置目录
    mkdir -p "$HOME/.docker"
    
    # 备份现有配置
    backup_config "$config_file"
    
    # 写入新配置
    cat > "$config_file" << 'EOF'
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.ccs.tencentyun.com",
    "https://registry.docker-cn.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF
    
    info "✓ 配置文件已创建: $config_file"
    
    warn "macOS用户需要手动重启Docker Desktop："
    echo "  1. 打开Docker Desktop"
    echo "  2. 点击右上角齿轮图标（Settings）"
    echo "  3. 进入 Docker Engine 选项"
    echo "  4. 确认配置已生效"
    echo "  5. 点击 Apply & Restart"
    echo ""
    echo "或者通过命令行："
    echo "  osascript -e 'quit app \"Docker\"'"
    echo "  sleep 2"
    echo "  open -a Docker"
}

# 验证Docker配置
verify_docker() {
    step "验证Docker配置..."
    
    if ! command -v docker &> /dev/null; then
        error "Docker未安装"
        exit 1
    fi
    
    info "Docker版本: $(docker --version)"
    
    echo ""
    info "当前镜像源配置："
    docker info 2>/dev/null | grep -A 10 "Registry Mirrors" || warn "无法获取镜像源信息"
}

# 配置Go代理
configure_go() {
    step "配置Go模块代理..."
    
    local shell_rc=""
    if [ -f "$HOME/.bashrc" ]; then
        shell_rc="$HOME/.bashrc"
    elif [ -f "$HOME/.zshrc" ]; then
        shell_rc="$HOME/.zshrc"
    fi
    
    if [ -n "$shell_rc" ]; then
        # 检查是否已配置
        if grep -q "GOPROXY" "$shell_rc"; then
            warn "Go代理已在 $shell_rc 中配置"
        else
            cat >> "$shell_rc" << 'EOF'

# Go模块代理配置（由setup-china-mirrors.sh添加）
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
export GOSUMDB=sum.golang.google.cn
EOF
            info "✓ Go代理配置已添加到: $shell_rc"
            info "  请运行: source $shell_rc"
        fi
    fi
    
    # 使用go env配置（如果go已安装）
    if command -v go &> /dev/null; then
        go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
        go env -w GOSUMDB=sum.golang.google.cn
        info "✓ Go环境变量已配置"
        info "  GOPROXY=$(go env GOPROXY)"
    fi
}

# 配置NPM镜像
configure_npm() {
    step "配置NPM镜像源..."
    
    if command -v npm &> /dev/null; then
        # 备份当前配置
        local current_registry=$(npm config get registry)
        info "当前NPM源: $current_registry"
        
        # 设置淘宝镜像
        npm config set registry https://registry.npmmirror.com
        
        info "✓ NPM镜像源已配置"
        info "  新的NPM源: $(npm config get registry)"
        
        # 提示安装nrm
        echo ""
        info "推荐安装nrm来管理NPM源："
        echo "  npm install -g nrm"
        echo "  nrm ls"
        echo "  nrm use taobao"
    else
        warn "NPM未安装，跳过配置"
    fi
}

# 测试镜像源速度
test_mirrors() {
    step "测试镜像源速度..."
    
    echo ""
    info "测试Docker Hub镜像..."
    for mirror in \
        "https://docker.mirrors.ustc.edu.cn" \
        "https://hub-mirror.c.163.com" \
        "https://mirror.ccs.tencentyun.com"
    do
        echo -n "  $mirror: "
        time_result=$(curl -o /dev/null -s -w '%{time_total}\n' "$mirror" 2>/dev/null || echo "失败")
        echo "${time_result}s"
    done
    
    echo ""
    info "测试Go代理..."
    for proxy in \
        "https://goproxy.cn" \
        "https://mirrors.aliyun.com/goproxy/" \
        "https://goproxy.io"
    do
        echo -n "  $proxy: "
        time_result=$(curl -o /dev/null -s -w '%{time_total}\n' "$proxy" 2>/dev/null || echo "失败")
        echo "${time_result}s"
    done
    
    echo ""
    info "测试NPM镜像..."
    for npm_mirror in \
        "https://registry.npmmirror.com" \
        "https://registry.npmjs.org"
    do
        echo -n "  $npm_mirror: "
        time_result=$(curl -o /dev/null -s -w '%{time_total}\n' "$npm_mirror" 2>/dev/null || echo "失败")
        echo "${time_result}s"
    done
}

# 显示帮助
show_help() {
    cat << EOF
Docker国内镜像源配置脚本

用法: $0 [选项]

选项:
    --docker        仅配置Docker镜像源
    --go            仅配置Go代理
    --npm           仅配置NPM镜像
    --test          测试镜像源速度
    --all           配置所有（默认）
    --help          显示此帮助信息

示例:
    $0              # 配置所有镜像源
    $0 --docker     # 仅配置Docker
    $0 --test       # 测试速度

EOF
}

# 主函数
main() {
    echo "======================================"
    echo "  Docker国内镜像源配置工具"
    echo "======================================"
    echo ""
    
    detect_os
    
    # 解析参数
    local config_docker=false
    local config_go=false
    local config_npm=false
    local run_test=false
    
    if [ $# -eq 0 ]; then
        # 默认配置所有
        config_docker=true
        config_go=true
        config_npm=true
    else
        while [[ $# -gt 0 ]]; do
            case $1 in
                --docker)
                    config_docker=true
                    shift
                    ;;
                --go)
                    config_go=true
                    shift
                    ;;
                --npm)
                    config_npm=true
                    shift
                    ;;
                --test)
                    run_test=true
                    shift
                    ;;
                --all)
                    config_docker=true
                    config_go=true
                    config_npm=true
                    shift
                    ;;
                --help)
                    show_help
                    exit 0
                    ;;
                *)
                    error "未知选项: $1"
                    show_help
                    exit 1
                    ;;
            esac
        done
    fi
    
    # 执行配置
    if [ "$config_docker" = true ]; then
        if [ "$OS" = "linux" ]; then
            configure_linux
        elif [ "$OS" = "macos" ]; then
            configure_macos
        else
            error "不支持的操作系统"
            exit 1
        fi
        verify_docker
    fi
    
    if [ "$config_go" = true ]; then
        configure_go
    fi
    
    if [ "$config_npm" = true ]; then
        configure_npm
    fi
    
    if [ "$run_test" = true ]; then
        test_mirrors
    fi
    
    echo ""
    echo "======================================"
    info "✓ 配置完成！"
    echo "======================================"
    echo ""
    info "后续步骤："
    echo "  1. 如果配置了Go代理，运行: source ~/.bashrc (或 ~/.zshrc)"
    echo "  2. 如果是macOS，需要重启Docker Desktop"
    echo "  3. 运行: docker-compose build 测试镜像拉取速度"
    echo ""
    info "查看详细文档: CHINA_MIRROR_CONFIG.md"
}

main "$@"
