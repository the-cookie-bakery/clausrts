// Package helpers contains small, dependency-free geometry and grid utilities.
package helpers

// Point is an integer position in tile space. (0,0) is the top-left tile,
// +X is right, +Y is down.
type Point struct {
	X, Y int
}

// Add returns p + q.
func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }

// Sub returns p - q.
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

// Equal reports whether p and q are the same tile.
func (p Point) Equal(q Point) bool { return p.X == q.X && p.Y == q.Y }

// Manhattan returns the Manhattan (taxicab) distance to q.
func (p Point) Manhattan(q Point) int {
	return abs(p.X-q.X) + abs(p.Y-q.Y)
}

// Chebyshev returns the Chebyshev (chessboard) distance to q.
func (p Point) Chebyshev(q Point) int {
	dx := abs(p.X - q.X)
	dy := abs(p.Y - q.Y)
	if dx > dy {
		return dx
	}
	return dy
}

// IsAdjacent reports whether q is one of the 8 neighbours of p.
func (p Point) IsAdjacent(q Point) bool {
	return p.Chebyshev(q) == 1
}

// DirectionTo returns the unit direction (one of the 8 compass directions)
// from p toward q, or (0,0) if p == q.
func (p Point) DirectionTo(q Point) Point {
	d := q.Sub(p)
	return Point{sign(d.X), sign(d.Y)}
}

// Compass is the set of 8 unit directions in tile space.
var Compass = [8]Point{
	{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1},
}

// Orthogonal is the set of 4 cardinal directions, used for corner-cut checks.
var Orthogonal = [4]Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sign(v int) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	default:
		return 0
	}
}

// ClampF clamps v into [lo, hi].
func ClampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// MaxF returns the larger of a and b.
func MaxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinF returns the smaller of a and b.
func MinF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
