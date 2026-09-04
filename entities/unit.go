package entities

import (
	"clausrts/configs"
	"clausrts/helpers"
)

// Unit is a single sub-unit token occupying exactly one tile. Units can be
// grouped into Formations (battalions) and transferred between them.
type Unit struct {
	ID      int
	Kind    string // UnitKind.ID, e.g. "tank"
	Faction FactionID
	Pos     helpers.Point
	Facing  helpers.Point

	HP    float64
	MaxHP float64

	// Dynamics are the live Clausewitz properties, recomputed each tick.
	Dynamics Dynamics

	// Order is the current instruction.
	Order Order

	// Movement interpolation state (grid-locked with smooth rendering).
	// RenderX/RenderY are the pixel position used for drawing; Pos is the
	// logical tile.
	RenderX, RenderY float64
	MoveFrom, MoveTo helpers.Point
	MoveProgress     float64 // 0..1 within the current tile step
	Moving           bool

	// AttackCooldown counts down from the kind's attack interval.
	AttackCooldown float64

	// Formation is the ID of the formation this unit belongs to (0 = none).
	Formation int
}

// NewUnit creates a Unit of the given kind owned by faction, at tile pos.
// It returns nil if the kind is unknown.
func NewUnit(id int, kindID string, faction FactionID, pos helpers.Point) *Unit {
	k, ok := UnitKinds[kindID]
	if !ok {
		return nil
	}
	u := &Unit{
		ID:      id,
		Kind:    kindID,
		Faction: faction,
		Pos:     pos,
		Facing:  helpers.Point{X: 0, Y: -1}, // face north by default
		HP:      k.HP,
		MaxHP:   k.HP,
	}
	u.RenderX, u.RenderY = TileCenter(pos)
	u.Dynamics.Moral = 1
	u.Dynamics.Supply = 1
	return u
}

// Kind returns the UnitKind definition for this unit.
func (u *Unit) KindDef() UnitKind { return UnitKinds[u.Kind] }

// IsAlive reports whether the unit has positive HP.
func (u *Unit) IsAlive() bool { return u.HP > 0 }

// TileCenter returns the pixel center of a tile (used by entities and systems).
func TileCenter(p helpers.Point) (float64, float64) {
	return float64(p.X*configs.TileSize) + configs.TileSize/2, float64(p.Y*configs.TileSize) + configs.TileSize/2
}
