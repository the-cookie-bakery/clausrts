package entities

func init() {
	RegisterBuildingKind(BuildingKind{
		ID:    "depot",
		Name:  "Supply Depot",
		Asset: "depot.png",

		HP: 250, Sight: 4,
		GravityWeight: 4,
		IsSupply:      true,
	})
}
