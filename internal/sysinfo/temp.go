package sysinfo

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func Temp() (float64, error) {
	return ZoneTemp(0)
}

func ZoneTemp(zone int) (float64, error) {
	tempBytes, err := os.ReadFile(fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", zone))
	if err != nil {
		return 0, fmt.Errorf("unable to read sys temp: %w", err)
	}

	tempStr := string(bytes.TrimSpace(tempBytes))
	temp, err := strconv.ParseUint(tempStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid temp returned %q: %w", tempStr, err)
	}

	return float64(temp) / 1000, nil
}
