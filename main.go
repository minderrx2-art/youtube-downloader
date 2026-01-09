package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"ytDownloader/src"
)

// Here be flags
type Config struct {
	c int
	u string
}

func parseConfig() Config {
	cfg := Config{}

	flag.IntVar(&cfg.c, "c", 3, "How many concurrent downloads allowed?")
	flag.StringVar(&cfg.u, "u", "", "Urls to download")

	flag.Parse()
	return cfg
}

func main() {
	cfg := parseConfig()
	ytdlp, err := src.SetupYTDLP()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Cleans up if CTRL-C terminate signal is sent
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

	if cfg.u != "" {
		urls = strings.Split(cfg.u, " ")
	} else {
		urls, err = src.ReadStdin()
		if err != nil {
			return
		}
	}

	if err := src.RunYTDLPConcurrent(ytdlp, urls, cfg.c); err != nil {
		return
	}
}
