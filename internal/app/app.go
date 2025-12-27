package app

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/rwcarlsen/goexif/exif"

	"github.com/luceast/yearcollage/internal/aspect"
	"github.com/luceast/yearcollage/internal/collect"
)

// Run orchestrates the YearCollage workflow (collect → sort → process → compose).
// It validates the config, gathers all supported images, resizes/crops them to
// the requested aspect ratio, and finally writes the collage to disk.
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

	imagePaths = sortImages(imagePaths, cfg.SortMode)

	log.Printf("Found %d images in %s", len(imagePaths), cfg.InputDir)
	for i, p := range imagePaths {
		if i >= 10 {
			log.Printf("... and %d more", len(imagePaths)-10)
			break
		}
		log.Printf("  %s", p)
	}

	columns := cfg.Columns
	var tileRatio float64

	if cfg.CollageAspect != "" {
		collageRatio, err := aspect.Parse(cfg.CollageAspect)
		if err != nil {
			return fmt.Errorf("invalid collage-aspect %q: %w", cfg.CollageAspect, err)
		}
		// When a collage aspect is provided we compute a column count that best
		// matches the overall shape and derive the tile aspect from it.
		columns = pickColumnsForCollage(len(imagePaths), collageRatio)
		if columns <= 0 {
			return fmt.Errorf("computed columns is non-positive")
		}
		rows := (len(imagePaths) + columns - 1) / columns
		tileRatio = collageRatio * float64(rows) / float64(columns)
		log.Printf("Collage aspect %s -> columns=%d, rows=%d, tile-aspect=%.4f (tile-aspect flag ignored)", cfg.CollageAspect, columns, rows, tileRatio)
	} else {
		ratio, err := aspect.Parse(cfg.TileAspect)
		if err != nil {
			return fmt.Errorf("invalid tile-aspect %q: %w", cfg.TileAspect, err)
		}
		tileRatio = ratio
		log.Printf("Tile aspect %s (from flag)", cfg.TileAspect)
	}

	tileWidth := cfg.TileWidth
	// Scale tile height from width so we always respect the intended ratio,
	// even when tileRatio came from collage-aspect inference.
	tileHeight := int(math.Round(float64(tileWidth) / tileRatio))
	if tileHeight <= 0 {
		return fmt.Errorf("computed tile height is non-positive; check tile/collage aspect")
	}

	rows := (len(imagePaths) + columns - 1) / columns
	canvasWidth := tileWidth * columns
	canvasHeight := tileHeight * rows
	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	for idx, path := range imagePaths {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open image %q: %w", path, err)
		}

		// Read the orientation before decoding so we can rewind and reuse the
		// same file handle for the actual pixel data.
		orientation := imageOrientation(f)
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			_ = f.Close()
			return fmt.Errorf("rewind image %q: %w", path, err)
		}

		img, _, err := image.Decode(f)
		_ = f.Close()
		if err != nil {
			return fmt.Errorf("decode image %q: %w", path, err)
		}

		img = normalizeOrientation(img, orientation)

		// Trim the photo so it fits the target aspect without stretching.
		cropped := cropToAspect(img, tileRatio)

		dst := image.NewRGBA(image.Rect(0, 0, tileWidth, tileHeight))
		draw.ApproxBiLinear.Scale(dst, dst.Bounds(), cropped, cropped.Bounds(), draw.Over, nil)

		col := idx % columns
		row := idx / columns
		offset := image.Pt(col*tileWidth, row*tileHeight)
		draw.Draw(canvas, image.Rectangle{Min: offset, Max: offset.Add(dst.Bounds().Size())}, dst, image.Point{}, draw.Src)
	}

	if err := saveImage(cfg.Output, canvas); err != nil {
		return err
	}

	log.Printf("Saved collage to %s (%dx%d)", cfg.Output, canvasWidth, canvasHeight)
	return nil
}

// cropToAspect returns a view of the image cropped to the target aspect ratio,
// preferring center crops so the main subject is likely preserved.
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

// saveImage picks an encoder based on the output extension and writes the image.
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

