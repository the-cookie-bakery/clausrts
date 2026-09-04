package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "infantry",
		Name:  "Infantry",
		Asset: "infantry.png",

		HP: 80, Attack: 12, Range: 1, Speed: 3, Sight: 5,
		AttackInterval: 1.0,

		Tags:     []CounterTag{TagInfantry, TagAntiArmor},
		Counters: []string{"airdefense", "cyber"},

		FrictionBase:  0.20,
		ReconPower:    0.15,
		GravityWeight: 1,
	})
}
