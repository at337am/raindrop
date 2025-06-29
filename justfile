# 列出所有可用任务
default:
    @just --list

# 整理依赖
tidy:
    @go mod tidy
    @echo "依赖整理完成"

# 清除构建目录
clean:
    @echo "正在清理..."
    @rm -rfv release/
    @echo "清理完成"

# 直接运行源码
run:
    @go run main.go -i ./README.md -m "Hello World"

# 构建二进制文件
build:
    @mkdir -p release/
    @go build -o ./release/rdrop
    @echo "构建完成 -> ./release/rdrop"
