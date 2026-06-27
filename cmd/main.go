package main

import (
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"ytgo/internal"

	"github.com/charmbracelet/huh/spinner"
)

type Setup_Result struct {
	Result *internal.YTDLP
	Err    error
}

func IsYouTubeURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}

	switch u.Host {
	case "youtube.com", "www.youtube.com":
		if u.Path != "/watch" {
			return false
		}
		id := u.Query().Get("v")
		return len(id) == 11

	case "youtu.be":
		id := strings.TrimPrefix(u.Path, "/")
		return len(id) == 11
	}

	return false
}
func main() {

	cfg, err := internal.PromptConfig()
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	var urls []string
	re := regexp.MustCompile(`&.*`)

	if cfg.Urls != "" {
		for _, url := range strings.Split(cfg.Urls, " ") {
			if IsYouTubeURL(url) {
				urls = append(urls, re.ReplaceAllString(url, ""))
			}
		}
	} else {
		urls, err = internal.ReadStdin()
		if err != nil {
			return
		}
	}

	if len(urls) < 1 {
		panic("No valid Urls presented")
	}

	setupChan := make(chan Setup_Result, 1)

	err = spinner.New().
		Title("downloading libs").
		Action(func() {
			ytdlp, err := internal.SetupYTDLP()
			time.Sleep(200)
			setupChan <- Setup_Result{
				Result: ytdlp,
				Err:    err,
			}
		}).
		Run()

	if err != nil {
		panic("failed to set up libs")
	}

	// YTDLP setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	setup := <-setupChan
	ytdlp, err := setup.Result, setup.Err

	go func() {
		<-sigChan
		os.RemoveAll(ytdlp.DirPath)
		os.Exit(0)
	}()

	defer os.RemoveAll(ytdlp.DirPath)

	if err := internal.RunYTDLPConcurrent(ytdlp, urls, cfg); err != nil {
		return
	}
}
