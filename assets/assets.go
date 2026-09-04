// Package assets embeds the game's 2D art and exposes it by key. Each entity
// kind has a matching PNG: entities/tank.go <-> assets/tank.png. Loading goes
// through assets.Image so placeholder art can be swapped for real art with no
// code changes.
package assets

import (
	"embed"
	"io/fs"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed *.png terrain/*.png
var files embed.FS

var (
	images map[string]*ebiten.Image
	once   sync.Once
)

// Load decodes every embedded PNG into GPU textures. It is safe to call from
// the game loop (deferred until the graphics context exists); subsequent calls
// are no-ops.
func Load() {
	once.Do(func() {
		images = map[string]*ebiten.Image{}
		err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			img, _, err := ebitenutil.NewImageFromFileSystem(files, path)
			if err != nil {
				log.Printf("assets: failed to load %s: %v", path, err)
				return nil
			}
			images[path] = img
			return nil
		})
		if err != nil {
			log.Printf("assets: walk error: %v", err)
		}
		log.Printf("assets: loaded %d images", len(images))
	})
}

// Image returns the decoded texture for name (e.g. "tank.png"), or nil.
func Image(name string) *ebiten.Image {
	return images[name]
}

// TerrainImage returns the texture for a terrain tile (e.g. "grass").
func TerrainImage(name string) *ebiten.Image {
	return images["terrain/"+name+".png"]
}
