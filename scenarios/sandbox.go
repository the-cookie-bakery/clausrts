// Package scenarios builds playable worlds. The sandbox scenario places two
// factions on a procedurally-decorated map for free play; the campaign
// scenarios (milestones M4/M5) add objectives, triggers and lessons.
package scenarios

import (
	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

const (
	mapW = 48
	mapH = 36
)

// Sandbox builds a world with two factions ready to fight.
func Sandbox(seed int64) *worlds.World {
	w := worlds.NewWorld(mapW, mapH, seed)

	player := entities.NewFaction(0, "Blue", [3]uint8{70, 140, 230})
	enemy := entities.NewFaction(1, "Red", [3]uint8{230, 70, 70})
	w.AddFaction(player)
	w.AddFaction(enemy)

	decorateTerrain(w)

	// Player base (top-left).
	playerUnits := map[string][]helpers.Point{
		"tank":       {{X: 5, Y: 8}, {X: 6, Y: 9}},
		"infantry":   {{X: 7, Y: 7}, {X: 8, Y: 8}},
		"artillery":  {{X: 9, Y: 6}},
		"drone":      {{X: 4, Y: 10}},
		"gps":        {{X: 8, Y: 6}},
		"airdefense": {{X: 6, Y: 7}},
		"engineer":   {{X: 9, Y: 9}},
	}
	placeUnits(w, 0, playerUnits)

	w.SpawnBuilding("hq", 0, helpers.Point{X: 4, Y: 4})
	w.SpawnBuilding("depot", 0, helpers.Point{X: 8, Y: 4})
	w.SpawnBuilding("radar", 0, helpers.Point{X: 4, Y: 8})
	w.SpawnBuilding("bunker", 0, helpers.Point{X: 6, Y: 5})
	w.SpawnBuilding("firewall", 0, helpers.Point{X: 5, Y: 6})

	// Enemy base (bottom-right).
	enemyUnits := map[string][]helpers.Point{
		"tank":       {{X: mapW - 6, Y: mapH - 9}, {X: mapW - 7, Y: mapH - 10}},
		"infantry":   {{X: mapW - 8, Y: mapH - 8}, {X: mapW - 9, Y: mapH - 9}},
		"artillery":  {{X: mapW - 10, Y: mapH - 7}},
		"drone":      {{X: mapW - 5, Y: mapH - 11}},
		"scrambler":  {{X: mapW - 9, Y: mapH - 7}},
		"airdefense": {{X: mapW - 7, Y: mapH - 8}},
		"cyber":      {{X: mapW - 10, Y: mapH - 10}},
	}
	placeUnits(w, 1, enemyUnits)

	w.SpawnBuilding("hq", 1, helpers.Point{X: mapW - 5, Y: mapH - 5})
	w.SpawnBuilding("depot", 1, helpers.Point{X: mapW - 9, Y: mapH - 5})
	w.SpawnBuilding("radar", 1, helpers.Point{X: mapW - 5, Y: mapH - 9})
	w.SpawnBuilding("bunker", 1, helpers.Point{X: mapW - 7, Y: mapH - 6})
	w.SpawnBuilding("firewall", 1, helpers.Point{X: mapW - 6, Y: mapH - 7})

	buildFormations(w)

	return w
}

func placeUnits(w *worlds.World, faction entities.FactionID, units map[string][]helpers.Point) {
	for kind, positions := range units {
		for _, p := range positions {
			w.SpawnUnit(kind, faction, p)
		}
	}
}

// buildFormations groups each faction's units into a single battalion for now;
// per-unit transfer between formations is a later milestone.
func buildFormations(w *worlds.World) {
	for _, f := range w.Factions {
		form := entities.NewFormation(w.AllocID(), f.ID)
		for _, u := range w.UnitsOf(f.ID) {
			form.Add(u.ID)
			u.Formation = form.ID
		}
		w.AddFormation(form)
	}
}

func decorateTerrain(w *worlds.World) {
	g := w.Grid

	// A river across the middle with two bridges.
	riverY := mapH / 2
	for x := 0; x < mapW; x++ {
		if x == mapW/3 || x == 2*mapW/3 {
			continue // bridges
		}
		g.Tiles[riverY][x].Terrain = worlds.TerrainWater
	}

	// Forest clusters.
	for i := 0; i < 60; i++ {
		x := w.RNG.Intn(mapW)
		y := w.RNG.Intn(mapH)
		if g.Tiles[y][x].Terrain == worlds.TerrainWater {
			continue
		}
		g.Tiles[y][x].Terrain = worlds.TerrainForest
	}

	// A central urban block.
	for y := riverY - 1; y <= riverY+1; y++ {
		for x := mapW/2 - 3; x <= mapW/2+3; x++ {
			if g.Tiles[y][x].Terrain == worlds.TerrainWater {
				continue
			}
			g.Tiles[y][x].Terrain = worlds.TerrainUrban
		}
	}

	// Mountains along the top and bottom edges.
	for x := 0; x < mapW; x++ {
		g.Tiles[0][x].Terrain = worlds.TerrainMountain
		g.Tiles[mapH-1][x].Terrain = worlds.TerrainMountain
	}

	// Roads connecting both bases vertically.
	for y := 0; y < mapH; y++ {
		if g.Tiles[y][8].Terrain == worlds.TerrainWater {
			continue
		}
		g.Tiles[y][8].Terrain = worlds.TerrainRoad
	}
}
