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

| Short Option | Full Option    | Description                                                                  | Default |
| :----------- | :------------- | :--------------------------------------------------------------------------- | :------ |
| `-m`         | `--message`    | The message content to send.                                                 | None    |
| `-c`         | `--content-file` | Specify the path to a file whose content will be sent as plain text.         | None    |
| `-d`         | `--dir`        | Specify the path to a directory whose files will be shared. Does not include dotfiles (files starting with `.`) and dot-prefixed subdirectories. | None    |
| `-p`         | `--port`       | Specify the port on which the server runs.                                   | 1130    |

### Usage Examples

**Share Specific Files**  
Share one or more files by directly following the command with their paths.

```bash
rdrop ./path/to/file1.txt ./path/to/image.jpg
```

**Send Message**  
Send plain text messages using the `-m` or `--message` option.

```bash
rdrop -m "This is a video link: https://youtu.be/xgXFCOoQJVc"
```

**Send File Content**  
Read the content of a specified file and send it as plain text using the `-c` or `--content-file` option.

```bash
rdrop -c ./path/to/iu_doc.txt
```

**Share Directory**  
Share a directory using the `-d` or `--dir` option. Note: When sharing a directory, dotfiles (files starting with `.`) and dot-prefixed subdirectories within the directory will not be included.

```bash
rdrop -d ./iu_folder
```

**Specify Server Port**  
Specify the server port for the application to run on using the `-p` or `--port` option.

```bash
rdrop -p 1993 -m "Love my IU"
```

**Combine All Options**  
All options can be combined to meet more complex requirements.

```bash
rdrop -m "Hi! These are today's materials" -c README.md -d ./downloads -p 9000 file.zip
```

## 🥺 Security Notice

raindrop is designed for **trusted local networks** like your home or office. It currently **does not provide any authentication or encryption**. Avoid using it over public or untrusted networks (such as coffee shops or airport Wi-Fi) when transferring sensitive data.
