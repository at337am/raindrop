# 默认目标：列出所有可用任务
default:
    @just --list

# 整理 go.mod 与 go.sum，清理无效依赖
tidy:
    go mod tidy

# 清除构建目录（release/）
clean:
    @echo "正在清理..."
    rm -rfv release/
    @echo "清理完成"

# 构建二进制文件，输出至 release/
build:
    @mkdir -p release/
    go build -o ./release/rdrop ./cmd/rdrop/

# 使用已构建的二进制运行（适用于测试或生产）
run: build
    ./release/rdrop -i ./README.md -m "Hello World"

# 直接运行源码（适用于开发调试）
run-dev:
    go run ./cmd/rdrop/main.go -i ./README.md -m "Hello World"

# 安装至 $GOBIN 或 $GOPATH/bin
install:
    go install ./cmd/rdrop/
    @echo "安装完成"
