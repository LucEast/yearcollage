package aspect

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse converts strings like "3:2" into a numeric aspect ratio (width / height).
func Parse(value string) (float64, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid aspect ratio format: %q", value)
	}

	w, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid width in aspect ratio: %w", err)
	}
	h, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid height in aspect ratio: %w", err)
	}
	if h == 0 {
		return 0, fmt.Errorf("height in aspect ratio cannot be zero")
	}

	return w / h, nil
}
