// Package systems contains the simulation systems. Each system is a function
// operating on *worlds.World; Tick runs them in a fixed order every tick.
package systems

import (
	"clausrts/configs"
	"clausrts/worlds"
)

// Tick advances the simulation by one fixed timestep. It is a no-op while
// paused (Spannung und Ruhe: plan while paused, resolve in real time).
func Tick(w *worlds.World) {
	if w.Paused {
		return
	}
	w.Tick++
	dt := configs.SecondsPerTick

	UpdateFog(w)
	UpdateDynamics(w)
	UpdateMovement(w, dt)
	UpdateCombat(w, dt)
	UpdateTrinity(w, dt)
}
