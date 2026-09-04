package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "drone",
		Name:  "Drone",
		Asset: "drone.png",

		HP: 40, Attack: 18, Range: 2, Speed: 6, Sight: 8,
		AttackInterval: 0.8,
		CanFly:         true,

		Tags:     []CounterTag{TagAir, TagRecon},
		Counters: []string{"tank", "artillery", "infantry"},

		// Drones are the eyes of the battlefield: they dissolve fog and lower
		// friction around themselves, but are fragile and jam-prone.
		FrictionBase:  0.10,
		ReconPower:    0.6,
		GravityWeight: 1,
	})
}
