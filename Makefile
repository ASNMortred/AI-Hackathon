.PHONY: help init build start stop restart logs status health clean backup deploy

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

## 显示帮助信息
help:
	@echo ''
	@echo '用法:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo '目标:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "  ${YELLOW}%-15s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${GREEN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## 初始化环境
init:
	@echo "检查系统要求..."
	@command -v docker >/dev/null 2>&1 || { echo "需要安装Docker"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "需要安装Docker Compose"; exit 1; }
	@echo "创建必要的目录..."
	@mkdir -p server/uploads
	@chmod 777 server/uploads
	@if [ ! -f .env ]; then \
		echo "创建.env文件..."; \
		cp .env.example .env 2>/dev/null || echo "警告: .env.example不存在"; \
	fi
	@echo "✓ 初始化完成"

## 构建Docker镜像
build:
	@echo "构建Docker镜像..."
	@docker-compose build --no-cache
	@echo "✓ 构建完成"

## 启动所有服务
start:
	@echo "启动服务..."
	@docker-compose up -d
	@echo "✓ 服务已启动"
	@sleep 3
	@make health

## 停止所有服务
stop:
	@echo "停止服务..."
	@docker-compose down
	@echo "✓ 服务已停止"

## 重启所有服务
restart:
	@echo "重启服务..."
	@docker-compose restart
	@echo "✓ 服务已重启"

## 查看服务日志
logs:
	@docker-compose logs -f

## 查看后端日志
logs-server:
	@docker-compose logs -f server

## 查看前端日志
logs-web:
	@docker-compose logs -f web

## 查看服务状态
status:
	@docker-compose ps

## 健康检查
health:
	@echo "检查服务健康状态..."
	@curl -sf http://localhost:8080/api/v1/health > /dev/null && echo "✓ 后端服务健康" || echo "✗ 后端服务未就绪"
	@curl -sf http://localhost > /dev/null && echo "✓ 前端服务健康" || echo "✗ 前端服务未就绪"

## 清理容器和镜像
clean:
	@echo "清理Docker资源..."
	@docker-compose down -v
	@docker rmi ai-hackathon_server:latest ai-hackathon_web:latest 2>/dev/null || true
	@echo "✓ 清理完成"

## 备份数据
backup:
	@echo "开始备份..."
	@mkdir -p backups
	@tar -czf backups/uploads_$$(date +%Y%m%d_%H%M%S).tar.gz server/uploads 2>/dev/null || true
	@tar -czf backups/configs_$$(date +%Y%m%d_%H%M%S).tar.gz server/configs
	@echo "✓ 备份完成，文件保存在backups目录"

## 完整部署（init + build + start）
deploy: init build start
	@echo "========================================="
	@echo "部署完成！"
	@echo "前端地址: http://localhost"
	@echo "后端API: http://localhost:8080/api/v1"
	@echo "========================================="

## 进入后端容器
shell-server:
	@docker-compose exec server sh

## 进入前端容器
shell-web:
	@docker-compose exec web sh

## 查看容器资源使用
stats:
	@docker stats --no-stream

## 重新构建并启动
rebuild: clean build start

## 开发模式（实时日志）
dev:
	@docker-compose up

## 拉取最新代码并重新部署
update:
	@git pull
	@make rebuild
	@echo "✓ 更新完成"
