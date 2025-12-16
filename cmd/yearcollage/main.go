package main

import (
	"log"

	flag "github.com/spf13/pflag"

	"github.com/luceast/yearcollage/internal/app"
)

// main wires CLI flags into a Config and hands control to the app package.
func main() {
	cfg := app.Config{}

	// CLI flags (lowercase/kebab to match README) with short aliases.
	flag.StringVarP(&cfg.InputDir, "input", "i", "", "Input directory containing images")
	flag.StringVarP(&cfg.Output, "output", "o", "collage.jpg", "Output collage file path")
	flag.StringVarP(&cfg.TileAspect, "tile-aspect", "a", "1:1", "Target tile aspect ratio, e.g. 1:1, 3:2, 4:3")
	flag.IntVarP(&cfg.TileWidth, "tile-width", "w", 400, "Tile width in pixels")
	flag.IntVarP(&cfg.Columns, "columns", "c", 20, "Number of columns in the collage grid")
	flag.StringVarP(&cfg.CollageAspect, "collage-aspect", "r", "", "Target aspect ratio for the final collage (overrides -columns if set)")
	flag.StringVarP(&cfg.SortMode, "sort", "s", "time", "Sort images by: time (file mod time), name (alphabetical), or exif (DateTimeOriginal/DateTimeDigitized)")

	flag.Parse()

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("yearcollage failed: %v", err)
	}
}
