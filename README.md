# raindrop (rd)

[English](README.md) | [简体中文](README.zh-CN.md)

raindrop is an extremely simple and easy-to-use command-line tool designed to help you instantly share files, text messages, or plain text content within a local area network (LAN). It is ideal for temporary or peer-to-peer data exchange scenarios.

It supports the following sharing modes, which can be flexibly combined:

- 📁 **File Sharing**: Share any number of files. Recipients can download them directly from their browser.
- 📄 **Plain Text Sharing**: Display the content of text files like `.txt`, `.md`, or `.log` directly in the recipient's browser without requiring a download.
- 💬 **Message Sending**: Send a short text message.

This tool is developed in Go, offering zero-configuration, lightweight efficiency, and cross-platform support. It runs as a single binary file.

## ✿ Why raindrop?

Inspired by the song "Rain Drop" by IU.  
Wanna listen?

- [Spotify](https://open.spotify.com/track/6tlMVCqZlmxfnjZt3OiHjE)
- [171210 IU Palette Rain Drop (YouTube)](https://youtu.be/xgXFCOoQJVc)

Note: For ease of use, the command-line tool is named `rd`. You will need to use `rd` to execute commands.

## 📦 How to Install?

This project provides a [`justfile`](https://github.com/casey/just) to simplify common tasks.

### Method 1: Using just (if installed)

```bash
just tidy     # Tidy dependencies
just build    # Build to ./release/rd
```

### Method 2: Using Go commands

```bash
go mod tidy
go build -o ./release/rd
```

## 🤯 How to Use?

When raindrop starts, it runs a temporary HTTP server on your device (listening on port 1130 by default). Simply open the URL displayed in the terminal in a browser on any other device (phone, tablet, computer) within the same LAN to access the shared content.

### Flag Descriptions

| Short Flag | Long Flag      | Description                                                              | Default |
| :--------- | :------------- | :----------------------------------------------------------------------- | :------ |
| `-m`       | `--message`    | The message content to be sent.                                          | None    |
| `-c`       | `--content-file` | Specifies the path to a file whose content will be sent as plain text. | None    |
| `-p`       | `--port`       | Specifies the port for the server to run on.                             | 1130    |

### Usage Examples

**Start directly**  
Shares all non-hidden files in the current directory by default.

```bash
rd
```

**Share specific files**  
Specify one or more files or directories.  
When sharing a directory, the program will exclude hidden files (those starting with a '.') by default and only share the top-level files within that directory. The contents of subdirectories are not included.

```bash
rd iu.txt iu.png iu_folder/
```

**Send a message**  
Send a plain text message using the `-m` or `--message` flag.

```bash
rd -m "Here is a video link: https://youtu.be/xgXFCOoQJVc"
```

**Send file content**  
Read the content of a specified file and send it as plain text using the `-c` or `--content-file` flag.

```bash
rd -c iu_wiki.md
```

**Specify server port** (default 1130)  
Specify the server port for the application using the `-p` or `--port` flag.

```bash
rd -p 1993 -m "Love my IU"
```

**Combine all flags**  
All flags can be used together to meet more complex needs.

```bash
rd iu_folder/ -m "Hi! Here is today's IU report" -c iu_wiki.md -p 1993
```

## 🥺 Security Tip

raindrop is designed for **trusted local area network (LAN)** environments (such as home or office networks). It currently **does not include any authentication or encryption mechanisms**. Please do not use it to transfer sensitive data over public or untrusted networks (like coffee shops or airport Wi-Fi).
