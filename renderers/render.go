package renderers

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"clausrts/assets"
	"clausrts/configs"
	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/uis"
	"clausrts/worlds"
)

const tile = float64(configs.TileSize)

var fogExploredAlpha = float32(0.45)

// Draw renders the whole world for the given observing faction.
func Draw(screen *ebiten.Image, w *worlds.World, cam *Camera, faction entities.FactionID, selected map[int]bool) {
	drawTiles(screen, w, cam, faction)
	drawBuildings(screen, w, cam, faction)
	drawUnits(screen, w, cam, faction, selected)
}

func drawTiles(screen *ebiten.Image, w *worlds.World, cam *Camera, faction entities.FactionID) {
	fog := w.Fog[faction]
	for y := 0; y < w.Grid.Height; y++ {
		for x := 0; x < w.Grid.Width; x++ {
			p := helpers.Point{X: x, Y: y}
			sx, sy := cam.TileToScreen(p)
			if sx+tile < 0 || sx > float64(cam.ScreenW) || sy+tile < 0 || sy > float64(cam.ScreenH) {
				continue
			}

			vis := byte(worlds.VisVisible)
			if fog != nil {
				vis = fog.At(p)
			}
			if vis == worlds.VisHidden {
				continue
			}

			img := assets.TerrainImage(terrainName(w.Grid.Tiles[y][x].Terrain))
			if img == nil {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(sx, sy)
			screen.DrawImage(img, op)

			if vis == worlds.VisExplored {
				drawAlpha(screen, sx, sy, fogExploredAlpha)
			}
		}
	}
}

func drawBuildings(screen *ebiten.Image, w *worlds.World, cam *Camera, faction entities.FactionID) {
	fog := w.Fog[faction]
	for _, b := range w.Buildings {
		if !b.IsAlive() {
			continue
		}
		// Hide enemy buildings the observer cannot currently see.
		if b.Faction != faction && fog != nil && fog.At(b.Pos) != worlds.VisVisible {
			continue
		}
		img := assets.Image(b.KindDef().Asset)
		if img == nil {
			continue
		}
		sx, sy := cam.TileToScreen(b.Pos)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx, sy)
		screen.DrawImage(img, op)
		drawFactionMarker(screen, b.Faction, sx, sy)
		if b.HP < b.MaxHP {
			drawHealthBar(screen, sx, sy, b.HP, b.MaxHP)
		}
	}
}

func drawUnits(screen *ebiten.Image, w *worlds.World, cam *Camera, faction entities.FactionID, selected map[int]bool) {
	// Painter's order: draw by row so lower units appear in front.
	units := make([]*entities.Unit, 0, len(w.Units))
	for _, u := range w.Units {
		if u.IsAlive() {
			units = append(units, u)
		}
	}

	for _, u := range units {
		fog := w.Fog[faction]
		if u.Faction != faction && fog != nil && fog.At(u.Pos) != worlds.VisVisible {
			continue
		}

		sx, sy := cam.WorldToScreen(u.RenderX, u.RenderY)
		if sx < -tile || sx > float64(cam.ScreenW)+tile || sy < -tile || sy > float64(cam.ScreenH)+tile {
			continue
		}

		img := assets.Image(u.KindDef().Asset)
		if img == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		// Rotate the image toward the facing tile. Per-kind customization
		// (e.g. a tank rotating only its turret) replaces this block later.
		angle := facingAngle(u.Facing)
		op.GeoM.Translate(-tile/2, -tile/2)
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(tile/2, tile/2)
		op.GeoM.Translate(sx-tile/2, sy-tile/2)
		screen.DrawImage(img, op)

		drawFactionMarker(screen, u.Faction, sx-tile/2, sy-tile/2)

		if selected[u.ID] {
			drawSelectionMarker(screen, sx, sy)
			drawAttackRadius(screen, sx, sy, u.KindDef().Range)
		}
		if u.HP < u.MaxHP {
			drawHealthBar(screen, sx-tile/2, sy-tile/2, u.HP, u.MaxHP)
		}
	}
}

func drawFactionMarker(screen *ebiten.Image, faction entities.FactionID, sx, sy float64) {
	f := factionColor(faction)
	vector.DrawFilledCircle(screen, float32(sx+tile/2), float32(sy+tile/2+14), 3, f, true)
}

// drawHealthBar draws a colour-coded bar and HP value above a damaged entity.
// (x, y) is the tile's top-left corner.
func drawHealthBar(screen *ebiten.Image, x, y, hp, maxHP float64) {
	if maxHP <= 0 || hp >= maxHP {
		return
	}
	frac := hp / maxHP
	if frac < 0 {
		frac = 0
	}

	const barH = 4.0
	bx, by := x, y-8

	vector.DrawFilledRect(screen, float32(bx), float32(by), float32(tile), float32(barH), color.RGBA{30, 30, 34, 235}, false)

	clr := color.RGBA{70, 200, 90, 255}
	switch {
	case frac < 0.25:
		clr = color.RGBA{225, 60, 60, 255}
	case frac < 0.5:
		clr = color.RGBA{235, 185, 45, 255}
	}
	vector.DrawFilledRect(screen, float32(bx), float32(by), float32(tile*frac), float32(barH), clr, false)

	label := fmt.Sprintf("%.0f", hp)
	tw := uis.TextWidth(label, 9)
	uis.DrawText(screen, label, bx+(tile-tw)/2, by-3, 9, color.RGBA{255, 255, 255, 255})
}

func drawSelectionMarker(screen *ebiten.Image, sx, sy float64) {
	vector.StrokeRect(screen, float32(sx-tile/2), float32(sy-tile/2), float32(tile), float32(tile), 2, color.RGBA{0, 220, 255, 255}, true)
}

// drawAttackRadius draws a faint filled circle (plus outline) showing the unit's
// attack distance. (sx, sy) is the unit's pixel center.
func drawAttackRadius(screen *ebiten.Image, sx, sy, rng float64) {
	if rng <= 0 {
		return
	}
	r := float32(rng * tile)
	cx, cy := float32(sx), float32(sy)
	vector.DrawFilledCircle(screen, cx, cy, r, color.RGBA{0, 220, 255, 12}, true)
	vector.StrokeCircle(screen, cx, cy, r, 1.5, color.RGBA{0, 220, 255, 90}, true)
}

// facingAngle returns the clockwise angle (radians) that rotates an
// up-pointing image onto the given tile-space direction.
func facingAngle(d helpers.Point) float64 {
	return math.Atan2(float64(d.X), -float64(d.Y))
}

func terrainName(t worlds.Terrain) string {
	switch t {
	case worlds.TerrainRoad:
		return "road"
	case worlds.TerrainForest:
		return "forest"
	case worlds.TerrainUrban:
		return "urban"
	case worlds.TerrainWater:
		return "water"
	case worlds.TerrainMountain:
		return "mountain"
	default:
		return "grass"
	}
}

func factionColor(id entities.FactionID) color.RGBA {
	switch id {
	case 1:
		return color.RGBA{230, 70, 70, 255}
	default:
		return color.RGBA{70, 140, 230, 255}
	}
}

// drawAlpha dims a tile by drawing a translucent black rectangle.
func drawAlpha(screen *ebiten.Image, sx, sy float64, alpha float32) {
	vector.DrawFilledRect(screen, float32(sx), float32(sy), float32(tile), float32(tile), color.RGBA{8, 8, 12, uint8(alpha * 255)}, false)
}
