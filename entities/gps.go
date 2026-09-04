package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "gps",
		Name:  "GPS Support",
		Asset: "gps.png",

		HP: 50, Attack: 5, Range: 1, Speed: 3, Sight: 6,
		AttackInterval: 1.0,

		Tags: []CounterTag{TagRecon},

		// GPS/satellite support dissolves fog, lowers friction and raises focus
		// across a wide radius — but it is exactly what scramblers jam.
		FrictionBase:  0.05,
		ReconPower:    0.9,
		GravityWeight: 1,
	})
}
