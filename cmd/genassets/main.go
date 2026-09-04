// Command genassets regenerates the placeholder PNG assets used by the game.
// Assets are simple geometric shapes; replace the PNGs with real art later —
// no code changes are required because loading goes through assets.Image.
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

const size = 32

type rgba struct{ r, g, b, a uint8 }

func img() *image.RGBA { return image.NewRGBA(image.Rect(0, 0, size, size)) }

func fillCircle(m *image.RGBA, cx, cy, r int, c rgba) {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= r*r {
				m.SetRGBA(x, y, color.RGBA{c.r, c.g, c.b, c.a})
			}
		}
	}
}

func fillRect(m *image.RGBA, x0, y0, x1, y1 int, c rgba) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if x >= 0 && x < size && y >= 0 && y < size {
				m.SetRGBA(x, y, color.RGBA{c.r, c.g, c.b, c.a})
			}
		}
	}
}

func fillTriangle(m *image.RGBA, p0, p1, p2 [2]int, c rgba) {
	sign := func(a, b, p [2]int) int {
		return (a[0]-p[0])*(b[1]-p[1]) - (b[0]-p[0])*(a[1]-p[1])
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			p := [2]int{x, y}
			d1 := sign(p0, p1, p)
			d2 := sign(p1, p2, p)
			d3 := sign(p2, p0, p)
			neg := d1 < 0 || d2 < 0 || d3 < 0
			pos := d1 > 0 || d2 > 0 || d3 > 0
			if !(neg && pos) {
				m.SetRGBA(x, y, color.RGBA{c.r, c.g, c.b, c.a})
			}
		}
	}
}

func save(name string, m *image.RGBA) {
	path := filepath.Join("assets", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, m); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", path)
}

// unit draws a circle body with a white direction wedge pointing north.
func unit(name string, c rgba) {
	m := img()
	fillCircle(m, 16, 16, 12, c)
	// direction wedge (points up / north)
	fillTriangle(m, [2]int{16, 3}, [2]int{11, 12}, [2]int{21, 12}, rgba{255, 255, 255, 255})
	save(name, m)
}

// building draws a filled square with an inset lighter square.
func building(name string, c rgba) {
	m := img()
	fillRect(m, 4, 4, 28, 28, c)
	fillRect(m, 9, 9, 23, 23, rgba{uint8(c.r / 3 * 4), uint8(c.g / 3 * 4), uint8(c.b / 3 * 4), 255})
	save(name, m)
}

// terrain draws a flat tile.
func terrain(name string, c rgba) {
	m := img()
	fillRect(m, 0, 0, size, size, c)
	save("terrain/"+name, m)
}

func main() {
	units := map[string]rgba{
		"infantry.png":  {0, 180, 80, 255},
		"tank.png":      {90, 130, 60, 255},
		"artillery.png": {200, 120, 40, 255},
		"drone.png":     {90, 200, 220, 255},
		"cyber.png":     {180, 80, 200, 255},
		"scrambler.png": {220, 180, 60, 255},
		"gps.png":       {70, 140, 250, 255},
		"airdefense.png":{120, 120, 200, 255},
		"engineer.png":  {180, 140, 90, 255},
	}
	for name, c := range units {
		unit(name, c)
	}

	buildings := map[string]rgba{
		"hq.png":       {230, 60, 60, 255},
		"depot.png":    {200, 150, 60, 255},
		"radar.png":    {60, 140, 240, 255},
		"bunker.png":   {120, 120, 120, 255},
		"firewall.png": {60, 200, 140, 255},
	}
	for name, c := range buildings {
		building(name, c)
	}

	terrains := map[string]rgba{
		"grass.png":    {96, 170, 96, 255},
		"road.png":     {150, 150, 150, 255},
		"forest.png":   {40, 110, 55, 255},
		"urban.png":    {160, 140, 115, 255},
		"water.png":    {55, 110, 185, 255},
		"mountain.png": {115, 115, 125, 255},
	}
	for name, c := range terrains {
		terrain(name, c)
	}
}
