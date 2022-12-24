package ledctl

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/buglloc/templed/internal/config"
)

func SetLeds(leds ...config.Led) error {
	for _, led := range leds {
		if err := SetLed(led); err != nil {
			return fmt.Errorf("unable to change led %q: %w", led.Device, err)
		}
	}

	return nil
}

func SetLed(led config.Led) error {
	ledCtl := filepath.Join("/sys/class/leds", led.Device, "brightness")
	f, err := os.OpenFile(ledCtl, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("unable to open brightness control: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(strconv.Itoa(led.Brightness)); err != nil {
		return fmt.Errorf("write to brightness control %q failed: %w", ledCtl, err)
	}
	return nil
}
