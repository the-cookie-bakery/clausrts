// Package ai provides the enemy commander. It is deliberately simple: a
// Clausewitzian commander seeks the enemy center of gravity (Schwerpunkt) and
// advances toward it, which is exactly the behaviour the game teaches players
// to do deliberately.
package ai

import (
	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/paths"
	"clausrts/worlds"
)

// EnemyCommander controls one faction.
type EnemyCommander struct {
	Faction entities.FactionID
}

// NewEnemyCommander returns a commander for the given faction.
func NewEnemyCommander(faction entities.FactionID) *EnemyCommander {
	return &EnemyCommander{Faction: faction}
}

// Step re-issues orders to idle units, advancing toward the enemy's center of
// gravity. Orders are re-evaluated periodically to avoid churn.
func (c *EnemyCommander) Step(w *worlds.World) {
	if w.Tick%90 != 0 {
		return
	}
	opponent := opponentFaction(w, c.Faction)
	if opponent == nil {
		return
	}
	target := opponent.Schwerpunkt
	if target == (helpers.Point{}) {
		return
	}

	for _, u := range w.UnitsOf(c.Faction) {
		// Leave units with an active move/attack order alone.
		if u.Order.Kind == entities.OrderMove || u.Order.Kind == entities.OrderAttack {
			continue
		}
		if u.Order.PathIdx < len(u.Order.Path) {
			continue
		}

		// Prefer striking a nearby enemy; otherwise advance toward Schwerpunkt.
		if t := nearestEnemyUnit(w, u, u.KindDef().Range+1); t != 0 {
			u.Order.Kind = entities.OrderAttack
			u.Order.TargetUnit = t
			u.Order.Path = nil
			u.Order.PathIdx = 0
			continue
		}

		extras := paths.ThreatExtras(w, c.Faction, 8.0, 4.0)
		path := paths.FindPath(w.Grid, u.Pos, target, u.KindDef().CanFly, extras)
		if path == nil {
			continue
		}
		u.Order.Kind = entities.OrderMove
		u.Order.Target = target
		u.Order.Path = path
		u.Order.PathIdx = 0
	}
}

func opponentFaction(w *worlds.World, self entities.FactionID) *entities.Faction {
	for id, f := range w.Factions {
		if id != self {
			return f
		}
	}
	return nil
}

func nearestEnemyUnit(w *worlds.World, u *entities.Unit, rng float64) int {
	best := 0
	bestDist := rng * rng
	for id, e := range w.Units {
		if e.Faction == u.Faction || !e.IsAlive() {
			continue
		}
		if d := worlds.SquareDistance(u.Pos, e.Pos); d <= bestDist {
			bestDist = d
			best = id
		}
	}
	return best
}
