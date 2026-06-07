# ============================================================
# 校园协作平台 - 根目录 Makefile
# ============================================================

.PHONY: help docker-up docker-down db-migrate server-run web-run test

help: ## 显示所有命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ===== Docker =====
docker-up: ## 启动 MySQL + Redis 容器
	docker compose up -d

docker-down: ## 停止所有容器
	docker compose down

docker-logs: ## 查看容器日志
	docker compose logs -f

docker-status: ## 查看容器状态
	docker compose ps

# ===== 后端 =====
server-run: ## 启动后端开发服务器
	cd server; go run cmd/main.go

server-build: ## 编译后端
	cd server; go build -o bin/server cmd/main.go

server-test: ## 运行后端测试
	cd server; go test -v ./internal/...

server-test-engine: ## 仅运行引擎单元测试
	cd server; go test -v -cover ./internal/engine/...

server-lint: ## 后端代码检查
	cd server; go vet ./...

server-fmt: ## 后端代码格式化
	cd server; go fmt ./...

# ===== 前端 =====
web-dev: ## 启动前端开发服务器
	cd web; pnpm dev

web-build: ## 前端生产构建
	cd web; pnpm build

web-lint: ## 前端代码检查
	cd web; pnpm lint

# ===== 数据库 =====
db-migrate-up: ## 执行所有数据库迁移
	cd server; migrate -path migrations -database "mysql://root:123456@tcp(localhost:3306)/campus_collab" up

db-migrate-down: ## 回滚最近一次迁移
	cd server; migrate -path migrations -database "mysql://root:123456@tcp(localhost:3306)/campus_collab" down 1

db-migrate-status: ## 查看迁移状态
	cd server; migrate -path migrations -database "mysql://root:123456@tcp(localhost:3306)/campus_collab" version

# ===== 工具 =====
seed: ## 填充测试种子数据
	cd server; go run scripts/seed.go
