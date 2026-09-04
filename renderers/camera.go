// Package renderers turns world state into pixels: camera, tiles, units,
// buildings, fog-of-war overlay and selection.
package renderers

import (
	"clausrts/configs"
	"clausrts/helpers"
)

// Camera maps world pixel coordinates to screen coordinates.
type Camera struct {
	X, Y    float64 // world-pixel coordinate of the screen's top-left
	Zoom    float64
	ScreenW int
	ScreenH int
}

// NewCamera returns a camera at origin with zoom 1.
func NewCamera() *Camera { return &Camera{Zoom: 1} }

// SetScreen records the current viewport size.
func (c *Camera) SetScreen(w, h int) { c.ScreenW, c.ScreenH = w, h }

// WorldToScreen converts world pixel coordinates to screen coordinates.
func (c *Camera) WorldToScreen(wx, wy float64) (float64, float64) {
	return (wx - c.X) * c.Zoom, (wy - c.Y) * c.Zoom
}

// ScreenToWorld converts screen coordinates to world pixel coordinates.
func (c *Camera) ScreenToWorld(sx, sy float64) (float64, float64) {
	return sx/c.Zoom + c.X, sy/c.Zoom + c.Y
}

// ScreenToTile converts screen coordinates to a tile coordinate.
func (c *Camera) ScreenToTile(sx, sy float64) helpers.Point {
	wx, wy := c.ScreenToWorld(sx, sy)
	return helpers.Point{X: int(wx) / configs.TileSize, Y: int(wy) / configs.TileSize}
}

// TileToScreen returns the screen pixel of a tile's top-left corner.
func (c *Camera) TileToScreen(p helpers.Point) (float64, float64) {
	return c.WorldToScreen(float64(p.X*configs.TileSize), float64(p.Y*configs.TileSize))
}

// Pan moves the camera by a screen-space delta.
func (c *Camera) Pan(sx, sy float64) {
	c.X -= sx / c.Zoom
	c.Y -= sy / c.Zoom
}

// CenterOn places the camera so the world point is centered on screen.
func (c *Camera) CenterOn(wx, wy float64) {
	c.X = wx - float64(c.ScreenW)/2/c.Zoom
	c.Y = wy - float64(c.ScreenH)/2/c.Zoom
}

// ZoomAt zooms by factor keeping the world point under (sx, sy) fixed.
func (c *Camera) ZoomAt(sx, sy, factor float64) {
	wx, wy := c.ScreenToWorld(sx, sy)
	c.Zoom *= factor
	if c.Zoom < configs.MinZoom {
		c.Zoom = configs.MinZoom
	}
	if c.Zoom > configs.MaxZoom {
		c.Zoom = configs.MaxZoom
	}
	c.X = wx - sx/c.Zoom
	c.Y = wy - sy/c.Zoom
}

// ClampToWorld keeps the camera within the world bounds so the view cannot
// pan into the void.
func (c *Camera) ClampToWorld(gridW, gridH int) {
	worldW := float64(gridW * configs.TileSize)
	worldH := float64(gridH * configs.TileSize)

	if c.X < 0 {
		c.X = 0
	}
	if c.Y < 0 {
		c.Y = 0
	}
	maxX := worldW - float64(c.ScreenW)/c.Zoom
	maxY := worldH - float64(c.ScreenH)/c.Zoom
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}
	if c.X > maxX {
		c.X = maxX
	}
	if c.Y > maxY {
		c.Y = maxY
	}
}
