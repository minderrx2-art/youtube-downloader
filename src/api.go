package src

import (
	"context"
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

func RunYTDLPConcurrent(ytdlp *YTDLP, urls []string, maxCon int) error {
	var wg sync.WaitGroup
	args, err := formatArgs()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err != nil {
		return err
	}

	renderer := NewRenderer()
	videoMap := make(map[string]*Video)

	for _, ytUrl := range urls {
		title, err := getTitle(ctx, ytdlp.FilePath, ytUrl)
		if err != nil {
			continue
		}
		if _, ok := videoMap[title]; !ok {
			videoMap[title] = &Video{
				bin:  ytdlp.FilePath,
				name: title,
				url:  ytUrl,
				slot: -1,
			}
		}
	}

	semaphore := make(chan struct{}, maxCon)

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

func formatArgs() ([]string, error) {
	downloadsPath, err := getDownloadsDir()
	if err != nil {
		return nil, err
	}
	args := []string{
		"-o",
		filepath.Join(downloadsPath, "%(title)s.%(ext)s"),
		"--progress-template",
		"%(progress._percent_str)s",
	}
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
