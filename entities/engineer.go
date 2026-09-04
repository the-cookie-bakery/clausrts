package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "engineer",
		Name:  "Engineer",
		Asset: "engineer.png",

		HP: 70, Attack: 8, Range: 1, Speed: 3, Sight: 5,
		AttackInterval: 1.0,

		Tags: []CounterTag{TagInfantry},

		// Engineers build and repair structures (bunker, firewall); they are
		// soft targets for artillery.
		FrictionBase:  0.20,
		GravityWeight: 1,
	})
}
