package app

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/luceast/yearcollage/internal/aspect"
	"github.com/luceast/yearcollage/internal/collect"
)

// Run orchestrates the YearCollage workflow (collect → sort → process → compose).
func Run(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Ensure the input path exists before walking it.
	info, err := os.Stat(cfg.InputDir)
	if err != nil {
		return fmt.Errorf("stat input dir %q: %w", cfg.InputDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("input path %q is not a directory", cfg.InputDir)
	}

	// Collect supported image files recursively.
	imagePaths, err := collect.Images(cfg.InputDir)
	if err != nil {
		return fmt.Errorf("collect images: %w", err)
	}
	if len(imagePaths) == 0 {
		return fmt.Errorf("no images found in %q", cfg.InputDir)
	}

	// Sort by filesystem mod time (oldest first). Errors are logged but do not stop sorting.
	sort.Slice(imagePaths, func(i, j int) bool {
		infoI, errI := os.Stat(imagePaths[i])
		if errI != nil {
			log.Printf("warn: stat %q: %v", imagePaths[i], errI)
			return false
		}
		infoJ, errJ := os.Stat(imagePaths[j])
		if errJ != nil {
			log.Printf("warn: stat %q: %v", imagePaths[j], errJ)
			return true
		}
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Parse target aspect ratio (e.g., "3:2" -> 1.5).
	ratio, err := aspect.Parse(cfg.TileAspect)
	if err != nil {
		return fmt.Errorf("invalid tile-aspect %q: %w", cfg.TileAspect, err)
	}

	tileWidth := cfg.TileWidth
	tileHeight := int(math.Round(float64(tileWidth) / ratio))
	if tileHeight <= 0 {
		return fmt.Errorf("computed tile height is non-positive; check tile-aspect %q", cfg.TileAspect)
	}

	log.Printf("Found %d images in %s", len(imagePaths), cfg.InputDir)
	for i, p := range imagePaths {
		if i >= 10 {
			log.Printf("... and %d more", len(imagePaths)-10)
			break
		}
		log.Printf("  %s", p)
	}

	rows := (len(imagePaths) + cfg.Columns - 1) / cfg.Columns
	canvasWidth := tileWidth * cfg.Columns
	canvasHeight := tileHeight * rows
	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	for idx, path := range imagePaths {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open image %q: %w", path, err)
		}

		img, _, err := image.Decode(f)
		_ = f.Close()
		if err != nil {
			return fmt.Errorf("decode image %q: %w", path, err)
		}

		cropped := cropToAspect(img, ratio)

		dst := image.NewRGBA(image.Rect(0, 0, tileWidth, tileHeight))
		draw.ApproxBiLinear.Scale(dst, dst.Bounds(), cropped, cropped.Bounds(), draw.Over, nil)

		col := idx % cfg.Columns
		row := idx / cfg.Columns
		offset := image.Pt(col*tileWidth, row*tileHeight)
		draw.Draw(canvas, image.Rectangle{Min: offset, Max: offset.Add(dst.Bounds().Size())}, dst, image.Point{}, draw.Src)
	}

	if err := saveImage(cfg.Output, canvas); err != nil {
		return err
	}

	log.Printf("Saved collage to %s (%dx%d)", cfg.Output, canvasWidth, canvasHeight)
	return nil
}

func cropToAspect(img image.Image, target float64) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w == 0 || h == 0 {
		return img
	}

	srcRatio := float64(w) / float64(h)
	if math.Abs(srcRatio-target) < 1e-9 {
		return img
	}

	var rect image.Rectangle
	if srcRatio > target {
		newW := int(math.Round(float64(h) * target))
		x0 := b.Min.X + (w-newW)/2
		rect = image.Rect(x0, b.Min.Y, x0+newW, b.Max.Y)
	} else {
		newH := int(math.Round(float64(w) / target))
		y0 := b.Min.Y + (h-newH)/2
		rect = image.Rect(b.Min.X, y0, b.Max.X, y0+newH)
	}

	if si, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}); ok {
		return si.SubImage(rect)
	}

	dst := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
	return dst
}

func saveImage(path string, img image.Image) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output %q: %w", path, err)
	}
	defer out.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		if err := png.Encode(out, img); err != nil {
			return fmt.Errorf("encode png %q: %w", path, err)
		}
	default:
		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("encode jpeg %q: %w", path, err)
		}
	}
	return nil
}
