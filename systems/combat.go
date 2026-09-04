package systems

import (
	"math"

	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

// target is a unit or building that can receive damage.
type target struct {
	id         int
	pos        helpers.Point
	faction    entities.FactionID
	hp         *float64
	maxHP      float64
	kind       string
	isBuilding bool
}

// UpdateCombat resolves auto-attacks: units engage the nearest enemy in range.
// Friction reduces accuracy (misses), focus multiplies damage, cover and
// friendly bunkers reduce incoming damage.
func UpdateCombat(w *worlds.World, dt float64) {
	// Morale recovery when not under immediate stress.
	for _, u := range w.Units {
		if u.IsAlive() {
			u.Dynamics.Moral = helpers.ClampF(u.Dynamics.Moral+dt*0.02, 0, 1)
		}
	}

	for _, u := range w.Units {
		if !u.IsAlive() || u.Moving {
			continue
		}
		k := u.KindDef()
		if k.Attack <= 0 {
			continue
		}
		u.AttackCooldown -= dt
		if u.AttackCooldown > 0 {
			continue
		}
		t := nearestEnemy(w, u, k.Range)
		if t == nil {
			continue
		}
		u.Facing = u.Pos.DirectionTo(t.pos)
		u.AttackCooldown = attackInterval(k)

		// Friction causes misses: real war is not war on paper.
		if u.Dynamics.Friction > 0 && w.RNG.Float64() < u.Dynamics.Friction*0.5 {
			continue
		}

		damage := k.Attack * counterMult(u, t)
		damage *= 0.7 + 0.6*u.Dynamics.Focus // concentrated forces hit harder
		damage *= 1 - w.Grid.Cover(t.pos)    // terrain cover
		if t.isBuilding {
			damage *= 1 - bunkerShield(w, t) // bunkers shield against fire
		}

		applyDamage(w, u, t, damage)
	}
}

func attackInterval(k entities.UnitKind) float64 {
	if k.AttackInterval <= 0 {
		return 1.0
	}
	return k.AttackInterval
}

// nearestEnemy returns the closest enemy unit or building within range, or nil.
func nearestEnemy(w *worlds.World, u *entities.Unit, rng float64) *target {
	bestDist := rng * rng
	var best *target
	for id, e := range w.Units {
		if e.Faction == u.Faction || !e.IsAlive() {
			continue
		}
		if d := worlds.SquareDistance(u.Pos, e.Pos); d <= bestDist {
			bestDist = d
			best = &target{id: id, pos: e.Pos, faction: e.Faction, hp: &e.HP, maxHP: e.MaxHP, kind: e.Kind}
		}
	}
	for id, b := range w.Buildings {
		if b.Faction == u.Faction || !b.IsAlive() {
			continue
		}
		if d := worlds.SquareDistance(u.Pos, b.Pos); d <= bestDist {
			bestDist = d
			best = &target{id: id, pos: b.Pos, faction: b.Faction, hp: &b.HP, maxHP: b.MaxHP, kind: b.Kind, isBuilding: true}
		}
	}
	return best
}

// findTargetByID returns a target wrapper for a unit or building ID, or nil if
// it is gone.
func findTargetByID(w *worlds.World, id int) *target {
	if u, ok := w.Units[id]; ok && u.IsAlive() {
		return &target{id: id, pos: u.Pos, hp: &u.HP, maxHP: u.MaxHP, kind: u.Kind}
	}
	if b, ok := w.Buildings[id]; ok && b.IsAlive() {
		return &target{id: id, pos: b.Pos, hp: &b.HP, maxHP: b.MaxHP, kind: b.Kind, isBuilding: true}
	}
	return nil
}

// counterMult returns the damage multiplier for the counter matrix
// (Wechselwirkung: every weapon has a counter).
func counterMult(u *entities.Unit, t *target) float64 {
	for _, c := range u.KindDef().Counters {
		if c == t.kind {
			return 2.0
		}
	}
	return 1.0
}

// bunkerShield returns the damage reduction from a friendly bunker on the tile.
func bunkerShield(w *worlds.World, t *target) float64 {
	for _, b := range w.Buildings {
		if b.KindDef().IsBunker && b.IsAlive() && b.Pos.Equal(t.pos) {
			return 0.7
		}
	}
	return 0
}

// applyDamage mutates the target's HP and handles deaths.
func applyDamage(w *worlds.World, attacker *entities.Unit, t *target, damage float64) {
	*t.hp -= damage
	if *t.hp > 0 {
		if !t.isBuilding {
			if u := w.Units[t.id]; u != nil {
				u.Dynamics.Moral = helpers.ClampF(u.Dynamics.Moral-damage/u.MaxHP*0.5, 0, 1)
			}
		}
		return
	}

	if t.isBuilding {
		w.RemoveBuilding(t.id)
		return
	}

	if f := w.Factions[t.faction]; f != nil {
		f.Losses++
	}
	pos := w.Units[t.id].Pos
	w.RemoveUnit(t.id)

	// Nearby friendly units lose morale when a comrade falls.
	for _, u := range w.Units {
		if u.Faction == t.faction && u.IsAlive() &&
			math.Sqrt(worlds.SquareDistance(u.Pos, pos)) <= 3 {
			u.Dynamics.Moral = helpers.ClampF(u.Dynamics.Moral-0.15, 0, 1)
		}
	}
}
