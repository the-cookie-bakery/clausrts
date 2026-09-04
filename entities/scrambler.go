package entities

func init() {
	RegisterUnitKind(UnitKind{
		ID:    "scrambler",
		Name:  "GPS Scrambler",
		Asset: "scrambler.png",

		HP: 60, Attack: 5, Range: 1, Speed: 3, Sight: 5,
		AttackInterval: 1.0,

		Tags:     []CounterTag{TagEW},
		Counters: []string{"drone", "gps"},

		// The scrambler is regional electronic warfare: it raises enemy fog and
		// friction in a radius, but is itself a soft, high-value target.
		FrictionBase:  0.15,
		JamPower:      0.7,
		GravityWeight: 1,
	})
}
