# raindrop (rdrop)

[English](README.md) | [简体中文](README.zh-CN.md)

raindrop 是一个极其简单易用的命令行工具，旨在帮助你在局域网内即时地共享文件、文本消息或纯文本内容，适用于临时性或点对点的数据交换场景。

支持以下分享模式，且这些模式可以灵活组合使用：

- 📁 **文件分享**：共享任意单个文件，接收方可以直接通过浏览器下载。
- 📄 **纯文本分享**：可以将 `.txt`、`.md` 或 `.log` 等文本文件的内容直接在接收方的浏览器中显示，无需下载。
- 💬 **消息发送**：发送一段简短的文本消息。

该工具基于 Go 语言开发，零配置、轻量高效，且跨平台运行，仅需一个独立的二进制文件即可使用。

## ✿ 为什么叫 raindrop？

灵感源自 IU 的歌曲《Rain Drop》  
想听听吗？

- [Spotify](https://open.spotify.com/track/6tlMVCqZlmxfnjZt3OiHjE)
- [171210 IU Palette Rain Drop (YouTube)](https://youtu.be/xgXFCOoQJVc)

注意：为了方便使用，命令行工具的名称为 `rdrop`，您需要使用 `rdrop` 来执行相关命令。

## 📦 如何安装？

本项目提供了 [`justfile`](https://github.com/casey/just)，用于简化常见操作。

### 方式一：使用 just（如已安装）

```bash
just tidy     # 整理依赖
just build    # 构建至 ./release/rdrop
```

### 方式二：使用 Go 命令

```bash
go mod tidy
go build -o ./release/rdrop
```

## 🤯 如何使用？

raindrop 启动后，会在你的设备上运行一个临时的 HTTP 服务器（默认监听端口为 1130）。你只需在局域网内的任何其他设备（手机、平板、电脑）的浏览器中打开终端上显示的 URL 即可访问共享的内容。

### 选项说明

| 选项   | 说明                       | 默认值    |
| :--- | :----------------------- | :----- |
| `-i` | 要共享的单个文件路径，接收方将看到一个下载链接  | 无      |
| `-I` | 文件路径，其内容将作为纯文本直接在浏览器中显示 | 无      |
| `-m` | 要发送的文本消息内容              | 无      |
| `-p` | 服务运行的端口                | `1130` |
| `-h` | 显示帮助信息                  | —      |

### 用法示例

仅共享一个文件：

```bash
# 以下两种方式等效
rdrop -i my_document.pdf
rdrop my_document.pdf
```

将日志文件内容显示为纯文本：

```bash
rdrop -I server-log-2025-06-27.log
```

发送一条临时消息：

```bash
rdrop -m "这是视频链接: https://youtu.be/xgXFCOoQJVc"
```

组合使用，分享文件并附上说明：

```bash
rdrop -i report.zip -m "这是今天的报告"
```

组合使用，同时分享文件、纯文本内容和消息：

```bash
rdrop -i design.sketch -I note.txt -m "这是设计稿和一些注意事项"
```

指定一个不同的端口：

```bash
rdrop -i my_video.mp4 -p 8080
```

## 🥺 安全提示

raindrop 专为**可信的局域网**（如家庭、办公室网络）环境设计。它目前**不包含任何身份验证或加密机制**。请勿在公共或不受信任的网络（如咖啡店、机场 Wi-Fi）上使用它来传输敏感数据。
