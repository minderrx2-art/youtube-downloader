package src

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type FilteredWriter struct {
	file     string
	slot     int
	renderer *MutexProgressRender
	patterns []*regexp.Regexp
}

type MutexProgressRender struct {
	mutex sync.Mutex
	lines int
}

// cmd.Stdout automatically calls .Write()
func (fw *FilteredWriter) Write(p []byte) (int, error) {
	text := string(p)

	for _, re := range fw.patterns {
		if re.MatchString(text) {
			fw.renderer.Update(&fw.slot, text, fw.file)
			break
		}
	}
	return len(p), nil
}

func NewFilteredWriter(file string, slot int, renderer *MutexProgressRender) *FilteredWriter {
	r := []rune(file)
	if len(r) > 60 {
		file = string(r[:60])
		file = file + "..."
	}
	return &FilteredWriter{
		file:     file,
		slot:     slot,
		renderer: renderer,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`\d+%`),
		},
	}
}

func NewRenderer() *MutexProgressRender {
	return &MutexProgressRender{lines: 0}
}

func (mpr *MutexProgressRender) Update(line *int, s string, file string) {
	// mutex.Lock = Only one go routine can touch this at a time.
	mpr.mutex.Lock()
	defer mpr.mutex.Unlock()

	if *line == -1 {
		fmt.Println()
		*line = mpr.lines
		mpr.lines++
	}

	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")

	up := mpr.lines - *line
	fmt.Printf("\033[%dA", up)
	fmt.Print("\r\033[2K")
	fmt.Printf("%s\t%s", file, s)
	fmt.Printf("\r\033[%dB", up)
}

// Need to inject custom text for example if download fails
