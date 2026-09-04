// Command clausrts is the game entrypoint. It wires the world, camera, input,
// systems, AI and UI into an ebiten Game and runs it.
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"clausrts/ai"
	"clausrts/assets"
	"clausrts/configs"
	"clausrts/entities"
	"clausrts/inputs"
	"clausrts/renderers"
	"clausrts/scenarios"
	"clausrts/systems"
	"clausrts/uis"
	"clausrts/worlds"
)

const playerFaction = entities.FactionID(configs.FactionPlayer)

// Game implements ebiten.Game.
type Game struct {
	world  *worlds.World
	cam    *renderers.Camera
	input  *inputs.Controller
	enemy  *ai.EnemyCommander
	sw, sh int
	centered bool
}

// NewGame assembles a fresh sandbox game.
func NewGame() *Game {
	w := scenarios.Sandbox(42)
	cam := renderers.NewCamera()
	return &Game{
		world: w,
		cam:   cam,
		input: inputs.NewController(cam, playerFaction),
		enemy: ai.NewEnemyCommander(entities.FactionID(configs.FactionEnemy)),
	}
}

func (g *Game) Update() error {
	assets.Load()

	g.input.Update(g.world)
	if !g.world.Paused {
		g.enemy.Step(g.world)
	}
	systems.Tick(g.world)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	renderers.Draw(screen, g.world, g.cam, playerFaction, g.input.Selected)
	g.input.Draw(screen)
	uis.DrawHUD(screen, g.world, playerFaction, g.input.Selected, g.world.Paused, g.sw, g.sh)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	g.sw, g.sh = outsideW, outsideH
	g.cam.SetScreen(outsideW, outsideH)
	if !g.centered {
		g.centered = true
		g.cam.CenterOn(4*configs.TileSize+configs.TileSize/2, 4*configs.TileSize+configs.TileSize/2)
	}
	return outsideW, outsideH
}

func main() {
	ebiten.SetWindowSize(configs.WindowWidth, configs.WindowHeight)
	ebiten.SetWindowTitle("CLAUSRTS")
	ebiten.SetTPS(configs.TicksPerSecond)
	ebiten.SetVsyncEnabled(true)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
