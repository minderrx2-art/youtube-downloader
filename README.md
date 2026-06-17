# youtube-downloader

My attempt to learn GOlang, concurrent YouTube video downloader written in Go using
[yt-dlp](https://github.com/yt-dlp/yt-dlp) under the hood.

## Features

- Concurrent video downloads
- Supports input via **stdin** or **manual URL**
- Configurable concurrency level
- Simple CLI interface

## Requirements

- Go (for building)

## Build

```bash
go build -o ytDownloader ./cmd/ytdownloader
```

## Usage
Download from a file (stdin)
```bash
cat my_url_file.txt | ./ytDownloader
```
Download with URL manually
```bash
./ytDownloader -u https://youtube.com/url/here
```
Set concurrency level
```bash
./ytDownloader -c 5
```

| Flag | Description                                   |
| ---: | --------------------------------------------- |
| `-u` | Manual URL input                              |
| `-c` | Number of concurrent downloads (default: `3`) |
| `-d` | Download directory (default )                            |