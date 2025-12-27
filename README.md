# YearCollage CLI

<p align="center">
  <a href="https://github.com/LucEast/yearcollage/releases">
    <img src="https://img.shields.io/github/v/release/LucEast/yearcollage?style=for-the-badge&label=latest&labelColor=363a4f&color=B4BEFE&logo=github&logoColor=cad3f5" alt="GitHub Release" />
  </a>
  <a href="https://github.com/LucEast/yearcollage/releases">
    <img src="https://img.shields.io/github/downloads/LucEast/yearcollage/total?style=for-the-badge&label=downloads&labelColor=363a4f&color=F9E2AF&logo=abdownloadmanager&logoColor=cad3f5" alt="Downloads" />
  </a>
  <a href="https://github.com/LucEast/yearcollage/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/LucEast/yearcollage/semantic-release.yml?branch=main&style=for-the-badge&label=CI&labelColor=363a4f&color=A6E3A1&logo=githubactions&logoColor=cad3f5" alt="CI Status" />
  </a>
  <a href="https://github.com/LucEast/yearcollage/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/LucEast/yearcollage?style=for-the-badge&labelColor=363a4f&color=FAB387&logo=open-source-initiative&logoColor=cad3f5" alt="License" />
  </a>
</p>

English | [Deutsch](README.de.md)

YearCollage is a Go CLI that turns a folder of photos into a single collage image. Images are collected recursively, sorted, cropped to a common aspect, scaled, and laid out in a grid.

## Install
- Go 1.20+ recommended
- Build: `go build -o yearcollage .`
- Run without building: `go run . -input ./photos -output collage.jpg`

## Quickstart
```bash
yearcollage -input ./bilder -output collage.jpg
```
The collage is written where you run the command unless you give an absolute or different relative `-output` path.

## Flags
| Flag (short) | Default | Description |
| --- | --- | --- |
| `-input`, `-i` | _required_ | Directory to scan for images (recursive). |
| `-output`, `-o` | `collage.jpg` | Output file path (extension controls JPEG/PNG). |
| `-tile-aspect`, `-a` | `1:1` | Aspect ratio for each tile (ignored if `-collage-aspect` is set). |
| `-tile-width`, `-w` | `400` | Tile width in pixels. Height is derived from aspect. |
| `-columns`, `-c` | `20` | Columns in the grid (ignored if `-collage-aspect` is set). |
| `-collage-aspect`, `-r` | _empty_ | Target aspect ratio for the whole collage; auto-picks columns and tile aspect. |
| `-sort`, `-s` | `time` | Sort mode: `time` (file mod time), `name` (alphabetical), `exif` (EXIF DateTime*). |

\* For `-sort exif`, EXIF DateTimeOriginal/DateTimeDigitized/DateTime are tried; falls back to file mod time if missing.

Supported inputs: `.jpg`, `.jpeg`, `.png`, `.webp`.

## Examples
- Fixed grid: `yearcollage -i ./bilder/2025 -o collage-2025.jpg -c 18 -w 360 -a 3:2`
- Auto columns by collage ratio: `yearcollage -i ./bilder/urlaub -o collage-urlaub.png -collage-aspect 16:9 -w 320`
- EXIF chronological: `yearcollage -i ./bilder -sort exif`

## Notes
- If you set `-collage-aspect`, the provided `-tile-aspect` is ignored; a tile aspect is derived to fit the target collage ratio.
- Images are laid out left→right, top→bottom.
- Output format: PNG if `-output` ends with `.png`, otherwise JPEG (quality 90).

## Development
- Format: `gofmt -w .`
- Lint: `go vet ./...`
- Tests: `go test ./...`
- Build: `go build -o yearcollage .`

## Translations
The default README is English. Add a translated copy like `README.de.md` and link it near the top, e.g.:
```
English | [Deutsch](README.de.md)
```
