package internal

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func withStdin(input string, fn func()) {
	old := os.Stdin
	defer func() { os.Stdin = old }()

	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	_, _ = w.WriteString(input)
	_ = w.Close()

	os.Stdin = r
	fn()
}

func TestReadStdin_ReadsLines(t *testing.T) {
	withStdin("a\nb\nc\n", func() {
		lines, err := ReadStdin()
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"a", "b", "c"}
		if !reflect.DeepEqual(lines, want) {
			t.Fatalf("got %v, want %v", lines, want)
		}
	})
}

func TestReadStdin_Timeout(t *testing.T) {
	old := os.Stdin
	defer func() { os.Stdin = old }()

	r, _, err := os.Pipe() // never write, never close
	if err != nil {
		t.Fatal(err)
	}

	os.Stdin = r

	start := time.Now()
	lines, err := ReadStdin()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}

	if len(lines) != 0 {
		t.Fatalf("expected no lines, got %v", lines)
	}

	if elapsed < 90*time.Millisecond {
		t.Fatalf("timeout fired too early: %v", elapsed)
	}
}
