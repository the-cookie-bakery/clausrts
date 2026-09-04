package systems

import (
	"testing"

	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/worlds"
)

func TestBuildingTakesDamage(t *testing.T) {
	w := worlds.NewWorld(10, 10, 1)
	w.AddFaction(entities.NewFaction(0, "a", [3]uint8{1, 1, 1}))
	w.AddFaction(entities.NewFaction(1, "b", [3]uint8{2, 2, 2}))

	b := w.SpawnBuilding("hq", 1, helpers.Point{X: 5, Y: 5})
	u := w.SpawnUnit("tank", 0, helpers.Point{X: 5, Y: 4})
	if b == nil || u == nil {
		t.Fatal("spawn failed")
	}
	u.Order.Kind = entities.OrderAttack
	u.Order.TargetUnit = b.ID

	startHP := b.HP
	for i := 0; i < 300; i++ {
		Tick(w)
	}
	if b.HP >= startHP {
		t.Fatalf("building took no damage: %v -> %v (faction govt %.2f)", startHP, b.HP, w.Faction(1).Government)
	}
	t.Logf("building HP: %v -> %v after 300 ticks", startHP, b.HP)
}
