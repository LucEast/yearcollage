package app

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"testing"
)

// TestNormalizeOrientation ensures each EXIF orientation code maps to the pixel
// arrangement described in the spec so the helpers stay stable.
func TestNormalizeOrientation(t *testing.T) {
	cases := []struct {
		name        string
		orientation int
		want        [][]int
	}{
		{
			name:        "orientation 1 (identity)",
			orientation: 1,
			want: [][]int{
				{0, 1},
				{2, 3},
				{4, 5},
			},
		},
		{
			name:        "orientation 2 (mirror horizontal)",
			orientation: 2,
			want: [][]int{
				{1, 0},
				{3, 2},
				{5, 4},
			},
		},
		{
			name:        "orientation 3 (rotate 180)",
			orientation: 3,
			want: [][]int{
				{5, 4},
				{3, 2},
				{1, 0},
			},
		},
		{
			name:        "orientation 4 (mirror vertical)",
			orientation: 4,
			want: [][]int{
				{4, 5},
				{2, 3},
				{0, 1},
			},
		},
		{
			name:        "orientation 5 (mirror horizontal + rotate 270 CW)",
			orientation: 5,
			want: [][]int{
				{0, 2, 4},
				{1, 3, 5},
			},
		},
		{
			name:        "orientation 6 (rotate 90 CW)",
			orientation: 6,
			want: [][]int{
				{4, 2, 0},
				{5, 3, 1},
			},
		},
		{
			name:        "orientation 7 (mirror horizontal + rotate 90 CW)",
			orientation: 7,
			want: [][]int{
				{5, 3, 1},
				{4, 2, 0},
			},
		},
		{
			name:        "orientation 8 (rotate 90 CCW)",
			orientation: 8,
			want: [][]int{
				{1, 3, 5},
				{0, 2, 4},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			img := makeLabeledImage()
			got := normalizeOrientation(img, tc.orientation)
			if diff := compareGrid(got, tc.want); diff != "" {
				t.Fatalf("normalizeOrientation mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func makeLabeledImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 2, 3))
	val := 0
	for y := 0; y < 3; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.NRGBA{uint8(val), 0, 0, 255})
			val++
		}
	}
	return img
}

func compareGrid(img image.Image, want [][]int) string {
	b := img.Bounds()
	got := make([][]int, b.Dy())
	for y := 0; y < b.Dy(); y++ {
		row := make([]int, b.Dx())
		for x := 0; x < b.Dx(); x++ {
			r, _, _, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			row[x] = int(r >> 8)
		}
		got[y] = row
	}
	if reflect.DeepEqual(got, want) {
		return ""
	}
	return fmt.Sprintf("got %v want %v", got, want)
}
