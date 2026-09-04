package entities

// BuildingKind is a static description of a building type.
type BuildingKind struct {
	ID    string
	Name  string
	Asset string

	HP    float64
	Sight float64

	// GravityWeight is the building's Schwerpunkt contribution. HQs are heavy.
	GravityWeight float64

	// Role flags.
	IsHQ     bool
	IsSupply bool // supply depot: shortens Verbindungslinien
	IsRadar  bool // reduces fog/friction in command radius
	IsBunker bool // shields against rockets/artillery
	IsFirewall bool // cyber defense
}

// BuildingKinds is the registry of all building kinds, keyed by ID.
var BuildingKinds = map[string]BuildingKind{}

// RegisterBuildingKind adds a BuildingKind to the registry.
func RegisterBuildingKind(k BuildingKind) { BuildingKinds[k.ID] = k }
