// Package configs holds global, immutable game configuration.
package configs

const (
	// TileSize is the side length in pixels of a single map tile. It is the
	// fundamental unit of the grid, pathfinding and rendering.
	TileSize = 32

	// WindowWidth and WindowHeight are the initial window dimensions.
	WindowWidth  = 1280
	WindowHeight = 720

	// TicksPerSecond is the fixed simulation rate. The game runs deterministic
	// logic ticks at this rate regardless of the display refresh rate.
	TicksPerSecond = 60

	// SecondsPerTick is the duration of a single simulation tick.
	SecondsPerTick = 1.0 / TicksPerSecond

	// Movement interpolates tile-to-tile in MoveDuration ticks (smoothed over
	// multiple ticks for fluid rendering while staying grid-locked).
	MoveTicks = 12
)

// Camera limits.
const (
	MinZoom = 0.5
	MaxZoom = 4.0
	ZoomStep = 1.1
)

// Faction identifiers used as array indices.
const (
	FactionPlayer = 0
	FactionEnemy  = 1
)
