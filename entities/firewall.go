package entities

func init() {
	RegisterBuildingKind(BuildingKind{
		ID:    "firewall",
		Name:  "Firewall",
		Asset: "firewall.png",

		HP: 150, Sight: 4,
		GravityWeight: 1,
		IsFirewall:    true,
	})
}
