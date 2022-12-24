package watcher

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/buglloc/templed/internal/config"
	"github.com/buglloc/templed/internal/ledctl"
	"github.com/buglloc/templed/internal/sysinfo"
)

type Watcher struct {
	curLed      string
	leds        []Led
	fallbackLed Led
	thermalZone int
	checkPeriod time.Duration
	ctx         context.Context
	cancelCtx   context.CancelFunc
	closed      chan struct{}
}

func NewWatcher(cfg *config.Config) (*Watcher, error) {
	leds := make([]Led, 0, len(cfg.Temps))
	var fallbackLed Led
	for temp, colorName := range cfg.Temps {
		color, ok := cfg.Colors[colorName]
		if !ok {
			return nil, fmt.Errorf("unknown color %q for temp %d", colorName, temp)
		}

		led := Led{
			Name:    colorName,
			MaxTemp: temp,
			Color:   color,
		}

		if temp == 0 {
			fallbackLed = led
			continue
		}

		leds = append(leds, led)
	}
	sort.Slice(leds, func(i, j int) bool {
		return leds[i].MaxTemp < leds[j].MaxTemp
	})

	if len(leds) == 0 {
		return nil, errors.New("no temps configured")
	}

	ctx, cancel := context.WithCancel(context.Background())
	out := &Watcher{
		curLed:      "",
		leds:        leds,
		fallbackLed: fallbackLed,
		thermalZone: cfg.ThermalZone,
		checkPeriod: cfg.Period,
		ctx:         ctx,
		cancelCtx:   cancel,
		closed:      make(chan struct{}),
	}

	return out, nil
}

func (w *Watcher) Watch() error {
	defer close(w.closed)

	log.Info().Msg("starts initial led sync")
	if err := w.SyncTemp(); err != nil {
		return fmt.Errorf("initial sync failed: %w", err)
	}

	log.Info().Msg("starts periodical watcher")
	ticker := time.NewTicker(w.checkPeriod)
	for {
		select {
		case <-w.ctx.Done():
			ticker.Stop()
			_ = w.setLed(w.fallbackLed, 0)
			return nil
		case <-ticker.C:
			if err := w.SyncTemp(); err != nil {
				log.Error().Err(err).Msg("sync failed")
				continue
			}
		}
	}
}

func (w *Watcher) SyncTemp() error {
	led, temp := w.tempLed()
	if led.Name == w.curLed {
		log.Info().Float64("temp", temp).Msg("nothing changed")
		return nil
	}

	return w.setLed(led, temp)
}

func (w *Watcher) setLed(led Led, temp float64) error {
	if led.Name == "" {
		return nil
	}

	log.Info().
		Float64("temp", temp).
		Int("max_temp", led.MaxTemp).
		Str("led", led.Name).
		Msg("led changed")

	w.curLed = led.Name
	return ledctl.SetLeds(led.Color...)
}

func (w *Watcher) tempLed() (Led, float64) {
	temp, err := sysinfo.ZoneTemp(w.thermalZone)
	if err != nil {
		log.Error().
			Err(err).
			Msg("unable to get temp, fallback color will be used")
		return w.fallbackLed, temp
	}

	var out Led
	for _, led := range w.leds {
		out = led
		if int(temp) <= led.MaxTemp {
			break
		}
	}
	return out, temp
}

func (w *Watcher) Shutdown(ctx context.Context) {
	w.cancelCtx()

	select {
	case <-ctx.Done():
	case <-w.closed:
	}
}
