# Developer Guide

## Architektur & Paketaufteilung

- `cmd/yearcollage/`: Einstieg, CLI-Flags, ruft `app.Run` auf.
- `internal/app/`: Orchestrierung (Input validieren, Images sammeln, sortieren, Aspect parsen, später Canvas/Output).
- `internal/aspect/`: Parsing von Seitenverhältnissen (`"3:2"` → `1.5`).
- `internal/collect/`: Rekursive Discovery erlaubter Bilddateien; Filter auf Extensions.
- Spätere Pakete: `internal/img` (load/crop/resize), `internal/collage` (Grid/Canvas/Save), optional `internal/exif`.

## Datenfluss (aktuell)

1. CLI-Flags werden in `Config` geschrieben und validiert.
2. Input-Ordner wird ge-`stat`-et, dann werden Bilder rekursiv gesammelt.
3. Sortierung nach `ModTime` (ältestes zuerst); Stat-Fehler werden geloggt, aber übersprungen.
4. Aspect-String wird geparst (noch ohne Anwendung auf Bilder).
5. TODO: Crop/Resize, Grid berechnen, Canvas zeichnen, Datei speichern.

## Build- & Test-Workflow

- `make build` — baut das Binary nach `bin/yearcollage`.
- `make run` — Beispielausführung mit Platzhaltern (`./bilder` anpassen).
- `make fmt` — `gofmt` auf `cmd/` und `internal/`.
- `make vet` — statische Checks.
- `make test` — führt die Go-Tests aus.

## Coding-Guidelines

- Flags lowercase/kebab (`-tile-aspect`). Export nur, was außerhalb des Pakets gebraucht wird.
- Früh validieren und explizit loggen; Fehler möglichst mit `%w` wrappen.
- Table-driven Tests in `*_test.go`; Fixtures klein halten (`testdata/`).
- Kommentare nur, wenn Verhalten nicht offensichtlich ist (z. B. Sortierung nach `ModTime` statt EXIF).

## Nächste sinnvolle Schritte

- Bild-Pipeline ergänzen (Aspect anwenden, Crop/Resize, Canvas, JPEG/PNG-Speichern) in separaten Paketen.
- Optional EXIF-Auswertung für korrektere Chronologie.
- `golangci-lint` ins Makefile hängen, sobald installiert.
- CLI-Doku und README synchron halten, wenn Flags/Defaults ändern.
