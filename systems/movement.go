package systems

import (
	"math"

	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

// UpdateMovement advances every unit along its order. Movement is grid-locked:
// a unit interpolates tile-to-tile and only changes logical position when a
// step completes, which keeps pathfinding and combat deterministic.
func UpdateMovement(w *worlds.World, dt float64) {
	for _, u := range w.Units {
		if !u.IsAlive() {
			continue
		}
		switch u.Order.Kind {
		case entities.OrderMove:
			stepMove(w, u, dt)
		case entities.OrderAttack:
			stepAttack(w, u, dt)
		}
	}
}

// stepMove advances a unit along its precomputed path.
func stepMove(w *worlds.World, u *entities.Unit, dt float64) {
	if u.Moving {
		advanceStep(w, u, dt)
		return
	}
	if u.Order.PathIdx >= len(u.Order.Path) {
		u.Order.Clear()
		return
	}
	next := u.Order.Path[u.Order.PathIdx]
	if !w.Grid.PassableFor(next, u.KindDef().CanFly) || blocked(w, next, u) {
		// Blocked: hold in place until the way clears (friction of real war).
		return
	}
	startStep(u, next)
	u.Order.PathIdx++
}

// stepAttack closes to the ordered target and attacks when in range.
func stepAttack(w *worlds.World, u *entities.Unit, dt float64) {
	tgt := findTargetByID(w, u.Order.TargetUnit)
	if tgt == nil {
		u.Order.Clear()
		return
	}
	if u.Moving {
		advanceStep(w, u, dt)
		return
	}
	dist := math.Sqrt(worlds.SquareDistance(u.Pos, tgt.pos))
	if dist <= u.KindDef().Range {
		u.Facing = u.Pos.DirectionTo(tgt.pos)
		return // combat system handles the actual attack
	}
	if !greedyStep(w, u, tgt.pos) {
		return // fully blocked
	}
}

// advanceStep moves the unit along the current tile step; on completion the
// logical tile position is committed.
func advanceStep(w *worlds.World, u *entities.Unit, dt float64) {
	speed := u.KindDef().Speed * (1 - 0.5*u.Dynamics.Friction)
	if speed <= 0 {
		speed = u.KindDef().Speed * 0.2
	}
	u.MoveProgress += speed * dt
	if u.MoveProgress >= 1 {
		// Step complete: commit the logical position and snap the render
		// position onto the new tile center. Rendering at progress 0 here
		// would snap the sprite back onto the previous tile.
		u.Pos = u.MoveTo
		u.MoveProgress = 0
		u.Moving = false
		u.RenderX, u.RenderY = entities.TileCenter(u.Pos)
		return
	}
	u.RenderX, u.RenderY = interpolate(u.MoveFrom, u.MoveTo, u.MoveProgress)
}

func interpolate(from, to helpers.Point, t float64) (float64, float64) {
	fx, fy := entities.TileCenter(from)
	tx, ty := entities.TileCenter(to)
	return fx + (tx-fx)*t, fy + (ty-fy)*t
}

func startStep(u *entities.Unit, dest helpers.Point) {
	u.MoveFrom = u.Pos
	u.MoveTo = dest
	u.MoveProgress = 0
	u.Moving = true
	u.Facing = u.Pos.DirectionTo(dest)
}

// greedyStep starts a single step toward target, choosing the passable,
// unblocked neighbour closest to it. Returns false if no neighbour is available.
func greedyStep(w *worlds.World, u *entities.Unit, target helpers.Point) bool {
	best := u.Pos
	bestDist := worlds.SquareDistance(u.Pos, target)
	for _, dir := range helpers.Compass {
		n := u.Pos.Add(dir)
		if !w.Grid.PassableFor(n, u.KindDef().CanFly) || blocked(w, n, u) {
			continue
		}
		if d := worlds.SquareDistance(n, target); d < bestDist {
			bestDist = d
			best = n
		}
	}
	if best.Equal(u.Pos) {
		return false
	}
	startStep(u, best)
	return true
}

// blocked reports whether tile p is occupied by any unit or building.
func blocked(w *worlds.World, p helpers.Point, self *entities.Unit) bool {
	if w.BuildingAt(p) != nil {
		return true
	}
	if o := w.UnitAt(p); o != nil && o.ID != self.ID {
		return true
	}
	return false
}
