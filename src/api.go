package src

import (
	"context"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func RunYTDLPConcurrent(ytdlp *YTDLP, urls []string) error {
	var wg sync.WaitGroup
	args, err := getArgs(urls)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err != nil {
		return err
	}

	semaphore := make(chan struct{}, 5)
	errors := make(chan error, 1)

	ytUrls := args[2:]

	renderer := NewRenderer(len(ytUrls) % 5)

	for i := 0; i < len(ytUrls) && i < 5; i++ {
		println(i)
	}

	for i, ytUrl := range ytUrls {
		wg.Add(1)
		slot := i % 5
		go func(slot int, s string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			if err := download(ctx, slot, renderer, ytdlp.FilePath, args[0], args[1], s); err != nil {

			}
		}(slot, ytUrl)
	}
	wg.Wait()

	// Basically unique channel only way to handle incoming events
	select {
	case err := <-errors:
		return err
	default:
		return nil
	}
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

func getTitle(ctx context.Context, bin, url string) (string, error) {
	cmd := ytdlpCmd(ctx, bin, "--get-title", url)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func download(ctx context.Context, slot int, renderer *MutexProgressRender, bin string, args ...string) error {
	cmd := ytdlpCmd(ctx, bin, args...)

	filtered := &FilteredWriter{
		line:     slot,
		renderer: renderer,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`\d+%`),
		},
	}

	cmd.Stdout = filtered
	cmd.Stderr = filtered

	return cmd.Run()
}

func ytdlpCmd(ctx context.Context, bin string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, bin, args...)
}
