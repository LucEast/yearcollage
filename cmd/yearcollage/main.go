package main

import (
	"flag"
	"log"

	"github.com/luceast/yearcollage/internal/app"
)

// main wires CLI flags into a Config and hands control to the app package.
func main() {
	cfg := app.Config{}

	// CLI flags (lowercase/kebab to match README)
	flag.StringVar(&cfg.InputDir, "input", "", "Input directory containing images")
	flag.StringVar(&cfg.Output, "output", "collage.jpg", "Output collage file path")
	flag.StringVar(&cfg.TileAspect, "tile-aspect", "1:1", "Target tile aspect ratio, e.g. 1:1, 3:2, 4:3")
	flag.IntVar(&cfg.TileWidth, "tile-width", 400, "Tile width in pixels")
	flag.IntVar(&cfg.Columns, "columns", 20, "Number of columns in the collage grid")

	flag.Parse()

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("yearcollage failed: %v", err)
	}
}
