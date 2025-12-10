package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)


//
// Config – Struktur zum Ablegen aller CLI-Parameter
//
type Config struct {
	InputDir	string // Folder with images
	OutputDir	string // Destinationpath of the collage
	TileAspect	string // desired aspect ratio e.g. "3:2"
	TileWidth	int // width of a tile in pixels
	Columns int // quantity of columns in the final grid
}

//
// parseFlags() - reads the cli-arguments and validates them
//
func parseFlags() Config {
	cfg := Config{}

	// define CLI-Flags 
	flag.StringVar(&cfg.InputDir, "Input", "", "Input directory containing images")
	flag.StringVar(&cfg.OutputFile, "output", "collage.jpg", "Output collage file path")
	flag.StringVar(&cfg.TileAspect, "tile-aspect", "1:1", "Target tile aspect ratio, e.g. 1:1, 3:2, 4:3")
	flag.IntVar(&cfg.TileWidth, "tile-width", 400, "Tile width in pixels")
	flag.IntVar(&cfg.Columns, "columns", 20, "Number of columns in the collage grid")

	flag.Parse() // Flags tatsächlich einlesen

	// Validierung
	if cfg.InputDir == "" {
		log.Fatal("missing required flag: -input")
	}

	return cfg
}

// main() - Einstiegspunkt des Programms
func main() {
	cfg := parseFlags()

	// Prüfen, ob dir Input-Ordner existiert
	info, err := os.Stat(cfg.InputDir)
	if err != nil {
		log.Fatalf("failed to stat input dir %q: %v", cfg.InputDir, err)
	}
	if !info.IsDir() {
		log.Fatalf("input path %q is not a directory", cfg.InputDir)
	}

	// Bilder aus Ordner + Unterordnern holen
	imagePaths, err := collectImages(cfg.InputDir)
	if err != nil {
		log.Fatalf("failed to collect images: %v", err)
	}

	if len(imagePaths) == 0 {
		log.Fatalf("no images found in %q", cfg.InputDir)
	}

	// Bilder nach Datum sortieren (älteste zuerst)
	sort.Slice(imagePaths, func(i, j int) bool {
    infoI, _ := os.Stat(imagePaths[i])
    infoJ, _ := os.Stat(imagePaths[j])
    return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Erfolgsmeldung
	fmt.Printf("Found %d images in %s\n", len(imagePaths), cfg.InputDir)

	// Nur die ersten 10 Bilder anzeigen (kein Spam)
	for i, p := range imagePaths {
		if i >= 10 {
			fmt.Printf("... and %d more\n", len(imagePaths)-10)
			break
		}
		fmt.Printf("  %s\n", p)
	}

	// --------------------------------------------------------------
	// TODO (bauen wir in den nächsten Schritten):
	// 1. tile-aspect parsen ("3:2" -> ratio 1.5)
	// 2. Bilder einzeln laden, auf ratio croppen & resizen
	// 3. Grid berechnen (rows = ceil(n / columns))
	// 4. große Canvas erstellen
	// 5. Bilder reinzeichnen
	// 6. finale Collage speichern
	// --------------------------------------------------------------

}

// collectImages() - durchsucht den Ordner rekursiv nach Bilddateien

func collectImages(root string) ([]string, error) {
	var images []string

	// Liste erlaubter Endungen (alles in Kleinbuchstaben)
	allowedExt := map[string]bool{
		".jpg": true,
		".jpeg": true,
		".png": true,
		".webp": true,
	}

	// WalkDir läuft den Ordner rekursiv durch
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// z.B. Permission denied - einfach weiter machen
			log.Printf("warnung: überspringt %q: %v", path, err)
			return nil 
		} 

		// Ordner ignorieren
		if d.IsDir() {
			return nil
		}

		//Dateiendung checken
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if allowedExt[ext] {
			images = append(images, path)
		}

		return nil
	})
	
	// Fehler des Walkers zurückgeben
	if err != nil {
		retrun nil, err
	}
	return images, nir
}

