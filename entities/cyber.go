package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "cyber",
		Name:  "Cyber Unit",
		Asset: "cyber.png",

		HP: 50, Attack: 25, Range: 1.5, Speed: 3, Sight: 5,
		AttackInterval: 1.5,

		Tags:     []CounterTag{TagCyber},
		Counters: []string{"gps", "scrambler", "radar", "depot"},

		FrictionBase:  0.15,
		GravityWeight: 1,
	})
}
