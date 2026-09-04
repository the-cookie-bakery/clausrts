// Package uis renders the heads-up display: the Trinity bars, minimap,
// selection panel and status text.
package uis

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

var (
	clrPeople     = color.RGBA{240, 200, 60, 255}
	clrMilitary   = color.RGBA{70, 150, 240, 255}
	clrGovernment = color.RGBA{200, 90, 200, 255}
	clrText       = color.RGBA{220, 220, 220, 255}
	clrPanel      = color.RGBA{20, 22, 30, 200}
	clrBarBG      = color.RGBA{40, 42, 52, 255}
)

// DrawHUD renders the complete interface.
func DrawHUD(screen *ebiten.Image, w *worlds.World, player entities.FactionID, selected map[int]bool, paused bool, sw, sh int) {
	drawTrinity(screen, w, player, sw)
	drawSelectionPanel(screen, w, player, selected, sh)
	drawStatus(screen, paused, w.Tick, sw, sh)
}

func drawTrinity(screen *ebiten.Image, w *worlds.World, player entities.FactionID, sw int) {
	f := w.Faction(player)
	if f == nil {
		return
	}
	x, y := 12.0, 12.0
	const bw, bh = 220.0, 16.0

	DrawText(screen, f.Name, x, y, 16, clrText)
	y += 22

	drawBar(screen, x, y, bw, bh, f.People, clrPeople, "People (Volk)")
	y += bh + 8
	drawBar(screen, x, y, bw, bh, f.Military, clrMilitary, "Military (Heer)")
	y += bh + 8
	drawBar(screen, x, y, bw, bh, f.Government, clrGovernment, "Government (Regierung)")

	// Enemy bars, collapsed, top-right.
	enemy := enemyFaction(w, player)
	if enemy != nil {
		ex := float64(sw) - bw - 12
		ey := 12.0
		DrawText(screen, enemy.Name, ex, ey, 16, clrText)
		ey += 22
		drawBar(screen, ex, ey, bw, bh, enemy.People, clrPeople, "People")
		ey += bh + 8
		drawBar(screen, ex, ey, bw, bh, enemy.Military, clrMilitary, "Military")
		ey += bh + 8
		drawBar(screen, ex, ey, bw, bh, enemy.Government, clrGovernment, "Government")
	}
}

func drawBar(screen *ebiten.Image, x, y, w, h, frac float64, clr color.Color, label string) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), clrBarBG, false)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w*helpers.ClampF(frac, 0, 1)), float32(h), clr, false)
	DrawText(screen, fmt.Sprintf("%s %.0f%%", label, frac*100), x, y+h+12, 12, clrText)
}

func drawSelectionPanel(screen *ebiten.Image, w *worlds.World, player entities.FactionID, selected map[int]bool, sh int) {
	if len(selected) == 0 {
		return
	}
	// Show details for the first selected unit.
	var first *entities.Unit
	for id := range selected {
		if u := w.Units[id]; u != nil && u.IsAlive() {
			first = u
			break
		}
	}
	if first == nil {
		return
	}

	const pw, ph = 240.0, 150.0
	x, y := 12.0, float64(sh)-ph-12
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(pw), float32(ph), clrPanel, false)

	k := first.KindDef()
	DrawText(screen, k.Name, x+10, y+20, 16, clrText)
	DrawText(screen, fmt.Sprintf("HP %.0f / %.0f", first.HP, first.MaxHP), x+10, y+42, 13, clrText)
	drawBar(screen, x+10, y+48, pw-20, 8, first.HP/first.MaxHP, clrMilitary, "")

	d := first.Dynamics
	row := func(label string, v float64, yy float64) {
		DrawText(screen, label, x+10, yy, 12, clrText)
		DrawText(screen, fmt.Sprintf("%.2f", v), x+pw-60, yy, 12, clrText)
	}
	row("Friction", d.Friction, y+72)
	row("Focus", d.Focus, y+88)
	row("Fog", d.Fog, y+104)
	row("Moral", d.Moral, y+120)
	row("Culmination", d.Culmination, y+136)
}

func drawStatus(screen *ebiten.Image, paused bool, tick uint64, sw, sh int) {
	status := fmt.Sprintf("tick %d", tick)
	if paused {
		status += "  [PAUSED]"
	}
	DrawText(screen, status, float64(sw)-140, float64(sh)-20, 13, clrText)
}

func enemyFaction(w *worlds.World, player entities.FactionID) *entities.Faction {
	for id, f := range w.Factions {
		if id != player {
			return f
		}
	}
	return nil
}
