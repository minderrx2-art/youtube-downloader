package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ytDownloader/internal"
)

func main() {
	cfg := internal.ParseConfig()
	ytdlp, err := internal.SetupYTDLP()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		os.RemoveAll(ytdlp.DirPath)
		os.Exit(0)
	}()

	defer os.RemoveAll(ytdlp.DirPath)

	if err != nil {
		return
	}

	var urls []string

	if cfg.Urls != "" {
		urls = strings.Split(cfg.Urls, " ")
	} else {
		urls, err = internal.ReadStdin()
		if err != nil {
			return
		}
	}

	if err := internal.RunYTDLPConcurrent(ytdlp, urls, cfg); err != nil {
		return
	}
}
