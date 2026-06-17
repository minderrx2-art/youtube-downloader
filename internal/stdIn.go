package internal

import (
	"bufio"
	"os"
	"time"
)

func readRoutine(lineCh chan<- string, errCh chan<- error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lineCh <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		errCh <- err
	}
	close(lineCh)
}

func ReadStdin() ([]string, error) {
	timeout := time.Millisecond * 100
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	lines := make([]string, 0)
	lineCh := make(chan string)
	errCh := make(chan error)

	go readRoutine(lineCh, errCh)

	for {
		select {
		case line, ok := <-lineCh:
			if !ok {
				return lines, nil
			}
			lines = append(lines, line)
			timer.Reset(timeout)
		case err := <-errCh:
			return nil, err
		case <-timer.C:
			return lines, nil
		}
	}

}
