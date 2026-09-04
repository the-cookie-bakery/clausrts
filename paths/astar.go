// Package paths implements A* pathfinding over the tile grid. The cost function
// deliberately includes Clausewitzian terms (friction, gravity/threat) so that
// units observably route around jammed and threatened areas.
package paths

import (
	"container/heap"
	"math"

	"clausrts/helpers"
	"clausrts/worlds"
)

// CostExtras is an optional function returning the additional cost of a tile on
// top of the terrain cost. It lets callers fold in friction and enemy gravity.
// Returning +Inf marks a tile as unenterable.
type CostExtras func(p helpers.Point) float64

type node struct {
	pos   helpers.Point
	g     float64
	f     float64
	index int
}

type priorityQueue []*node

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].f < pq[j].f }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *priorityQueue) Push(x any)         { n := x.(*node); n.index = len(*pq); *pq = append(*pq, n) }
func (pq *priorityQueue) Pop() any           { old := *pq; n := old[len(old)-1]; old[len(old)-1] = nil; *pq = old[:len(old)-1]; return n }

// FindPath returns the sequence of tiles from `from` to `to` (inclusive of `to`,
// exclusive of `from`) or nil if no path exists. Movement is 8-directional with
// corner-cut prevention. The extras function adds Clausewitzian cost on top of
// terrain cost; pass nil for terrain-only cost. The flying flag lets the
// path cross water (drones).
func FindPath(g *worlds.Grid, from, to helpers.Point, flying bool, extras CostExtras) []helpers.Point {
	if !g.InBounds(to) || !g.PassableFor(to, flying) {
		return nil
	}
	if from.Equal(to) {
		return []helpers.Point{to}
	}

	start := &node{pos: from}
	open := &priorityQueue{start}
	heap.Init(open)

	cameFrom := map[helpers.Point]helpers.Point{}
	gScore := map[helpers.Point]float64{from: 0}
	fScore := map[helpers.Point]float64{from: heuristic(from, to)}

	for open.Len() > 0 {
		current := heap.Pop(open).(*node)
		if current.pos.Equal(to) {
			return reconstruct(cameFrom, to)
		}

		for _, dir := range helpers.Compass {
			next := current.pos.Add(dir)
			if !g.PassableFor(next, flying) {
				continue
			}
			// Corner-cut prevention: diagonal moves require both adjacent
			// orthogonal tiles to be passable.
			if dir.X != 0 && dir.Y != 0 {
				if !g.PassableFor(helpers.Point{X: current.pos.X + dir.X, Y: current.pos.Y}, flying) ||
					!g.PassableFor(helpers.Point{X: current.pos.X, Y: current.pos.Y + dir.Y}, flying) {
					continue
				}
			}

			var extra float64
			if extras != nil {
				extra = extras(next)
				if math.IsInf(extra, 1) {
					continue
				}
			}

			step := 1.0
			if dir.X != 0 && dir.Y != 0 {
				step = math.Sqrt2
			}

			tentative := gScore[current.pos] + g.TerrainCost(next)*step + extra
			if known, ok := gScore[next]; ok && tentative >= known {
				continue
			}

			cameFrom[next] = current.pos
			gScore[next] = tentative
			fScore[next] = tentative + heuristic(next, to)

			n := &node{pos: next, g: tentative, f: fScore[next]}
			heap.Push(open, n)
		}
	}

	return nil
}

func heuristic(a, b helpers.Point) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx+dy*dy)
}

func reconstruct(cameFrom map[helpers.Point]helpers.Point, to helpers.Point) []helpers.Point {
	path := []helpers.Point{to}
	cur := to
	for {
		prev, ok := cameFrom[cur]
		if !ok {
			break
		}
		path = append(path, prev)
		cur = prev
	}
	// Reverse to get from -> to order, excluding the start tile.
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path[1:]
}
