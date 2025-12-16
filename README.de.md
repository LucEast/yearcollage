# YearCollage CLI (Deutsch)

<p align="center">
  <a href="https://github.com/LucEast/obsidian-current-view/releases">
    <img src="https://img.shields.io/github/v/release/LucEast/yearcollage?style=for-the-badge&label=latest&labelColor=363a4f&color=B4BEFE&logo=github&logoColor=cad3f5" alt="GitHub Release" />
  </a>
  <a href="https://github.com/LucEast/obsidian-current-view/releases">
    <img src="https://img.shields.io/github/downloads/LucEast/yearcollage/total?style=for-the-badge&label=downloads&labelColor=363a4f&color=F9E2AF&logo=abdownloadmanager&logoColor=cad3f5" alt="Downloads" />
  </a>
  <a href="https://github.com/LucEast/obsidian-current-view/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/LucEast/yearcollage/semantic-release.yml?branch=main&style=for-the-badge&label=CI&labelColor=363a4f&color=A6E3A1&logo=githubactions&logoColor=cad3f5" alt="CI Status" />
  </a>
  <a href="https://github.com/LucEast/obsidian-current-view/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/LucEast/yearcollage?style=for-the-badge&labelColor=363a4f&color=FAB387&logo=open-source-initiative&logoColor=cad3f5" alt="License" />
  </a>
</p>

[Deutsch](README.de.md) | [English](README.md)

YearCollage ist ein Go-CLI-Tool, das einen Foto-Ordner zu einer einzigen Collage zusammenstellt. Bilder werden rekursiv gesammelt, sortiert, auf ein einheitliches Seitenverhaeltnis gecroppt, skaliert und in einem Grid angeordnet.

## Installation
- Go 1.20+ empfohlen
- Bauen: `go build -o yearcollage .`
- Direkt laufen lassen: `go run . -input ./bilder -output collage.jpg`

## Schnellstart
```bash
yearcollage -input ./bilder -output collage.jpg
```
Die Collage wird dort abgelegt, wo du den Befehl ausfuehrst, ausser du gibst einen absoluten oder anderen relativen `-output` Pfad an.

## Flags
| Flag (Kurz) | Default | Beschreibung |
| --- | --- | --- |
| `-input`, `-i` | _required_ | Verzeichnis fuer Bilder (rekursiv). |
| `-output`, `-o` | `collage.jpg` | Ausgabedatei (Endung steuert JPEG/PNG). |
| `-tile-aspect`, `-a` | `1:1` | Seitenverhaeltnis pro Kachel (wird ignoriert, wenn `-collage-aspect` gesetzt ist). |
| `-tile-width`, `-w` | `400` | Kachelbreite in Pixeln; Hoehe wird vom Seitenverhaeltnis abgeleitet. |
| `-columns`, `-c` | `20` | Spaltenanzahl (ignoriert, wenn `-collage-aspect` gesetzt ist). |
| `-collage-aspect`, `-r` | _leer_ | Ziel-Seitenverhaeltnis der gesamten Collage; Spalten und Kachel-Aspect werden automatisch bestimmt. |
| `-sort`, `-s` | `time` | Sortierung: `time` (Dateizeit), `name` (alphabetisch), `exif` (EXIF DateTime*). |

\* Bei `-sort exif` werden DateTimeOriginal/DateTimeDigitized/DateTime gelesen; faellt auf Dateizeit zurueck, wenn nicht vorhanden.

Unterstuetzte Eingaben: `.jpg`, `.jpeg`, `.png`, `.webp`.

## Beispiele
- Fixes Grid: `yearcollage -i ./bilder/2025 -o collage-2025.jpg -c 18 -w 360 -a 3:2`
- Spalten automatisch ueber Collage-Aspect: `yearcollage -i ./urlaub -o collage-urlaub.png -collage-aspect 16:9 -w 320`
- Chronologisch nach EXIF: `yearcollage -i ./bilder -sort exif`

## Hinweise
- Wenn `-collage-aspect` gesetzt ist, wird `-tile-aspect` ignoriert; ein passender Tile-Aspect wird abgeleitet.
- Layout: links→rechts, oben→unten.
- Ausgabeformat: PNG bei `.png`, sonst JPEG (Qualitaet 90).

## Entwicklung
- Formatierung: `gofmt -w .`
- Checks: `go vet ./...`
- Tests: `go test ./...`
- Build: `go build -o yearcollage .`
