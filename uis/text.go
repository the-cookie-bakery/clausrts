package uis

import (
	"bytes"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	faceSource     *text.GoTextFaceSource
	faceSourceOnce sync.Once
	faceSourceErr  error
)

func fontSource() *text.GoTextFaceSource {
	faceSourceOnce.Do(func() {
		faceSource, faceSourceErr = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	})
	return faceSource
}

// Face returns a scalable text.Face at the given pixel size.
func Face(size float64) text.Face {
	return &text.GoTextFace{Source: fontSource(), Size: size}
}

// DrawText draws a left-aligned string at (x, y) in the given colour.
func DrawText(screen *ebiten.Image, str string, x, y, size float64, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, Face(size), op)
}

// TextWidth returns the rendered width of str at the given size.
func TextWidth(str string, size float64) float64 {
	return text.Advance(str, Face(size))
}
