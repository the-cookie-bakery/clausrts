package systems

import (
	"clausrts/helpers"
	"clausrts/worlds"
)

// UpdateFog recomputes each faction's visibility grid from the sight radius of
// its units and buildings. Previously-seen tiles degrade to "explored".
func UpdateFog(w *worlds.World) {
	for fid, fog := range w.Fog {
		// Demote currently-visible tiles to explored.
		for y := 0; y < fog.H; y++ {
			for x := 0; x < fog.W; x++ {
				if fog.Values[y][x] == worlds.VisVisible {
					fog.Values[y][x] = worlds.VisExplored
				}
			}
		}
		for _, u := range w.Units {
			if u.Faction == fid && u.IsAlive() {
				reveal(fog, u.Pos, u.KindDef().Sight)
			}
		}
		for _, b := range w.Buildings {
			if b.Faction == fid && b.IsAlive() {
				reveal(fog, b.Pos, b.KindDef().Sight)
			}
		}
	}
}

// reveal marks every tile within the (Euclidean) sight radius as visible.
func reveal(fog *worlds.Fog, center helpers.Point, radius float64) {
	r := int(radius)
	r2 := radius * radius
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if float64(dx*dx+dy*dy) > r2 {
				continue
			}
			fog.Set(center.Add(helpers.Point{X: dx, Y: dy}), worlds.VisVisible)
		}
	}
}
