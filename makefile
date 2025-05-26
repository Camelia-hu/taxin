# 项目名称
PROJECT_NAME := taxin
# 可执行文件名称
BINARY_NAME := $(PROJECT_NAME)
# 目标端口
PORT := 50051
# 环境变量文件
ENV_FILE := $(CURDIR)/.env
# 测试包路径
TEST_PACKAGE := github.com/Camelia-hu/taxin/service
# Docker Compose 文件路径
DOCKER_COMPOSE_FILE := docker-compose.yaml

# 默认目标
all: build

# 清理生成的文件
clean:
	rm -f $(BINARY_NAME)

# 启动 Docker Compose 服务
start-docker:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# 停止 Docker Compose 服务
stop-docker:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# 构建项目
build: start-docker clean
	go build -o $(BINARY_NAME) ./cmd

# 运行项目
run: build
	@if [ -f $(ENV_FILE) ]; then \
		. $(ENV_FILE); \
	fi
	./$(BINARY_NAME)

# 部署项目（在后台运行）
deploy: build
	@if [ -f $(ENV_FILE) ]; then \
		. $(ENV_FILE); \
	fi
	nohup ./$(BINARY_NAME) > $(PROJECT_NAME).log 2>&1 &
	@echo "项目已部署，日志文件: $(PROJECT_NAME).log"

# 停止项目
stop: stop-docker
	@pkill -f $(BINARY_NAME) || true
	@echo "项目已停止"

## 运行单元测试
#test:
#	go test -v $(TEST_PACKAGE)

# 生成性能分析数据 需要下载Graphviz
pprof: build
	@echo "正在启动 pprof 服务，若程序已运行请忽略启动信息..."
	@if [ -f $(ENV_FILE) ]; then \
		. $(ENV_FILE); \
	fi
	@echo "正在通过 pprof 分析 CPU 性能..."
	@go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile
	@echo "若要分析内存性能，请在新终端执行：go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap"

# 显示帮助信息
help:
	@echo "可用的目标:"
	@echo "  all       - 构建项目"
	@echo "  clean     - 清理生成的文件"
	@echo "  start-docker - 启动 Docker Compose 服务"
	@echo "  stop-docker  - 停止 Docker Compose 服务"
	@echo "  build     - 启动 Docker Compose 并构建项目"
	@echo "  run       - 运行项目"
	@echo "  deploy    - 部署项目（在后台运行）"
	@echo "  stop      - 停止 Docker Compose 服务和项目"
	@echo "  test      - 运行单元测试"
	@echo "  pprof     - 生成性能分析数据"
	@echo "  help      - 显示帮助信息"

.PHONY: all clean start-docker stop-docker build run deploy stop test pprof help