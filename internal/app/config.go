package app

import "fmt"

// Config holds all CLI parameters.
type Config struct {
	InputDir   string
	Output     string
	TileAspect string
	TileWidth  int
	Columns    int
}

// Validate ensures required flags are provided and values make sense for the renderer.
func (c Config) Validate() error {
	if c.InputDir == "" {
		return fmt.Errorf("missing required flag: -input")
	}
	if c.TileWidth <= 0 {
		return fmt.Errorf("tile-width must be greater than zero")
	}
	if c.Columns <= 0 {
		return fmt.Errorf("columns must be greater than zero")
	}
	return nil
}
