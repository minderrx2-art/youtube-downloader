package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type YTDLP struct {
	FilePath string
	DirPath  string
}

type ytdlpResult struct {
	Result *YTDLP
	Err    error
}

func SetupYTDLP() (*YTDLP, error) {
	tempDir, err := os.MkdirTemp("", "ytgo-*")
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(tempDir, "yt-dlp")
	err = downloadYTDLP(filePath)
	if err != nil {
		return nil, err
	}
	return &YTDLP{FilePath: filePath, DirPath: tempDir}, nil
}

func downloadYTDLP(fileName string) error {
	downloadURL := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	res, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	if err := os.Chmod(fileName, 0755); err != nil {
		return err
	}
	return nil
}

func getDownloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "linux":
		return getLinuxDirectory(home)
	case "darwin":
		return filepath.Join(home, "Downloads"), nil
	case "windows":
		return filepath.Join(home, "Downloads"), nil
	default:
		return "", fmt.Errorf("Unsupported OS")
	}
}

func getLinuxDirectory(home string) (string, error) {
	downDir := filepath.Join(home, "Downloads")
	// Poke to check existance
	info, err := os.Stat(downDir)

	if err == nil && info.IsDir() {
		return downDir, nil
	}
	// Make new Dir if ~/Downloads don't exist
	err = os.Mkdir(filepath.Join(home, "ytgo"), 0755)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "ytgo"), nil
}
