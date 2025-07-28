# raindrop (rdrop)

[English](README.md) | [简体中文](README.zh-CN.md)

**raindrop** is an extremely simple and user-friendly command-line tool designed to help you instantly share files, text messages, or plain content within a local area network (LAN). It's ideal for temporary or point-to-point data exchange scenarios.

The following sharing modes are supported and can be flexibly combined:

* 📁 **File Sharing**: Share any number of files that recipients can directly download via their browser.
* 📄 **Plain Text Sharing**: Display the contents of `.txt`, `.md`, `.log`, and other text files directly in the browser without downloading.
* 💬 **Message Sending**: Send a short text message.

This tool is developed in Go, zero-config, lightweight, and cross-platform. It runs with just a standalone binary.

## ✿ Why is it called *raindrop*?

The name is inspired by IU's song **"Rain Drop"**.
Want to give it a listen?

* [Spotify](https://open.spotify.com/track/6tlMVCqZlmxfnjZt3OiHjE)
* [171210 IU Palette Rain Drop (YouTube)](https://youtu.be/xgXFCOoQJVc)

Note: For convenience, the command-line tool is named `rdrop`. Use `rdrop` to execute related commands.

## 📦 How to Install?

This project provides a [`justfile`](https://github.com/casey/just) to simplify common operations.

### Method 1: Using `just` (if installed)

```bash
just tidy     # Tidy up dependencies
just build    # Build to ./release/rdrop
```

### Method 2: Using Go commands

```bash
go mod tidy
go build -o ./release/rdrop
```

## 🤯 How to Use?

Once started, raindrop runs a temporary HTTP server on your device (default port is `1130`). Simply open the URL shown in your terminal from any other device (phone, tablet, PC) on the same LAN to access the shared content.

### Option Descriptions

| Short | Long             | Description                                                          | Default |
| :---- | :--------------- | :------------------------------------------------------------------- | :------ |
| `-m`  | `--message`      | The message content to send.                                         | None    |
| `-c`  | `--content-file` | Specify the path of a file whose content will be sent as plain text. | None    |
| `-p`  | `--port`         | Specify the port on which the server will run.                       | 1130    |

### Usage Examples

**Share specific files**
Directly follow the command with the paths to files or folders to be shared. You can share one or more files or directories.

```bash
rdrop iu.txt iu.png iu_folder/
```

**Send a message**
Send a plain text message using the `-m` or `--message` option.

```bash
rdrop -m "Here's a video link: https://youtu.be/xgXFCOoQJVc"
```

**Send file content**
Use the `-c` or `--content-file` option to send the content of a specified file as plain text.

```bash
rdrop -c iu_wiki.md
```

**Specify server port** (default 1130)
Use the `-p` or `--port` option to set the server's listening port.

```bash
rdrop -p 1993 -m "Love my IU"
```

**Combine all options**
All options can be used together to accommodate more complex needs.

```bash
rdrop iu_folder/ -m "Hi! These are today's IU reports" -c iu_wiki.md -p 1993
```

## 🥺 Security Notice

raindrop is designed for **trusted local networks** (e.g., home or office networks). It **does not include any authentication or encryption mechanisms** at this time. Do not use it on public or untrusted networks (e.g., cafes, airport Wi-Fi) to transmit sensitive data.
