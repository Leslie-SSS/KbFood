.PHONY: help build run test clean docker sqlc test-unit test-integration test-e2e test-all

# 默认目标
help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make run             - Run the application"
	@echo "  make test            - Run unit tests"
	@echo "  make test-unit       - Run unit tests only"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-e2e        - Run E2E tests"
	@echo "  make test-all        - Run all tests (unit + integration + e2e)"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make docker          - Build Docker image"
	@echo "  make sqlc            - Generate sqlc code"

# 构建应用
build:
	@echo "Building..."
	@go build -o bin/main ./cmd/server
	@echo "Build complete: bin/main"

# 运行应用
run:
	@echo "Running..."
	@go run ./cmd/server/main.go

# 运行单元测试
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.out -tags='!integration !e2e' ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 运行集成测试 (需要数据库)
test-integration:
	@echo "Running integration tests..."
	@go test -v -race -tags=integration ./test/integration/...
	@echo "Integration tests complete"

# 运行 E2E 测试 (需要完整环境)
test-e2e:
	@echo "Running E2E tests..."
	@go test -v -race -tags=e2e ./test/e2e/...
	@echo "E2E tests complete"

# 运行所有测试
test-all:
	@echo "Running all tests..."
	@$(MAKE) test-unit
	@$(MAKE) test-integration
	@$(MAKE) test-e2e
	@echo "All tests complete"

# 运行测试 (默认单元测试)
test: test-unit

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# 生成 sqlc 代码
sqlc:
	@echo "Generating sqlc code..."
	@sqlc generate
	@echo "sqlc generation complete"

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# 代码检查
lint:
	@echo "Running linters..."
	@golangci-lint run ./...

# 构建 Docker 镜像
docker:
	@echo "Building Docker image..."
	@docker build -t kbfood:latest -f deployments/Dockerfile .
	@echo "Docker build complete"

# 依赖管理
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Wire 依赖注入生成
wire:
	@echo "Generating wire code..."
	@cd cmd/server && wire
	@echo "Wire generation complete"
