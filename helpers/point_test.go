package helpers

import "testing"

func TestDirectionTo(t *testing.T) {
	cases := []struct {
		from, to, want Point
	}{
		{Point{0, 0}, Point{0, -5}, Point{0, -1}},
		{Point{0, 0}, Point{3, -3}, Point{1, -1}},
		{Point{0, 0}, Point{-2, 0}, Point{-1, 0}},
		{Point{4, 4}, Point{4, 4}, Point{0, 0}},
	}
	for _, c := range cases {
		if got := c.from.DirectionTo(c.to); got != c.want {
			t.Errorf("DirectionTo(%v,%v) = %v, want %v", c.from, c.to, got, c.want)
		}
	}
}

func TestIsAdjacent(t *testing.T) {
	if !(Point{0, 0}).IsAdjacent(Point{1, 1}) {
		t.Error("diagonal neighbour should be adjacent")
	}
	if (Point{0, 0}).IsAdjacent(Point{2, 0}) {
		t.Error("two tiles away should not be adjacent")
	}
}

func TestClampF(t *testing.T) {
	if ClampF(1.5, 0, 1) != 1 {
		t.Error("ClampF should clamp to upper bound")
	}
	if ClampF(-1, 0, 1) != 0 {
		t.Error("ClampF should clamp to lower bound")
	}
	if ClampF(0.5, 0, 1) != 0.5 {
		t.Error("ClampF should preserve in-range values")
	}
}
