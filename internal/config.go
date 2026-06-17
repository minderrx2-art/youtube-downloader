package internal

import (
	"strconv"

	"github.com/charmbracelet/huh"
)

type Config struct {
	Concurrency int
	Urls        string
	Directory   string
}

func PromptConfig() (Config, error) {
	defaultDir, err := getDownloadsDir()

	cfg := Config{
		Concurrency: 1,
		Directory:   defaultDir,
	}

	concurrency := strconv.Itoa(cfg.Concurrency)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Concurrency").
				Value(&concurrency),

			huh.NewInput().
				Title("URLs (space separated)").
				Value(&cfg.Urls),

			huh.NewInput().
				Title("Download Directory").
				Value(&cfg.Directory),
		),
	)

	if err := form.Run(); err != nil {
		return Config{}, err
	}

	n, err := strconv.Atoi(concurrency)
	if err != nil {
		n = 3
	}

	cfg.Concurrency = n

	return cfg, nil
}
