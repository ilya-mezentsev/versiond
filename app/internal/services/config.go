package services

import "time"

type (
	Config struct {
		DefaultVersion string
		Cache          Cache
		Version        Version
		OnChangeCmds   OnChangeCmds
	}

	Cache struct {
		Dir string
	}

	Version struct {
		Source  Source
		Monitor Monitor
	}

	Source struct {
		Address string
		Timout  time.Duration
	}

	Monitor struct {
		Period time.Duration
	}

	OnChangeCmds struct {
		Before string
		Main   string
		After  string
	}
)
