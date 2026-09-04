// Package commanders defines the commander abstraction: any actor that turns
// high-level intent into unit orders. The player commander is driven by
// inputs; the AI commander by the ai package.
package commanders

import (
	"clausrts/worlds"
)

// Commander issues orders to the units it controls each step.
type Commander interface {
	// Step issues orders. Called once per simulation tick.
	Step(w *worlds.World)
}
