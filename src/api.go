package src

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type FilteredWriter struct {
	patterns []*regexp.Regexp
}

func RunYTDLPSequential(ytdlp *YTDLP, urls []string) error {
	args, err := getArgs(urls)
	if err != nil {
		return err
	}
	if err := runYTDLP(ytdlp.FilePath, args...); err != nil {
		return err
	}
	// How do I do I summarise this output without it being as verbose ???
	return nil
}

// Buffered channel with 5 slots,
// Adds 1 and blocks when channel is full,
// Evicts from channel.
func RunYTDLPConcurrent(ytdlp *YTDLP, urls []string) error {
	var wg sync.WaitGroup
	args, err := getArgs(urls)
	if err != nil {
		return err
	}
	semaphore := make(chan bool, 5)
	for _, ytUrl := range args[2:] {
		wg.Add(1)
		semaphore <- true
		go func(s string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			if err := runYTDLP(ytdlp.FilePath, args[0], args[1], s); err != nil {
				// How do I handle concurrent errors??
				// I can't return them watafak??
			}
		}(ytUrl)
	}
	wg.Wait()
	return nil
}

func getArgs(urls []string) ([]string, error) {
	downloadsPath, err := getDownloadsDir()
	if err != nil {
		return nil, err
	}
	args := []string{
		"-o",
		filepath.Join(downloadsPath, "%(title)s.%(ext)s"),
	}
	args = append(args, urls...)
	return args, nil
}

func getVideoTitle(ydlpPath string, url string) (string, error) {
	cmd := exec.Command(ydlpPath, "--get-title", url)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func runYTDLP(ydlpPath string, args ...string) error {
	if len(args) < 3 {
		return fmt.Errorf("URL was not given to YTDLP")
	}
	cmd := exec.Command(
		ydlpPath,
		args...,
	)
	handleOutput(cmd)
	return cmd.Run()
}

// cmd.Stdout automatically calls .Write()
func (fw *FilteredWriter) Write(p []byte) (n int, err error) {
	line := string(p)

	for _, pattern := range fw.patterns {
		if pattern.MatchString(line) {
			fmt.Print(line)
			break
		}
	}
	return len(p), nil
}

func handleOutput(stdo *exec.Cmd) {
	filtered := &FilteredWriter{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`\d+%`),
		},
	}

	stdo.Stdout = filtered
	stdo.Stderr = filtered
}
