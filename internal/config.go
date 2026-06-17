package internal

import "flag"

type Config struct {
	Concurrency int
	Urls        string
	Directory   string
}

func ParseConfig() Config {
	cfg := Config{}

	flag.IntVar(&cfg.Concurrency, "c", 3, "How many concurrent downloads allowed?")
	flag.StringVar(&cfg.Urls, "u", "", "Urls to download")
	flag.StringVar(&cfg.Directory, "d", "", "Directory to download to (default: Downloads directory)")

	flag.Parse()
	return cfg
}
