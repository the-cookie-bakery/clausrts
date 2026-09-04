package entities

import (
	"clausrts/helpers"
)

// OrderKind enumerates the orders a formation can execute.
type OrderKind int

const (
	// OrderHold means "stay and engage targets in range".
	OrderHold OrderKind = iota
	// OrderMove means "path to Target tile".
	OrderMove
	// OrderAttack means "path to TargetUnit and destroy it".
	OrderAttack
)

// Order is an instruction assigned to a unit. Movement orders carry a
// precomputed tile Path so the movement system stays deterministic.
type Order struct {
	Kind    OrderKind
	Target  helpers.Point
	TargetUnit int
	Path    []helpers.Point
	PathIdx int
}

// Clear resets the order to Hold with no path.
func (o *Order) Clear() {
	o.Kind = OrderHold
	o.Target = helpers.Point{}
	o.TargetUnit = 0
	o.Path = nil
	o.PathIdx = 0
}
