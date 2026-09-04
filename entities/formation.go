package entities

// Formation groups sub-units into a battalion. Formations carry the aggregate
// Clausewitz Dynamics that are computed from their composition, and are the
// primary object of player control (select a formation, order it, merge/split).
type Formation struct {
	ID     int
	Faction FactionID

	// Units holds the IDs of the member sub-units.
	Units []int

	// Dynamics are the formation-level Clausewitz properties, aggregated from
	// member units each tick.
	Dynamics Dynamics
}

// NewFormation returns an empty formation.
func NewFormation(id int, faction FactionID) *Formation {
	return &Formation{ID: id, Faction: faction}
}

// Add appends a unit ID to the formation.
func (f *Formation) Add(unitID int) { f.Units = append(f.Units, unitID) }

// Remove removes the first occurrence of unitID, returning true if removed.
func (f *Formation) Remove(unitID int) bool {
	for i, id := range f.Units {
		if id == unitID {
			f.Units = append(f.Units[:i], f.Units[i+1:]...)
			return true
		}
	}
	return false
}

// Contains reports whether unitID is a member.
func (f *Formation) Contains(unitID int) bool {
	for _, id := range f.Units {
		if id == unitID {
			return true
		}
	}
	return false
}

// Size returns the number of sub-units in the formation.
func (f *Formation) Size() int { return len(f.Units) }
