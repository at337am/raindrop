# raindrop (rdrop)

[English](README.md) | [简体中文](README.zh-CN.md)

raindrop is an ultra-simple and easy-to-use command-line tool designed to help you instantly share files, text messages, or plain text content over your local network. It’s perfect for quick, temporary, or peer-to-peer data exchanges.

It supports several sharing modes, which you can combine freely:

* 📁 **File Sharing**: Share any single file, allowing recipients to download it directly via their browser.
* 📄 **Plain Text Sharing**: Display the content of text files like `.txt`, `.md`, or `.log` directly in the recipient’s browser without requiring a download.
* 💬 **Message Sending**: Send a short text message.

This tool is developed in Go, requiring zero configuration, lightweight and efficient, and runs cross-platform with just a single standalone binary.

## ✿ Why is it called raindrop?

The name is inspired by IU's song *Rain Drop*.  
Want to listen?

- [Spotify](https://open.spotify.com/track/6tlMVCqZlmxfnjZt3OiHjE)
- [171210 IU Palette Rain Drop (YouTube)](https://youtu.be/xgXFCOoQJVc)

Note: For convenience, the command-line tool is named `rdrop`. Please use `rdrop` to run the related commands.

## 📦 How to Install?

This project includes a [`justfile`](https://github.com/casey/just) to simplify common tasks.

### Method 1: Using just (if already installed)

```bash
just tidy     # tidy up dependencies
just build    # builds to ./release/rdrop
```

### Method 2: Using Go commands

```bash
go mod tidy
go build -o ./release/rdrop
```

## 🤯 How to Use?

When launched, raindrop runs a temporary HTTP server on your device (default port 1130). On any other device on the same local network—be it a phone, tablet, or computer—you just open the URL shown in the terminal in a browser to access the shared content.

### Option Overview

| Option | Description                                                             | Default |
| ------ | ----------------------------------------------------------------------- | ------- |
| `-i`   | Path to a single file to share (downloadable link)                      | None    |
| `-I`   | Path to a file whose content will be shown as plain text in the browser | None    |
| `-m`   | Text message to send                                                    | None    |
| `-p`   | Port to run the service on                                              | `1130`  |
| `-h`   | Show help information                                                   | —       |

### Usage Examples

Share just one file:

```bash
# Both commands do the same thing
rdrop -i my_document.pdf
rdrop my_document.pdf
```

Show a log file as plain text in the browser:

```bash
rdrop -I server-log-2025-06-27.log
```

Send a temporary message:

```bash
rdrop -m "Here’s a video link: https://youtu.be/xgXFCOoQJVc"
```

Combine file sharing with a message:

```bash
rdrop -i report.zip -m "Here’s today’s report"
```

Share a file, plain text content, and a message all at once:

```bash
rdrop -i design.sketch -I note.txt -m "Here’s the design draft and some notes"
```

Run the service on a different port:

```bash
rdrop -i my_video.mp4 -p 8080
```

## 🥺 Security Notice

raindrop is designed for **trusted local networks** like your home or office. It currently **does not provide any authentication or encryption**. Avoid using it over public or untrusted networks (such as coffee shops or airport Wi-Fi) when transferring sensitive data.
