// Package entities defines the game's data model: units, buildings, formations,
// factions and the Clausewitz "Dynamics" properties that the game teaches.
package entities

// Dynamics carries the live Clausewitz properties of a formation. These are the
// values the player can observe and influence, and are updated by the systems
// package each tick.
type Dynamics struct {
	// Friction (Friktion) in [0,1]. Higher means "war on paper" diverges from
	// real war: slower, less reliable. Raised by fog, distance from command,
	// jamming and rough terrain; lowered by reconnaissance and good command.
	Friction float64

	// Focus (Ökonomie der Kräfte) in [0,1]. Higher means forces concentrated
	// in space and time; multiplies attack effectiveness.
	Focus float64

	// Gravity is the formation's Schwerpunkt weight: how much this force is a
	// center of gravity of the conflict. Larger forces and HQs exert more.
	Gravity float64

	// Fog (Nachrichten) in [0,1]. The uncertainty about the immediate
	// surroundings: 0 = clear, 1 = blind.
	Fog float64

	// Moral in [0,1] is cohesion / the moralische Hauptpotenzen. Under stress
	// (casualties, friction) moral collapses and the formation stops fighting.
	Moral float64

	// Supply in [0,1] reflects the Verbindungslinien back to a supply source.
	Supply float64

	// Culmination in [0,1] is progress toward the Kulminationspunkt: past 1.0
	// the offensive has exhausted itself and defense is stronger.
	Culmination float64
}

// Reset returns d zeroed out.
func (d *Dynamics) Reset() {
	*d = Dynamics{}
}
