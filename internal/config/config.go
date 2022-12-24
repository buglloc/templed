package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Led struct {
	Device     string `yaml:"dev"`
	Brightness int    `yaml:"brightness"`
}

type Config struct {
	ThermalZone int              `yaml:"thermal_zone"`
	Period      time.Duration    `yaml:"period"`
	Colors      map[string][]Led `yaml:"colors"`
	Temps       map[int]string   `yaml:"temps"`
}

func LoadConfig(configs ...string) (*Config, error) {
	out := &Config{
		ThermalZone: 0,
		Period:      30 * time.Second,
		Colors: map[string][]Led{
			"fallback": {
				{
					Device:     "surround:blue",
					Brightness: 0,
				},
				{
					Device:     "surround:green",
					Brightness: 0,
				},
				{
					Device:     "surround:red",
					Brightness: 0,
				},
			},
			"cool": {
				{
					Device:     "surround:blue",
					Brightness: 0,
				},
				{
					Device:     "surround:green",
					Brightness: 1,
				},
				{
					Device:     "surround:red",
					Brightness: 1,
				},
			},
			"warm": {
				{
					Device:     "surround:blue",
					Brightness: 1,
				},
				{
					Device:     "surround:green",
					Brightness: 0,
				},
				{
					Device:     "surround:red",
					Brightness: 1,
				},
			},
			"hot": {
				{
					Device:     "surround:blue",
					Brightness: 0,
				},
				{
					Device:     "surround:green",
					Brightness: 0,
				},
				{
					Device:     "surround:red",
					Brightness: 1,
				},
			},
		},
		Temps: map[int]string{
			0:    "fallback",
			50:   "cool",
			60:   "warm",
			1000: "hot",
		},
	}

	if len(configs) == 0 {
		return out, nil
	}

	for _, cfgPath := range configs {
		err := func() error {
			f, err := os.Open(cfgPath)
			if err != nil {
				return fmt.Errorf("unable to open config file: %w", err)
			}
			defer func() { _ = f.Close() }()

			if err := yaml.NewDecoder(f).Decode(&out); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}

			return nil
		}()
		if err != nil {
			return nil, fmt.Errorf("unable to load config %q: %w", cfgPath, err)
		}
	}

	return out, nil
}
