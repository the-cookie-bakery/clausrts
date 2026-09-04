package entities

func init() {
	RegisterBuildingKind(BuildingKind{
		ID:    "hq",
		Name:  "Headquarters",
		Asset: "hq.png",

		HP: 400, Sight: 8,
		GravityWeight: 10,
		IsHQ:          true,
	})
}
