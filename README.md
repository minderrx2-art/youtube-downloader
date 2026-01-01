# youtube-downloader

Concurrent youtube video downloader using [yt-](https://github.com/yt-dlp/yt-dlp)

Compile with
go build .

Run via stdin:
cat my_url_file.txt | ytDownloader

Run via manual url input
./ytDownloader -u youtube/url/here

Supported Flags
-u = manual URL input.
-c = number of concurrent downloads allowed, default = 3
