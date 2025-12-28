package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"ytDownloader/src"
)

// Find a way to provide custom console log output
// Like
// ## 0/4 downloaded

// Find a way to concurrently download them all at the same time
// clean up temp files after CTRL-C
// Check if files already exist in Downloads
func main() {
	urls, err := src.ReadStdin()
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
		fmt.Println(err)
		return
	}

	concurr := flag.Bool("c", false, "Use concurrent downloading?")
	flag.Parse()
	if *concurr {
		if err := src.RunYTDLPConcurrent(ytdlp, urls); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if err := src.RunYTDLPSequential(ytdlp, urls); err != nil {
			fmt.Println(err)
			return
		}
	}
}
