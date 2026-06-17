package internal

import (
	"os"
	"path"
	"runtime"
	"testing"
)

func TestSetupYTDLP_CreatesExecutable(t *testing.T) {
	ytdlp, err := SetupYTDLP()
	if err != nil {
		t.Fatal(err)
	}

	if ytdlp.FilePath == "" {
		t.Fatal("FilePath is empty")
	}
	if ytdlp.DirPath == "" {
		t.Fatal("DirPath is empty")
	}

	info, err := os.Stat(ytdlp.FilePath)
	if err != nil {
		t.Fatalf("yt-dlp file does not exist: %v", err)
	}

	// check executable bit (unix only)
	if runtime.GOOS != "windows" {
		if info.Mode()&0111 == 0 {
			t.Fatalf("yt-dlp is not executable: mode=%v", info.Mode())
		}
	}
}
func TestGetLinuxDirectory_UsesDownloadsIfExists(t *testing.T) {
	home := t.TempDir()
	downloads := path.Join(home, "Downloads")

	if err := os.Mkdir(downloads, 0755); err != nil {
		t.Fatal(err)
	}

	dir, err := getLinuxDirectory(home)
	if err != nil {
		t.Fatal(err)
	}

	if dir != downloads {
		t.Fatalf("got %q, want %q", dir, downloads)
	}
}
func TestGetLinuxDirectory_CreatesFallback(t *testing.T) {
	home := t.TempDir()

	dir, err := getLinuxDirectory(home)
	if err != nil {
		t.Fatal(err)
	}

	want := path.Join(home, "ytgo")
	if dir != want {
		t.Fatalf("got %q, want %q", dir, want)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("%q is not a directory", dir)
	}
}
