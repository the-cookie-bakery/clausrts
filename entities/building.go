package entities

import (
	"clausrts/configs"
	"clausrts/helpers"
)

// Building is a static structure occupying a single tile.
type Building struct {
	ID      int
	Kind    string
	Faction FactionID
	Pos     helpers.Point

	HP    float64
	MaxHP float64
}

// NewBuilding creates a Building of the given kind. Returns nil if unknown.
func NewBuilding(id int, kindID string, faction FactionID, pos helpers.Point) *Building {
	k, ok := BuildingKinds[kindID]
	if !ok {
		return nil
	}
	return &Building{
		ID:      id,
		Kind:    kindID,
		Faction: faction,
		Pos:     pos,
		HP:      k.HP,
		MaxHP:   k.HP,
	}
}

// KindDef returns the BuildingKind definition.
func (b *Building) KindDef() BuildingKind { return BuildingKinds[b.Kind] }

// IsAlive reports whether the building has positive HP.
func (b *Building) IsAlive() bool { return b.HP > 0 }

// TileCenter returns the pixel center of the building's tile.
func (b *Building) TileCenter() (float64, float64) {
	return float64(b.Pos.X*configs.TileSize) + configs.TileSize/2,
		float64(b.Pos.Y*configs.TileSize) + configs.TileSize/2
}
