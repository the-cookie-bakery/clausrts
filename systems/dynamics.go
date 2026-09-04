package systems

import (
	"math"

	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

// Clausewitzian radii, in tiles.
const (
	commandRadius     = 12.0
	reconRadius       = 9.0
	jamRadius         = 9.0
	focusRadius       = 4.0
	culminationRadius = 22.0
	fogSampleRadius   = 5.0
)

// UpdateDynamics recomputes every unit's Clausewitz properties: friction, fog,
// focus, gravity, supply and culmination.
func UpdateDynamics(w *worlds.World) {
	for _, u := range w.Units {
		if !u.IsAlive() {
			continue
		}
		k := u.KindDef()
		d := &u.Dynamics

		// Fog: uncertainty about the immediate surroundings (Nachrichten).
		d.Fog = fogAround(w, u.Faction, u.Pos, fogSampleRadius)

		// Friction: base + terrain + command distance + jamming - recon + fog.
		fric := k.FrictionBase + w.Grid.FrictionMod(u.Pos)
		fric += nearestCommandDistance(w, u) / commandRadius * 0.5
		fric += enemyJam(w, u)
		fric -= friendlyRecon(w, u)
		fric += d.Fog * 0.4
		d.Friction = helpers.ClampF(fric, 0, 1)

		// Focus: concentration of friendly forces (Ökonomie der Kräfte) plus
		// recon/GPS support.
		friends := countFriendsNear(w, u, focusRadius)
		d.Focus = helpers.ClampF(float64(friends)/6.0+friendlyRecon(w, u)*0.3, 0, 1)

		// Gravity: the unit's Schwerpunkt weight (scales with health).
		d.Gravity = k.GravityWeight * (u.HP / u.MaxHP)

		// Supply & culmination: Verbindungslinien back to a supply source.
		dist := nearestSupplyDistance(w, u)
		d.Culmination = helpers.ClampF(dist/culminationRadius, 0, 1)
		d.Supply = 1 - d.Culmination
	}
}

// fogAround returns the fraction of nearby tiles not currently visible to the
// given faction.
func fogAround(w *worlds.World, faction entities.FactionID, pos helpers.Point, radius float64) float64 {
	fog := w.Fog[faction]
	if fog == nil {
		return 0
	}
	r := int(radius)
	seen, total := 0, 0
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if float64(dx*dx+dy*dy) > radius*radius {
				continue
			}
			total++
			if fog.At(pos.Add(helpers.Point{X: dx, Y: dy})) == worlds.VisVisible {
				seen++
			}
		}
	}
	if total == 0 {
		return 0
	}
	return 1 - float64(seen)/float64(total)
}

// nearestCommandDistance returns the distance in tiles to the nearest friendly
// HQ or radar; large if none (i.e. friction grows without command).
func nearestCommandDistance(w *worlds.World, u *entities.Unit) float64 {
	best := commandRadius * 4
	for _, b := range w.BuildingsOf(u.Faction) {
		k := b.KindDef()
		if !k.IsHQ && !k.IsRadar {
			continue
		}
		if d := math.Sqrt(worlds.SquareDistance(u.Pos, b.Pos)); d < best {
			best = d
		}
	}
	return best
}

// enemyJam returns the strongest enemy jamming effect acting on the unit.
func enemyJam(w *worlds.World, u *entities.Unit) float64 {
	best := 0.0
	for _, e := range w.Units {
		if e.Faction == u.Faction || !e.IsAlive() {
			continue
		}
		if e.KindDef().JamPower <= 0 {
			continue
		}
		d := math.Sqrt(worlds.SquareDistance(u.Pos, e.Pos))
		if d > jamRadius {
			continue
		}
		effect := e.KindDef().JamPower * (1 - d/jamRadius)
		if effect > best {
			best = effect
		}
	}
	return best
}

// friendlyRecon returns the strongest friendly reconnaissance effect on the unit.
func friendlyRecon(w *worlds.World, u *entities.Unit) float64 {
	best := 0.0
	for _, f := range w.Units {
		if f.Faction != u.Faction || !f.IsAlive() {
			continue
		}
		if f.KindDef().ReconPower <= 0 {
			continue
		}
		d := math.Sqrt(worlds.SquareDistance(u.Pos, f.Pos))
		if d > reconRadius {
			continue
		}
		effect := f.KindDef().ReconPower * (1 - d/reconRadius)
		if effect > best {
			best = effect
		}
	}
	return best
}

// countFriendsNear counts friendly units within the given radius.
func countFriendsNear(w *worlds.World, u *entities.Unit, radius float64) int {
	n := 0
	for _, f := range w.Units {
		if f.ID == u.ID || f.Faction != u.Faction || !f.IsAlive() {
			continue
		}
		if math.Sqrt(worlds.SquareDistance(u.Pos, f.Pos)) <= radius {
			n++
		}
	}
	return n
}

// nearestSupplyDistance returns the distance to the nearest friendly supply
// depot, falling back to the faction's HQ, then the map edge.
func nearestSupplyDistance(w *worlds.World, u *entities.Unit) float64 {
	best := math.Inf(1)
	for _, b := range w.BuildingsOf(u.Faction) {
		if b.KindDef().IsSupply {
			best = math.Min(best, math.Sqrt(worlds.SquareDistance(u.Pos, b.Pos)))
		}
	}
	if !math.IsInf(best, 1) {
		return best
	}
	for _, b := range w.BuildingsOf(u.Faction) {
		if b.KindDef().IsHQ {
			best = math.Min(best, math.Sqrt(worlds.SquareDistance(u.Pos, b.Pos)))
		}
	}
	if !math.IsInf(best, 1) {
		return best
	}
	return 0
}
