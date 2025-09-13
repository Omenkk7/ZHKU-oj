# 校园Java-OJ系统构建配置

.PHONY: help build clean test run-server run-judger run-worker docker-build docker-up docker-down

# 默认目标
help:
	@echo "校园Java-OJ系统构建命令："
	@echo "  build        - 构建所有服务"
	@echo "  clean        - 清理构建文件"
	@echo "  test         - 运行测试"
	@echo "  run-server   - 运行Web服务器"
	@echo "  run-judger   - 运行判题服务"
	@echo "  run-worker   - 运行异步任务处理器"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-up    - 启动Docker服务"
	@echo "  docker-down  - 停止Docker服务"

# 构建变量
BINARY_DIR := bin
SERVER_BINARY := $(BINARY_DIR)/server
JUDGER_BINARY := $(BINARY_DIR)/judger
WORKER_BINARY := $(BINARY_DIR)/worker

# Go构建参数
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LDFLAGS := -s -w

# 构建所有服务
build: clean
	@echo "构建Web服务器..."
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(SERVER_BINARY) ./cmd/server
	@echo "构建判题服务..."
	@go build -ldflags "$(LDFLAGS)" -o $(JUDGER_BINARY) ./cmd/judger
	@echo "构建异步任务处理器..."
	@go build -ldflags "$(LDFLAGS)" -o $(WORKER_BINARY) ./cmd/worker
	@echo "构建完成！"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf $(BINARY_DIR)
	@echo "清理完成！"

# 运行测试
test:
	@echo "运行单元测试..."
	@go test -v ./tests/unit/...
	@echo "运行集成测试..."
	@go test -v ./tests/integration/...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 运行Web服务器
run-server:
	@echo "启动Web服务器..."
	@CONFIG_PATH=configs/config.yaml go run cmd/server/main.go

# 运行判题服务
run-judger:
	@echo "启动判题服务..."
	@CONFIG_PATH=configs/config.yaml go run cmd/judger/main.go

# 运行异步任务处理器
run-worker:
	@echo "启动异步任务处理器..."
	@CONFIG_PATH=configs/config.yaml go run cmd/worker/main.go

# 代码格式化
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	@golangci-lint run

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod tidy
	@go mod download

# Docker相关命令
docker-build:
	@echo "构建Docker镜像..."
	@docker build -f deployments/docker/Dockerfile.web -t campus-oj-web:latest .
	@docker build -f deployments/docker/Dockerfile.judger -t campus-oj-judger:latest .

docker-up:
	@echo "启动Docker服务..."
	@cd deployments/docker && docker-compose up -d

docker-down:
	@echo "停止Docker服务..."
	@cd deployments/docker && docker-compose down

docker-logs:
	@echo "查看Docker服务日志..."
	@cd deployments/docker && docker-compose logs -f

# 开发环境快速启动
dev-setup: deps docker-up
	@echo "等待服务启动..."
	@sleep 30
	@echo "开发环境已准备就绪！"

# 生产环境部署
deploy-prod:
	@echo "部署到生产环境..."
	@cd deployments/docker && docker-compose -f docker-compose.prod.yml up -d

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	@go run cmd/migration/main.go

# 生成API文档
docs:
	@echo "生成API文档..."
	@swag init -g cmd/server/main.go -o ./docs

# 安全扫描
security-scan:
	@echo "运行安全扫描..."
	@gosec ./...

# 性能测试
benchmark:
	@echo "运行性能测试..."
	@go test -bench=. -benchmem ./...

# 完整的CI流程
ci: deps fmt lint test security-scan
	@echo "CI流程完成！"