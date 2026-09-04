package paths

import (
	"math"

	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

// ThreatExtras returns a CostExtras that adds enemy gravity (Schwerpunkt
// influence) as a routing cost. Units will observably route around enemy
// concentrations, teaching the center-of-gravity concept.
func ThreatExtras(w *worlds.World, faction entities.FactionID, radius, weight float64) CostExtras {
	return func(p helpers.Point) float64 {
		var c float64
		for _, u := range w.Units {
			if u.Faction == faction || !u.IsAlive() {
				continue
			}
			d := math.Sqrt(worlds.SquareDistance(u.Pos, p))
			if d > radius {
				continue
			}
			c += u.Dynamics.Gravity / (d + 1) * weight
		}
		return c
	}
}