// pickColumnsForCollage picks a column count for a target collage aspect.
// It prefers grids that keep the inferred tile aspect near 1:1 to minimize cropping.
func pickColumnsForCollage(numImages int, targetCollageRatio float64) int {
	if numImages <= 0 {
		return 0
	}

	ideal := math.Sqrt(float64(numImages) * targetCollageRatio)
	best := clampInt(int(math.Round(ideal)), 1, numImages)
	bestScore := math.Abs(tileAspectFromGrid(numImages, best, targetCollageRatio) - 1.0)

	for delta := -3; delta <= 3; delta++ {
		c := clampInt(int(math.Round(ideal))+delta, 1, numImages)
		score := math.Abs(tileAspectFromGrid(numImages, c, targetCollageRatio) - 1.0)
		if score < bestScore || (score == bestScore && c < best) {
			best = c
			bestScore = score
		}
	}

	return best
}

// tileAspectFromGrid derives the tile aspect ratio implied by a collage grid.
func tileAspectFromGrid(numImages, columns int, collageRatio float64) float64 {
	rows := (numImages + columns - 1) / columns
	return collageRatio * float64(rows) / float64(columns)
}

// clampInt bounds v to [min, max].
func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// sortImages orders image paths according to the chosen sort mode.
func sortImages(paths []string, mode string) []string {
	switch mode {
	case "", "time":
		sort.Slice(paths, func(i, j int) bool {
			infoI, errI := os.Stat(paths[i])
			if errI != nil {
				log.Printf("warn: stat %q: %v", paths[i], errI)
				return false
			}
			infoJ, errJ := os.Stat(paths[j])
			if errJ != nil {
				log.Printf("warn: stat %q: %v", paths[j], errJ)
				return true
			}
			return infoI.ModTime().Before(infoJ.ModTime())
		})
	case "name":
		sort.Strings(paths)
	case "exif":
		type item struct {
			path string
			ts   time.Time
		}
		items := make([]item, 0, len(paths))
		for _, p := range paths {
			t := exifTime(p)
			items = append(items, item{path: p, ts: t})
		}
		sort.Slice(items, func(i, j int) bool {
			if items[i].ts.Equal(items[j].ts) {
				return items[i].path < items[j].path
			}
			return items[i].ts.Before(items[j].ts)
		})
		paths = paths[:0]
		for _, it := range items {
			paths = append(paths, it.path)
		}
	default:
		log.Printf("warn: unknown sort mode %q, falling back to time", mode)
		return sortImages(paths, "time")
	}
	return paths
}

// imageOrientation extracts the EXIF orientation flag and returns a value between
// 1 and 8 (per the TIFF/EXIF spec). When the file has no EXIF block or the tag
// is missing we default to 1 (top-left).
func imageOrientation(rs io.ReadSeeker) int {
	if rs == nil {
		return 1
	}
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		return 1
	}
	x, err := exif.Decode(rs)
	if err != nil {
		return 1
	}
	field, err := x.Get(exif.Orientation)
	if err != nil {
		return 1
	}
	val, err := field.Int(0)
	if err != nil {
		return 1
	}
	orientation := int(val)
	if orientation < 1 || orientation > 8 {
		return 1
	}
	return orientation
}

// exifTime extracts the best-effort EXIF timestamp, falling back to modtime.
func exifTime(path string) time.Time {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("warn: open for exif %q: %v", path, err)
		return modTimeOrZero(path)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return modTimeOrZero(path)
	}

	if tm, err := x.DateTime(); err == nil {
		return tm
	}
	for _, tag := range []exif.FieldName{exif.DateTimeOriginal, exif.DateTimeDigitized} {
		if field, err := x.Get(tag); err == nil {
			if s, err := field.StringVal(); err == nil {
				if tm, ok := parseExifTimeString(s); ok {
					return tm
				}
			}
		}
	}
	return modTimeOrZero(path)
}

// parseExifTimeString handles a handful of timestamp formats commonly seen in EXIF.
func parseExifTimeString(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		"2006:01:02 15:04:05",
		time.RFC3339,
		time.RFC3339Nano,
	}
	for _, layout := range layouts {
		if tm, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return tm, true
		}
	}
	return time.Time{}, false
}

// modTimeOrZero reports file modification time or zero on error, logging the issue.
func modTimeOrZero(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		log.Printf("warn: stat %q: %v", path, err)
		return time.Time{}
	}
	return info.ModTime()
}
