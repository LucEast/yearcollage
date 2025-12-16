package collect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImages(t *testing.T) {
	// Use a temporary tree so the test never touches real user files.
	root := t.TempDir()

	// Create a mix of supported and unsupported files across nested directories.
	files := []string{
		"a.jpg",
		"b.png",
		"c.txt",
		"nested/d.webp",
		"nested/e.JPEG",
	}

	for _, name := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatalf("write file %s: %v", path, err)
		}
	}

	// Collect all images under root; should ignore unsupported files.
	paths, err := Images(root)
	if err != nil {
		t.Fatalf("Images returned error: %v", err)
	}

	// Map of paths we expect to see from the collector.
	expected := map[string]bool{
		filepath.Join(root, "a.jpg"):            true,
		filepath.Join(root, "b.png"):            true,
		filepath.Join(root, "nested", "d.webp"): true,
		filepath.Join(root, "nested", "e.JPEG"): true,
	}

	// Quick size check before verifying specific entries.
	if len(paths) != len(expected) {
		t.Fatalf("Images() returned %d files, want %d", len(paths), len(expected))
	}

	// Each returned path should be one of the expected image files.
	for _, p := range paths {
		if !expected[p] {
			t.Fatalf("unexpected path %q", p)
		}
	}
}
