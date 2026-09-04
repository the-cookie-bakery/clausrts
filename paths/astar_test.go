package paths

import (
	"testing"

	"clausrts/helpers"
	"clausrts/worlds"
)

func TestFindPathStraightLine(t *testing.T) {
	g := worlds.NewGrid(5, 5)
	path := FindPath(g, helpers.Point{X: 0, Y: 0}, helpers.Point{X: 4, Y: 4}, false, nil)
	if path == nil {
		t.Fatal("expected a path")
	}
	if path[len(path)-1] != (helpers.Point{X: 4, Y: 4}) {
		t.Fatalf("path must end at target, got %v", path[len(path)-1])
	}
	if len(path) < 4 {
		t.Fatalf("expected at least 4 steps, got %d", len(path))
	}
}

func TestFindPathObstacle(t *testing.T) {
	g := worlds.NewGrid(6, 3)
	// Wall of mountains across x=2, with a gap at y=0.
	g.Tiles[1][2].Terrain = worlds.TerrainMountain
	g.Tiles[2][2].Terrain = worlds.TerrainMountain
	path := FindPath(g, helpers.Point{X: 0, Y: 1}, helpers.Point{X: 5, Y: 1}, false, nil)
	if path == nil {
		t.Fatal("expected a path around the wall")
	}
	for _, p := range path {
		if !g.Passable(p) {
			t.Fatalf("path crosses impassable tile %v", p)
		}
	}
}

func TestFindPathNoRoute(t *testing.T) {
	g := worlds.NewGrid(3, 1)
	g.Tiles[0][1].Terrain = worlds.TerrainWater
	if path := FindPath(g, helpers.Point{X: 0, Y: 0}, helpers.Point{X: 2, Y: 0}, false, nil); path != nil {
		t.Fatalf("expected no path across water, got %v", path)
	}
}

func TestFindPathFlyingOverWater(t *testing.T) {
	g := worlds.NewGrid(3, 1)
	g.Tiles[0][1].Terrain = worlds.TerrainWater
	if path := FindPath(g, helpers.Point{X: 0, Y: 0}, helpers.Point{X: 2, Y: 0}, false, nil); path != nil {
		t.Fatalf("non-flying units must not cross water")
	}
	if path := FindPath(g, helpers.Point{X: 0, Y: 0}, helpers.Point{X: 2, Y: 0}, true, nil); path == nil {
		t.Fatalf("flying units should cross water")
	}
}
