package src

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/mattn/go-runewidth"
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
	width int
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

// If filename is over width trims it byte by byte, otherwise pads it
func padOrTrim(s string, width int) string {
	if runewidth.StringWidth(s) > width {
		for runewidth.StringWidth(s+"…") > width {
			s = s[:len(s)-1]
		}
		return s + "…"
	}
	return s + strings.Repeat(" ", width-runewidth.StringWidth(s))
}

func NewFilteredWriter(file string, slot int, renderer *MutexProgressRender) *FilteredWriter {
	lineWidth := 45
	return &FilteredWriter{
		file:     padOrTrim(file, lineWidth),
		slot:     slot,
		renderer: renderer,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`\d+%`),
		},
	}
}

func NewRenderer() *MutexProgressRender {
	return &MutexProgressRender{
		lines: 0,
		width: 45,
	}
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
	fmt.Printf("%-*s %s", mpr.width, file, s)
	fmt.Printf("\r\033[%dB", up)
}

// Need to inject custom text for example if download fails
