// Package worlds holds the mutable simulation state: the grid map, factions,
// units, buildings and the fixed-timestep tick loop.
package worlds

import (
	"clausrts/helpers"
)

// Terrain is the kind of ground a tile has.
type Terrain int

const (
	TerrainGrass Terrain = iota
	TerrainRoad
	TerrainForest
	TerrainUrban
	TerrainWater
	TerrainMountain
)

// Tile is a single map cell.
type Tile struct {
	Terrain Terrain
}

// Grid is the tile map.
type Grid struct {
	Width, Height int
	Tiles         [][]Tile
}

// NewGrid allocates an all-grass grid of the given dimensions.
func NewGrid(w, h int) *Grid {
	g := &Grid{Width: w, Height: h, Tiles: make([][]Tile, h)}
	for y := range g.Tiles {
		g.Tiles[y] = make([]Tile, w)
		for x := range g.Tiles[y] {
			g.Tiles[y][x].Terrain = TerrainGrass
		}
	}
	return g
}

// InBounds reports whether p is inside the grid.
func (g *Grid) InBounds(p helpers.Point) bool {
	return p.X >= 0 && p.X < g.Width && p.Y >= 0 && p.Y < g.Height
}

// Passable reports whether a non-flying unit may stand on tile p.
func (g *Grid) Passable(p helpers.Point) bool { return g.PassableFor(p, false) }

// PassableFor reports whether a unit may enter tile p. Flying units (drones)
// can cross water; mountains remain impassable to all units.
func (g *Grid) PassableFor(p helpers.Point, flying bool) bool {
	if !g.InBounds(p) {
		return false
	}
	switch g.Tiles[p.Y][p.X].Terrain {
	case TerrainMountain:
		return false
	case TerrainWater:
		return flying
	default:
		return true
	}
}

// TerrainCost returns the movement cost of tile p (lower is faster).
func (g *Grid) TerrainCost(p helpers.Point) float64 {
	if !g.InBounds(p) {
		return 1
	}
	switch g.Tiles[p.Y][p.X].Terrain {
	case TerrainRoad:
		return 0.8
	case TerrainForest:
		return 2.0
	case TerrainUrban:
		return 1.5
	default:
		return 1.0
	}
}

// FrictionMod returns the friction contribution of the terrain at tile p.
// Rough ground raises friction; roads lower it.
func (g *Grid) FrictionMod(p helpers.Point) float64 {
	if !g.InBounds(p) {
		return 0
	}
	switch g.Tiles[p.Y][p.X].Terrain {
	case TerrainRoad:
		return -0.15
	case TerrainForest:
		return 0.25
	case TerrainUrban:
		return 0.35
	default:
		return 0
	}
}

// Cover returns the damage-reduction factor (0..1) provided by the terrain.
func (g *Grid) Cover(p helpers.Point) float64 {
	if !g.InBounds(p) {
		return 0
	}
	switch g.Tiles[p.Y][p.X].Terrain {
	case TerrainForest:
		return 0.3
	case TerrainUrban:
		return 0.5
	default:
		return 0
	}
}
