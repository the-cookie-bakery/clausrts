package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "airdefense",
		Name:  "Air Defense",
		Asset: "airdefense.png",

		HP: 90, Attack: 35, Range: 4, Speed: 3, Sight: 7,
		AttackInterval: 1.0,

		Tags:     []CounterTag{TagAntiAir},
		Counters: []string{"drone", "artillery"},

		FrictionBase:  0.20,
		GravityWeight: 2,
	})
}
