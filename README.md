# 📸 YearCollage – Create Large Image Collages Automatically

**YearCollage** ist ein Go‑basiertes CLI-Tool, mit dem du automatisch große Collagen aus Bildern generieren kannst – perfekt für **Jahresrückblicke**, **Poster**, **Fotowände** oder Social-Media-Projekte.

Das Tool:

* liest alle Bilder aus einem Ordner (rekursiv),
* sortiert sie **chronologisch nach Aufnahmedatum**,
* passt jedes Bild auf ein einheitliches Seitenverhältnis an (Cropping),
* ordnet alle Bilder als Grid an (links → rechts, oben → unten),
* erzeugt am Ende eine Collage als JPEG/PNG.

---

## 🚀 Features

* 📁 **Ordner einlesen** (inkl. Unterordner)
* 🕒 **Sortierung nach Dateidatum** (älteste zuerst)
* 🖼️ **Resize & Crop** auf festes Seitenverhältnis (z. B. 1:1, 3:2, 4:3)
* 🧱 **Grid-Platzierung** nach Spalten & Reihen
* 🖼️ **Output als Bilddatei** (z. B. `collage.jpg`)
* ⚙️ **Konfigurierbar über Flags**

---

## 📦 Installation

Du brauchst Go (Version 1.20 oder neuer).

```bash
git clone https://github.com/luceast/yearcollage
cd yearcollage
go build -o yearcollage
```

---

## 🛠️ Usage

### Minimal

```bash
yearcollage -input ./bilder
```

### Voller Befehl

```bash
yearcollage \
  -input ./bilder/2025 \   # Kurzform: -i
  -output collage-2025.jpg \ # Kurzform: -o
  -tile-aspect 3:2 \       # Kurzform: -a (wird ignoriert, wenn -collage-aspect gesetzt ist)
  -tile-width 400 \        # Kurzform: -w
  -columns 20              # Kurzform: -c (oder -collage-aspect 16:9 für Auto-Berechnung)
  -sort time               # Kurzform: -s (time | name)
```

### Parameter

| Flag (long/short) | Beschreibung                                         |
| ----------------- | ---------------------------------------------------- |
| `-input`, `-i`    | Pfad zum Bilder-Ordner (**required**)                |
| `-output`, `-o`   | Zieldatei für die Collage (default: `collage.jpg`)   |
| `-tile-aspect`, `-a` | Seitenverhältnis für jedes Bild (z. B. `1:1`, `3:2`) – wird ignoriert, wenn `-collage-aspect` genutzt wird |
| `-tile-width`, `-w`  | Breite jedes einzelnen Bildes im Grid                |
| `-columns`, `-c`     | Anzahl der Spalten im finalen Grid (ignoriert, wenn `-collage-aspect` gesetzt ist) |
| `-collage-aspect`, `-r` | Ziel-Seitenverhältnis der gesamten Collage; Spalten & Kachel-Seitenverhältnis werden automatisch berechnet |
| `-sort`, `-s`        | Sortierung der Bilder: `time` (Datei-Modtime, älteste zuerst), `name` (alphabetisch) oder `exif` (EXIF DateTime\*) |

\* Bei `-sort exif` wird auf EXIF-Tags (DateTimeOriginal/DateTimeDigitized/DateTime) zurückgegriffen, andernfalls auf Dateizeit.

---

## 🧠 Internes Funktionsprinzip

1. **Bilder finden:**

   * alle Dateien im Ordner sammeln
   * Endungen filtern (`jpg`, `jpeg`, `png`, `webp`)

2. **Nach Datum sortieren:**

   * kleinster Zeitstempel → erstes Bild
   * Ergebnis: Collage verläuft chronologisch

3. **Bilder verarbeiten:**

   * laden
   * auf Seitenverhältnis croppen
   * auf feste Breite skalieren

4. **Canvas erzeugen:**

   * Gesamtbreite = `columns * tileWidth`
   * Höhe ergibt sich dynamisch aus Anzahl der Bilder

5. **Bilder platzieren:**

   * Zeile für Zeile
   * Pixelgenau

6. **Als JPEG/PNG speichern**

---

## 📚 TODO / Next Steps

* [x] `tile-aspect` Parser implementieren
* [x] cropping-Funktion (`cropToAspect`)
* [x] resizing-Funktion
* [x] Canvas erstellen & Bilder zeichnen
* [x] Output speichern
* [x] optional: EXIF‑Date statt File‑Date nutzen
* [ ] optional: Rand & Abstand zwischen Kacheln einführen
* [ ] optional: Hintergrundfarbe wählbar machen

---

## 🤝 Contribution

PRs sind jederzeit willkommen – besonders beim Bild-Processing und bei Optimierungen für Performance.

---

## 📄 License

MIT License

---

Wenn du möchtest, können wir als Nächstes die README weiter strukturieren, Diagramme einbauen oder eine richtige Projektstruktur (`cmd/`, `pkg/`, `internal/`) anlegen. 😊
