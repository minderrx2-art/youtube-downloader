package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ytgo/internal"

	"github.com/charmbracelet/huh/spinner"
)

type Setup_Result struct {
	Result *internal.YTDLP
	Err    error
}

func main() {

	setupChan := make(chan Setup_Result, 1)

	err := spinner.New().
		Title("Downloading video...").
		Action(func() {
			ytdlp, err := internal.SetupYTDLP()
			setupChan <- Setup_Result{
				Result: ytdlp,
				Err:    err,
			}
		}).
		Run()

	if err != nil {
		panic("Can't set up YTDLP")
	}

	setup := <-setupChan
	ytdlp, err := setup.Result, setup.Err

	cfg, err := internal.PromptConfig()
	if err != nil {
		return
	}
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
