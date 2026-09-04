package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "tank",
		Name:  "Tank",
		Asset: "tank.png",

		HP: 220, Attack: 45, Range: 2.5, Speed: 4, Sight: 6,
		AttackInterval: 1.2,

		Tags:     []CounterTag{TagArmor},
		Counters: []string{"infantry", "artillery", "airdefense", "cyber"},

		// Armour is powerful but heavy: high inherent friction, and its long
		// advances culminate quickly (see systems/culmination.go).
		FrictionBase:  0.35,
		GravityWeight: 3,
	})
}
