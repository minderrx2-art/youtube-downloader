package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Worker struct {
	ctx       context.Context
	wg        *sync.WaitGroup
	renderer  *MutexProgressRender
	semaphore chan struct{}
}

type Video struct {
	bin  string
	name string
	url  string
	slot int
}

type TitleResult struct {
	url   string
	title string
}

func RunYTDLPConcurrent(ytdlp *YTDLP, urls []string, cfg Config) error {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	titleChan := make(chan TitleResult, cfg.Concurrency)
	for _, url := range urls {
		go getTitle(ctx, ytdlp.FilePath, url, titleChan)
	}
	videoMap := make(map[string]*Video)

	for range urls {
		titleResult := <-titleChan
		if titleResult.title == "" {
			warn("Skipping Invalid URL " + titleResult.url)
			continue
		}

		if _, ok := videoMap[titleResult.title]; !ok {
			videoMap[titleResult.title] = &Video{
				bin:  ytdlp.FilePath,
				name: titleResult.title,
				url:  titleResult.url,
				slot: -1,
			}
		}
	}

	renderer := NewRenderer()
	semaphore := make(chan struct{}, cfg.Concurrency)
	args, err := formatArgs(cfg.Directory)

	if err != nil {
		return err
	}

	for _, video := range videoMap {
		wg.Add(1)
		worker := newWorker(ctx, &wg, semaphore, renderer)
		go worker.start(video, args...)
	}
	wg.Wait()
	return err
}

func newWorker(ctx context.Context, wg *sync.WaitGroup, sem chan struct{}, renderer *MutexProgressRender) *Worker {
	return &Worker{
		ctx:       ctx,
		wg:        wg,
		semaphore: sem,
		renderer:  renderer,
	}
}

func (worker *Worker) start(video *Video, args ...string) {
	defer worker.wg.Done()
	worker.semaphore <- struct{}{}
	defer func() { <-worker.semaphore }()
	if err := download(worker.ctx, worker.renderer, video, args...); err != nil {

	}
}

func formatArgs(downPath string) ([]string, error) {
	if downPath == "" {
		defaultPath, err := getDownloadsDir()
		if err != nil {
			return nil, err
		}
		downPath = defaultPath
	} else {
		info, err := os.Stat(downPath)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("directory is not a directory")
		}
	}
	args := []string{
		"-o",
		filepath.Join(downPath, "%(title)s.%(ext)s"),
		"--progress-template",
		"%(progress._percent_str)s",
	}
	return args, nil
}

func getTitle(ctx context.Context, bin, url string, titleChan chan<- TitleResult) {
	cmd := ytdlpCmd(ctx, bin, "--get-title", url)

	out, err := cmd.Output()
	if err != nil {
		titleChan <- TitleResult{
			url: url,
		}
		return
	}

	titleChan <- TitleResult{
		title: strings.TrimSpace(string(out)),
		url:   url,
	}
}

func download(ctx context.Context, renderer *MutexProgressRender, video *Video, args ...string) error {
	args = append(args, video.url)
	cmd := ytdlpCmd(ctx, video.bin, args...)
	filtered := NewFilteredWriter(
		video.name,
		video.slot,
		renderer,
	)
	cmd.Stdout = filtered
	cmd.Stderr = filtered

	return cmd.Run()
}

func ytdlpCmd(ctx context.Context, bin string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, bin, args...)
}
