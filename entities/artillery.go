package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "artillery",
		Name:  "Rocket Artillery",
		Asset: "artillery.png",

		HP: 90, Attack: 60, Range: 8, Speed: 2, Sight: 6,
		AttackInterval: 2.5,

		Tags:     []CounterTag{TagArtillery},
		Counters: []string{"infantry", "tank", "engineer"},

		FrictionBase:  0.25,
		GravityWeight: 2,
	})
}
