package entities

func init() {
	RegisterBuildingKind(BuildingKind{
		ID:    "bunker",
		Name:  "Bunker",
		Asset: "bunker.png",

		HP: 350, Sight: 4,
		GravityWeight: 2,
		IsBunker:      true,
	})
}
