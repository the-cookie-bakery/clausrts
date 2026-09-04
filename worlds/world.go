package worlds

import (
	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/utils"
)

// Vis states for a fog tile.
const (
	VisHidden   byte = 0 // never seen
	VisVisible  byte = 1 // currently seen
	VisExplored byte = 2 // seen before, currently out of sight
)

// World owns all mutable simulation state.
type World struct {
	Grid *Grid

	Units     map[int]*entities.Unit
	Buildings map[int]*entities.Building
	Formations map[int]*entities.Formation

	Factions map[entities.FactionID]*entities.Faction

	// Fog holds a per-faction visibility grid.
	Fog map[entities.FactionID]*Fog

	Tick   uint64
	Paused bool

	RNG *utils.RNG

	nextID int
}

// Fog is a per-faction visibility grid over the map.
type Fog struct {
	W, H   int
	Values [][]byte
}

// NewFog allocates an all-hidden fog grid.
func NewFog(w, h int) *Fog {
	f := &Fog{W: w, H: h, Values: make([][]byte, h)}
	for y := range f.Values {
		f.Values[y] = make([]byte, w)
	}
	return f
}

// At returns the visibility of tile p.
func (f *Fog) At(p helpers.Point) byte {
	if p.X < 0 || p.X >= f.W || p.Y < 0 || p.Y >= f.H {
		return VisHidden
	}
	return f.Values[p.Y][p.X]
}

// Set sets the visibility of tile p.
func (f *Fog) Set(p helpers.Point, v byte) {
	if p.X < 0 || p.X >= f.W || p.Y < 0 || p.Y >= f.H {
		return
	}
	f.Values[p.Y][p.X] = v
}

// NewWorld creates an empty world over a grid of the given size.
func NewWorld(w, h int, seed int64) *World {
	world := &World{
		Grid:       NewGrid(w, h),
		Units:      map[int]*entities.Unit{},
		Buildings:  map[int]*entities.Building{},
		Formations: map[int]*entities.Formation{},
		Factions:   map[entities.FactionID]*entities.Faction{},
		Fog:        map[entities.FactionID]*Fog{},
		RNG:        utils.NewRNG(seed),
	}
	return world
}

// AllocID returns the next unique entity ID.
func (w *World) AllocID() int {
	w.nextID++
	return w.nextID
}

// AddFaction registers a faction and allocates its fog grid.
func (w *World) AddFaction(f *entities.Faction) {
	w.Factions[f.ID] = f
	w.Fog[f.ID] = NewFog(w.Grid.Width, w.Grid.Height)
}

// AddUnit inserts a unit into the world and returns it.
func (w *World) AddUnit(u *entities.Unit) *entities.Unit {
	w.Units[u.ID] = u
	return u
}

// SpawnUnit creates a unit of the given kind at tile pos for faction.
func (w *World) SpawnUnit(kindID string, faction entities.FactionID, pos helpers.Point) *entities.Unit {
	u := entities.NewUnit(w.AllocID(), kindID, faction, pos)
	if u == nil {
		return nil
	}
	return w.AddUnit(u)
}

// AddBuilding inserts a building into the world and returns it.
func (w *World) AddBuilding(b *entities.Building) *entities.Building {
	w.Buildings[b.ID] = b
	return b
}

// SpawnBuilding creates a building of the given kind at tile pos for faction.
func (w *World) SpawnBuilding(kindID string, faction entities.FactionID, pos helpers.Point) *entities.Building {
	b := entities.NewBuilding(w.AllocID(), kindID, faction, pos)
	if b == nil {
		return nil
	}
	return w.AddBuilding(b)
}

// AddFormation registers a formation.
func (w *World) AddFormation(f *entities.Formation) *entities.Formation {
	w.Formations[f.ID] = f
	return f
}

// UnitAt returns the live unit occupying tile p, or nil.
func (w *World) UnitAt(p helpers.Point) *entities.Unit {
	for _, u := range w.Units {
		if u.IsAlive() && u.Pos == p {
			return u
		}
	}
	return nil
}

// BuildingAt returns the building occupying tile p, or nil.
func (w *World) BuildingAt(p helpers.Point) *entities.Building {
	for _, b := range w.Buildings {
		if b.IsAlive() && b.Pos == p {
			return b
		}
	}
	return nil
}

// Occupied reports whether tile p is blocked by a unit or building.
func (w *World) Occupied(p helpers.Point) bool {
	return w.UnitAt(p) != nil || w.BuildingAt(p) != nil
}

// RemoveUnit deletes a unit from the world.
func (w *World) RemoveUnit(id int) { delete(w.Units, id) }

// RemoveBuilding deletes a building from the world.
func (w *World) RemoveBuilding(id int) { delete(w.Buildings, id) }

// SquareDistance returns the squared Euclidean distance in tile space.
func SquareDistance(a, b helpers.Point) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return dx*dx + dy*dy
}

// Faction returns the faction with the given ID.
func (w *World) Faction(id entities.FactionID) *entities.Faction { return w.Factions[id] }

// UnitsOf returns the live units belonging to a faction.
func (w *World) UnitsOf(faction entities.FactionID) []*entities.Unit {
	var out []*entities.Unit
	for _, u := range w.Units {
		if u.Faction == faction && u.IsAlive() {
			out = append(out, u)
		}
	}
	return out
}

// BuildingsOf returns the live buildings belonging to a faction.
func (w *World) BuildingsOf(faction entities.FactionID) []*entities.Building {
	var out []*entities.Building
	for _, b := range w.Buildings {
		if b.Faction == faction && b.IsAlive() {
			out = append(out, b)
		}
	}
	return out
}
