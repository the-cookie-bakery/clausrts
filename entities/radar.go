package entities

func init() {
	RegisterBuildingKind(BuildingKind{
		ID:    "radar",
		Name:  "Radar / Comms",
		Asset: "radar.png",

		HP: 200, Sight: 12,
		GravityWeight: 3,
		IsRadar:       true,
	})
}
