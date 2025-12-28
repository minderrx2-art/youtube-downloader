package src

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"runtime"
)

type YTDLP struct {
	FilePath string
	DirPath  string
}

//go:embed bin/yt-dlp
var ytDLPBin []byte

// Saves binary to temporary directory and returns the binary path, and path to temp dir
func SetupYTDLP() (*YTDLP, error) {
	tempDir, err := os.MkdirTemp("", "ytdownloader-*")
	if err != nil {
		return nil, err
	}
	filePath := path.Join(tempDir, "yt-dlp")
	if err := os.WriteFile(filePath, ytDLPBin, 0755); err != nil {
		return nil, err
	}
	return &YTDLP{FilePath: filePath, DirPath: tempDir}, nil
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
		return path.Join(home, "Downloads"), nil
	case "windows":
		return path.Join(home, "Downloads"), nil
	default:
		return "", fmt.Errorf("Unsupported OS")
	}
}

func getLinuxDirectory(home string) (string, error) {
	downDir := path.Join(home, "Downloads")
	// Poke to check existance
	info, err := os.Stat(downDir)
	if info.IsDir() && err == nil {
		return downDir, nil
	}
	// Make new Dir if ~/Downloads don't exist
	err = os.Mkdir(path.Join(home, "ytDownloader"), 0755)
	if err != nil {
		return "", err
	}
	return path.Join(home, "ytDownloader"), nil
}
