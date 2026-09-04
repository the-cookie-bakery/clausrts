package entities

// CounterTag categorises units for the counter matrix. The tags let the combat
// system resolve "what counters what" (Wechselwirkung: every weapon has a
// counter).
type CounterTag int

const (
	TagArmor        CounterTag = iota // tanks, armoured vehicles
	TagAntiArmor                       // ATGM infantry, anti-tank
	TagInfantry                        // soft infantry
	TagArtillery                       // rockets / long-range indirect fire
	TagAir                             // drones, aircraft
	TagAntiAir                         // air defense
	TagCyber                           // cyber offensive
	TagCyberDefense                    // firewall
	TagEW                              // electronic warfare / GPS scrambler
	TagRecon                           // GPS / satellite / drones recon
	TagStructure                       // buildings
)

// UnitKind is a static description of a unit type. It is registered by the
// per-entity files (entities/tank.go etc.) via RegisterUnitKind.
type UnitKind struct {
	ID    string // unique key, matches the asset basename, e.g. "tank"
	Name  string // human-readable name
	Asset string // asset filename, e.g. "tank.png"

	HP     float64 // max hit points
	Attack float64 // damage per hit
	Range  float64 // attack range in tiles (0 = melee/adjacent)
	Speed  float64 // movement speed in tiles per second
	Sight  float64 // visibility radius in tiles

	// CanFly lets the unit pass over water (drones).
	CanFly bool

	// AttackInterval is the seconds between attacks (0 = use default).
	AttackInterval float64

	Tags     []CounterTag
	Counters []string // UnitKind IDs this unit is strong against

	// Clausewitz modifiers.
	FrictionBase float64 // inherent friction of this unit type
	ReconPower   float64 // reduces nearby fog and friction
	JamPower     float64 // raises nearby enemy fog and friction
	GravityWeight float64 // contribution to the Schwerpunkt (center of gravity)
}

// UnitKinds is the registry of all unit kinds, keyed by ID.
var UnitKinds = map[string]UnitKind{}

// RegisterUnitKind adds a UnitKind to the registry. Called from per-entity
// files' init() functions.
func RegisterUnitKind(k UnitKind) { UnitKinds[k.ID] = k }
