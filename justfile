# 列出所有可用任务
default:
    @just --list

cp:
    fd -e go -x sh -c 'echo "===== {} ====="; cat {}; echo' | wl-copy

# 整理依赖
tidy:
    @go mod tidy
    @echo "依赖整理完成"

# 清除构建目录
clean:
    @echo "正在清理..."
    @rm -rfv release/
    @echo "清理完成"

# 构建二进制文件
build:
    @mkdir -p release/
    @go build -o ./release/rdrop
    @echo "构建完成 -> ./release/rdrop"
