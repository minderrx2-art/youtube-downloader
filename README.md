# youtube-downloader

My attempt to learn GOlang, concurrent YouTube video downloader written in Go using
[yt-dlp](https://github.com/yt-dlp/yt-dlp) under the hood.

## Features

- Concurrent video downloads
- Configurable concurrency level
- Simple CLI interface

## Requirements

- Go (for building)

## Build

```bash
go build -o ytgo ./cmd
```

## Usage
<img width="544" height="203" alt="image" src="https://github.com/user-attachments/assets/ad3127df-b9b2-4641-98c4-67de35deb45b" />

| Flag | Description                                                     |
| ---: | ---------------------------------------------                   |
| `Urls to download` | Manual URL input                                  |
| `Concurrency level` | Number of concurrent downloads (default: `1`)    |
| `Default directory` | Download directory (default )                    |
