package watcher

import "github.com/buglloc/templed/internal/config"

type Led struct {
	MaxTemp int
	Name    string
	Color   []config.Led
}
