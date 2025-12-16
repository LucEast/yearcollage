package app

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "columns provided",
			cfg:     Config{InputDir: "in", TileWidth: 100, Columns: 2},
			wantErr: false,
		},
		{
			name:    "collage aspect provided",
			cfg:     Config{InputDir: "in", TileWidth: 100, CollageAspect: "16:9"},
			wantErr: false,
		},
		{
			name:    "missing columns and collage aspect",
			cfg:     Config{InputDir: "in", TileWidth: 100},
			wantErr: true,
		},
		{
			name:    "negative columns",
			cfg:     Config{InputDir: "in", TileWidth: 100, Columns: -1},
			wantErr: true,
		},
		{
			name:    "invalid sort",
			cfg:     Config{InputDir: "in", TileWidth: 100, Columns: 1, SortMode: "weird"},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPickColumnsForCollage(t *testing.T) {
	cases := []struct {
		name   string
		n      int
		target float64
		want   int
	}{
		{"ten images 16:9", 10, 16.0 / 9.0, 5},
		{"ten images 1:1", 10, 1.0, 4},
		{"single image", 1, 1.0, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pickColumnsForCollage(tc.n, tc.target)
			if got != tc.want {
				t.Fatalf("pickColumnsForCollage(%d, %.3f) = %d, want %d", tc.n, tc.target, got, tc.want)
			}
		})
	}
}

func TestRunWithCollageAspectOverridesTileAspect(t *testing.T) {
	tmp := t.TempDir()

	// Create a few simple PNGs to feed into the pipeline.
	for i := 0; i < 6; i++ {
		path := filepath.Join(tmp, fmt.Sprintf("img-%02d.png", i))
		if err := writeSolidPNG(path, 40+10*i, 30+5*i, color.RGBA{uint8(20 * i), 0, 200, 255}); err != nil {
			t.Fatalf("write image %s: %v", path, err)
		}
	}

	outPath := filepath.Join(tmp, "out.png")
	cfg := Config{
		InputDir:      tmp,
		Output:        outPath,
		TileAspect:    "3:2", // should be ignored when collage-aspect is set
		TileWidth:     100,
		Columns:       0,
		CollageAspect: "1:1",
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decode output: %v", err)
	}

	b := img.Bounds()
	ratio := float64(b.Dx()) / float64(b.Dy())
	if math.Abs(ratio-1.0) > 0.05 {
		t.Fatalf("output aspect = %.3f, want close to 1.0", ratio)
	}
}

func writeSolidPNG(path string, w, h int, c color.Color) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
