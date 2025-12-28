package src

import (
	"fmt"
	"regexp"
	"sync"
)

type FilteredWriter struct {
	line     int
	renderer *MutexProgressRender
	patterns []*regexp.Regexp
}

type MutexProgressRender struct {
	mutex sync.Mutex
	lines int
}

// cmd.Stdout automatically calls .Write()
func (fw *FilteredWriter) Write(p []byte) (int, error) {
	line := string(p)

	for _, re := range fw.patterns {
		if re.MatchString(line) {
			fw.renderer.Update(fw.line, line)
			break
		}
	}
	return len(p), nil
}

func NewRenderer(lines int) *MutexProgressRender {
	// No remainder = need all lines
	if lines == 0 {
		lines = 5
	}
	return &MutexProgressRender{lines: lines}
}

func (r *MutexProgressRender) Update(line int, s string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// move cursor up from bottom
	up := r.lines - line
	fmt.Printf("\033[%dA", up)

	// clear line and print
	fmt.Print("\r\033[2K")
	fmt.Print(s)

	// move cursor back down
	fmt.Printf("\r\033[%dB", up)
}

// Need to inject custom text for example if download fails
