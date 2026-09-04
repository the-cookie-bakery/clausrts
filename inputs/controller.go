// Package inputs translates raw ebiten input into selection and orders, and
// drives the camera. It is the "player commander" input surface.
package inputs

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"clausrts/entities"
	"clausrts/helpers"
	"clausrts/paths"
	"clausrts/renderers"
	"clausrts/worlds"
)

// Controller owns the camera and the current selection.
type Controller struct {
	Cam      *renderers.Camera
	Player   entities.FactionID
	Selected map[int]bool

	dragging  bool
	dragStart helpers.Point // screen coords
	dragCur   helpers.Point
	panning   bool
	panLast   helpers.Point
}

// NewController returns a Controller for the given player faction.
func NewController(cam *renderers.Camera, player entities.FactionID) *Controller {
	return &Controller{
		Cam:      cam,
		Player:   player,
		Selected: map[int]bool{},
	}
}

// Update processes input for one frame. It runs even while paused so the player
// can keep planning (Spannung und Ruhe).
func (c *Controller) Update(w *worlds.World) {
	c.updateCamera()
	c.Cam.ClampToWorld(w.Grid.Width, w.Grid.Height)
	c.updateSelection(w)
	c.updateOrders(w)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		w.Paused = !w.Paused
	}
}

func (c *Controller) updateCamera() {
	const panSpeed = 8.0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		c.Cam.Pan(0, -panSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		c.Cam.Pan(0, panSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		c.Cam.Pan(-panSpeed, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		c.Cam.Pan(panSpeed, 0)
	}

	mx, my := ebiten.CursorPosition()
	if _, wy := ebiten.Wheel(); wy != 0 {
		if wy > 0 {
			c.Cam.ZoomAt(float64(mx), float64(my), 1.1)
		} else {
			c.Cam.ZoomAt(float64(mx), float64(my), 1/1.1)
		}
	}

	// Middle-drag to pan.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		c.panning = true
		c.panLast = helpers.Point{X: mx, Y: my}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		c.panning = false
	}
	if c.panning {
		dx := float64(mx - c.panLast.X)
		dy := float64(my - c.panLast.Y)
		c.Cam.Pan(dx, dy)
		c.panLast = helpers.Point{X: mx, Y: my}
	}
}

func (c *Controller) updateSelection(w *worlds.World) {
	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.dragging = true
		c.dragStart = helpers.Point{X: mx, Y: my}
		c.dragCur = c.dragStart
	}
	if c.dragging {
		c.dragCur = helpers.Point{X: mx, Y: my}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.dragging = false
		// A click (no drag) selects the unit under the cursor; a drag selects
		// all player units within the box.
		if c.dragStart.Chebyshev(c.dragCur) <= 2 {
			c.clickSelect(w, mx, my)
		} else {
			c.boxSelect(w)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		c.Selected = map[int]bool{}
	}
}

func (c *Controller) clickSelect(w *worlds.World, mx, my int) {
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
		// Shift-click adds/removes a single unit.
		if u := unitAtScreen(w, c.Cam, c.Player, mx, my); u != nil {
			if c.Selected[u.ID] {
				delete(c.Selected, u.ID)
			} else {
				c.Selected[u.ID] = true
			}
		}
		return
	}
	c.Selected = map[int]bool{}
	if u := unitAtScreen(w, c.Cam, c.Player, mx, my); u != nil {
		c.Selected[u.ID] = true
	}
}

func (c *Controller) boxSelect(w *worlds.World) {
	x0 := min(c.dragStart.X, c.dragCur.X)
	x1 := max(c.dragStart.X, c.dragCur.X)
	y0 := min(c.dragStart.Y, c.dragCur.Y)
	y1 := max(c.dragStart.Y, c.dragCur.Y)

	c.Selected = map[int]bool{}
	for _, u := range w.UnitsOf(c.Player) {
		sx, sy := c.Cam.WorldToScreen(u.RenderX, u.RenderY)
		if sx >= float64(x0) && sx <= float64(x1) && sy >= float64(y0) && sy <= float64(y1) {
			c.Selected[u.ID] = true
		}
	}
}

func (c *Controller) updateOrders(w *worlds.World) {
	if len(c.Selected) == 0 {
		return
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		return
	}
	mx, my := ebiten.CursorPosition()
	tile := c.Cam.ScreenToTile(float64(mx), float64(my))
	if !w.Grid.InBounds(tile) {
		return
	}

	// Right-clicking an enemy issues an attack order; otherwise a move order.
	var attackTarget int
	if b := w.BuildingAt(tile); b != nil && b.Faction != c.Player {
		attackTarget = b.ID
	} else if u := w.UnitAt(tile); u != nil && u.Faction != c.Player {
		attackTarget = u.ID
	}

	for id := range c.Selected {
		u := w.Units[id]
		if u == nil || !u.IsAlive() {
			continue
		}
		if attackTarget != 0 {
			u.Order.Kind = entities.OrderAttack
			u.Order.TargetUnit = attackTarget
			u.Order.Path = nil
			u.Order.PathIdx = 0
			continue
		}
		if tile.Equal(u.Pos) {
			u.Order.Clear()
			continue
		}
		extras := paths.ThreatExtras(w, c.Player, 8.0, 4.0)
		path := paths.FindPath(w.Grid, u.Pos, tile, u.KindDef().CanFly, extras)
		if path == nil {
			continue
		}
		u.Order.Kind = entities.OrderMove
		u.Order.Target = tile
		u.Order.Path = path
		u.Order.PathIdx = 0
	}
}

// Draw renders the selection overlay (drag box).
func (c *Controller) Draw(screen *ebiten.Image) {
	if c.dragging {
		x0 := float32(min(c.dragStart.X, c.dragCur.X))
		y0 := float32(min(c.dragStart.Y, c.dragCur.Y))
		w := float32(abs(c.dragCur.X - c.dragStart.X))
		h := float32(abs(c.dragCur.Y - c.dragStart.Y))
		vector.StrokeRect(screen, x0, y0, w, h, 1, color.RGBA{0, 220, 255, 200}, true)
	}
}

func unitAtScreen(w *worlds.World, cam *renderers.Camera, faction entities.FactionID, mx, my int) *entities.Unit {
	var best *entities.Unit
	bestDist := 32.0 * 32.0 // one tile radius
	for _, u := range w.UnitsOf(faction) {
		sx, sy := cam.WorldToScreen(u.RenderX, u.RenderY)
		dx := sx - float64(mx)
		dy := sy - float64(my)
		if d := dx*dx + dy*dy; d < bestDist {
			bestDist = d
			best = u
		}
	}
	return best
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
