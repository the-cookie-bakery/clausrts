package systems

import (
	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

// UpdateTrinity aggregates the three Clausewitz Trinity bars per faction and
// recomputes each faction's center of gravity (Schwerpunkt).
func UpdateTrinity(w *worlds.World, dt float64) {
	for _, f := range w.Factions {
		// Military (Heer): the remaining fraction of fighting strength.
		var hp, maxHP float64
		for _, u := range w.UnitsOf(f.ID) {
			hp += u.HP
			maxHP += u.MaxHP
		}
		if maxHP > 0 {
			f.Military = hp / maxHP
		} else {
			f.Military = 0
		}

		// People (Volk): war-weariness from time and casualties, slowly
		// recovering toward unity in quiet periods.
		f.People -= dt * 0.0005
		f.People -= float64(f.Losses) * 0.0001
		f.People += (1 - f.People) * dt * 0.002
		f.People = helpers.ClampF(f.People, 0, 1)

		// Government (Regierung): political will tracks the integrity of the
		// seat of policy (HQ) — the war is an instrument of politics.
		if hq := ownHQ(w, f.ID); hq != nil {
			target := 0.5 + 0.5*(hq.HP/hq.MaxHP)
			f.Government += (target - f.Government) * 0.02
		}
		f.Government = helpers.ClampF(f.Government, 0, 1)

		// Center of gravity: gravity-weighted mean position of all forces.
		var sx, sy, g float64
		for _, u := range w.UnitsOf(f.ID) {
			gx := u.Dynamics.Gravity
			sx += float64(u.Pos.X) * gx
			sy += float64(u.Pos.Y) * gx
			g += gx
		}
		for _, b := range w.BuildingsOf(f.ID) {
			gx := b.KindDef().GravityWeight * (b.HP / b.MaxHP)
			sx += float64(b.Pos.X) * gx
			sy += float64(b.Pos.Y) * gx
			g += gx
		}
		if g > 0 {
			f.Schwerpunkt = helpers.Point{X: int(sx/g + 0.5), Y: int(sy/g + 0.5)}
		}
		f.TotalGravity = g
	}
}

// ownHQ returns the faction's headquarters building, if any.
func ownHQ(w *worlds.World, fid entities.FactionID) *entities.Building {
	for _, b := range w.BuildingsOf(fid) {
		if b.KindDef().IsHQ {
			return b
		}
	}
	return nil
}
