package internal

import (
	"bufio"
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Worker struct {
	ctx       context.Context
	wg        *sync.WaitGroup
	semaphore chan struct{}
	send      func(tea.Msg)
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

	for i := range urls {
		titleResult := <-titleChan
		if titleResult.title == "" {
			continue
		}

		if _, ok := videoMap[titleResult.title]; !ok {
			videoMap[titleResult.title] = &Video{
				bin:  ytdlp.FilePath,
				name: titleResult.title,
				url:  titleResult.url,
				slot: i,
			}
		}
	}

	semaphore := make(chan struct{}, cfg.Concurrency)
	args, err := formatArgs(cfg.Directory)

	if err != nil {
		return err
	}

	titles := slices.Collect(maps.Keys(videoMap))

	teaProgram := tea.NewProgram(NewOutput(titles, len(titles)))
	for _, video := range videoMap {
		wg.Add(1)
		worker := newWorker(ctx, &wg, semaphore, teaProgram)
		go worker.start(video, args...)
	}
	_, err = teaProgram.Run()

	wg.Wait()
	return err
}

func newWorker(ctx context.Context, wg *sync.WaitGroup, sem chan struct{}, teaProgram *tea.Program) *Worker {
	return &Worker{
		ctx:       ctx,
		wg:        wg,
		semaphore: sem,
		send: func(msg tea.Msg) {
			teaProgram.Send(msg)
		},
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
		"--newline",
		"--progress-template",
		"%(progress._percent_str)s",
	}
	return args, nil
}

func getTitle(ctx context.Context, bin, url string, titleChan chan<- TitleResult) {
	cmd := exec.CommandContext(ctx, bin, "--get-title", url)

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

func (w *Worker) start(video *Video, args ...string) {
	defer w.wg.Done()
	w.semaphore <- struct{}{}
	defer func() { <-w.semaphore }()
	// Delay to let bubble tea initialize
	//
	time.Sleep(100 * time.Millisecond)
	if err := w.download(w.ctx, video, args...); err != nil {

	}
}

func (w *Worker) download(ctx context.Context, video *Video, args ...string) error {
	args = append(args, video.url)
	cmd := exec.CommandContext(ctx, video.bin, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasSuffix(line, "%") {
			continue
		}

		pct := strings.TrimSuffix(line, "%")

		w.send(VideoDebug{
			id:       video.slot,
			progress: pct,
			message:  line,
		})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}
