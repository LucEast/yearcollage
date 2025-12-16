package collect

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Images walks the root directory recursively and returns supported image paths.
func Images(root string) ([]string, error) {
	var images []string

	// Allowlist of extensions (case-insensitive).
	allowedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	// WalkDir avoids following symlinks and reports errors via callback.
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("warn: skipping %q: %v", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if allowedExt[ext] {
			images = append(images, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return images, nil
}
