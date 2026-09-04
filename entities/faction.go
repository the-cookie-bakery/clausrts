package entities

import (
	"clausrts/helpers"
)

// FactionID identifies a side in the conflict.
type FactionID int

// Faction is a side in the conflict. Its three bars are the Clausewitz Trinity:
// People (passion), Military (chance), Government (reason).
type Faction struct {
	ID   FactionID
	Name string

	// Trinity bars, each in [0,1]. Letting any reach 0 ends the game.
	People     float64
	Military   float64
	Government float64

	// Cumulative friendly casualties, used to compute People decay.
	Losses int

	// Schwerpunkt is the faction's computed center of gravity (weighted mean of
	// unit/building positions by gravity). TotalGravity is the total weight.
	Schwerpunkt helpers.Point
	TotalGravity float64

	// ColorRGB is the faction's primary colour for rendering (0..255).
	ColorRGB [3]uint8
}

// NewFaction returns a Faction with full Trinity bars.
func NewFaction(id FactionID, name string, rgb [3]uint8) *Faction {
	return &Faction{
		ID:         id,
		Name:       name,
		People:     1,
		Military:   1,
		Government: 1,
		ColorRGB:   rgb,
	}
}
