package app

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ilya-mezentsev/versiond/app/internal/services"
)

type (
	config struct {
		DefaultVersion string                `json:"default_version"`
		Cache          services.Cache        `json:"cache"`
		Version        version               `json:"version"`
		OnChangeCmds   services.OnChangeCmds `json:"on_change_cmds"`
	}

	version struct {
		Source  source  `json:"source"`
		Monitor monitor `json:"monitor"`
	}

	source struct {
		Address string   `json:"address"`
		Timeout Duration `json:"timeout"`
	}

	monitor struct {
		Period Duration `json:"period"`
	}

	Duration struct {
		time.Duration
	}
)

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unexpected duration with type %T and value %v", value, value)
	}
}

func mustParseConfig(filePath string) services.Config {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("unable to read confif file: %s, %v", filePath, err))
	}

	var cfg config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal config: %v", err))
	}

	return services.Config{
		DefaultVersion: cfg.DefaultVersion,
		Cache:          cfg.Cache,
		Version: services.Version{
			Source: services.Source{
				Address: cfg.Version.Source.Address,
				Timout:  cfg.Version.Source.Timeout.Duration,
			},
			Monitor: services.Monitor{
				Period: cfg.Version.Monitor.Period.Duration,
			},
		},
		OnChangeCmds: cfg.OnChangeCmds,
	}
}
