# Architecture — CLAUSRTS

## 1. Package layout (plural names)

```
clausrts/
├── cmd/
│   ├── clausrts/main.go        # entrypoint: window setup, run game
│   └── genassets/main.go       # regenerate placeholder PNG assets
├── configs/                    # constants: tile size (32), window, TPS, balance
├── worlds/                     # World state: grid map, factions, formations, fixed tick
├── entities/                   # Unit, Formation, Building, Faction, kinds, Dynamics
├── systems/                    # simulation systems (one concern each)
├── paths/                      # AStar over the 32px grid
├── commanders/                 # Commander interface (player + AI)
├── ai/                         # enemy commander
├── renderers/                  # camera, tiles, units, fog overlay, selection
├── uis/                        # HUD, Trinity bars, minimap, text, teaching layer
├── inputs/                     # mouse/keyboard, selection, orders, hotkeys, pause
├── scenarios/                  # scenario definitions, objectives, win/lose
├── data/                       # static definitions (units/buildings/terrain/quotes)
├── assets/                     # go:embed 2D art (PNG) + font
├── helpers/                    # grid math, geometry
├── utils/                      # deterministic RNG, timing, logging
├── docs/                       # this folder
└── tests/                      # integration/scenario tests
```

## 2. Data flow

```
inputs (mouse/kb) ──► player orders ──► systems ──► worlds (state)
ai (enemy commander) ──► orders ─────────────┘
renderers ◄── worlds (read)
uis       ◄── worlds (read)
```

- `worlds.World` owns all mutable state (grid, units, formations, buildings,
  factions, Trinity bars, fog grid).
- `systems` are stateless functions run once per tick by `systems.Tick` (fixed
  order: fog → dynamics → movement → combat → trinity).
- `inputs.Controller` turns raw input into selection + `Order`s for the player;
  `ai.EnemyCommander` does the same for the enemy. Both write orders directly
  onto units; the `commanders` package defines the `Commander` interface for a
  future refactor into pluggable commanders.
- `renderers`/`uis` are read-only views over `worlds`.

## 3. Coordinate / grid model

- Tile size **32 px** (`configs.TileSize`).
- World is `Width × Height` tiles. Tiles are `[X, Y]` with `(0,0)` top-left.
- `helpers.Point` (tile coordinates) plus `DirectionTo`, `Chebyshev`,
  `Manhattan` and clamp helpers. Pixel↔tile conversion lives in
  `renderers.Camera` (world-pixel ↔ screen) and `entities.TileCenter`
  (tile → world-pixel center).
- Units are grid-locked: they occupy exactly one tile and move tile-to-tile
  with smooth interpolation (a `MoveFrom`/`MoveTo` pair + `MoveProgress`).
- Facing is an 8-direction unit vector in tile space.
- Passability is **per-unit**: `worlds.Grid.PassableFor(p, flying)` lets flying
  units (drones, `UnitKind.CanFly`) cross water; mountains stay impassable.

## 4. Pathfinding (`paths`)

AStar on the grid with a cost function that deliberately encodes Clausewitz:

```
cost(node) = terrainCost(tile)
           + threatWeight * enemyGravity(node)   (route around enemy Schwerpunkte)
```

- `paths.FindPath(grid, from, to, flying, extras)` returns a tile path;
  `extras` (typically `paths.ThreatExtras`) adds enemy-gravity routing cost.
- Friction is applied elsewhere: it **slows** movement (`systems/movement.go`)
  and causes combat misses, rather than altering path cost.
- Flying units (drones) cross water; mountains are impassable to all units.

This makes gravity *observable*: units visibly route around enemy
concentrations.

## 5. Simulation loop & determinism

- Fixed timestep: ebiten `Update` runs at `configs.TicksPerSecond` (60);
  `systems.Tick` advances one tick of `configs.SecondsPerTick` each call.
- **Pause** (`Space`) stops `systems.Tick` and the AI step but keeps rendering
  + input live (Spannung und Ruhe: plan while paused, resolve in real time).
- `utils.RNG` is a seeded PRNG (World-owned) so replays/tests are deterministic.

## 6. Rendering & rotation

- `renderers.Camera` provides world↔screen transform (pan = WASD / arrows /
  middle-drag, zoom = scroll), clamped to the map.
- Terrain tiles are drawn from `assets` (embedded PNG, 32×32).
- Units are drawn from `assets` and **rotated toward their facing tile**; the
  per-kind hook is `drawUnits` in `renderers` (kinds such as `tank` will later
  draw a static hull + rotating turret — see [assets.md](assets.md)).
- **Health bars**: damaged units and buildings draw a colour-coded bar (green →
  amber → red) with the numeric HP value above them.
- **Attack radius**: selected units draw a faint filled circle (alpha 12) plus
  an outline in the selection colour, radius = `kind.Range × tile`.
- Fog-of-war is drawn as a translucent overlay per tile from the per-faction
  fog grid; enemy entities outside `VisVisible` are hidden.

## 7. Assets & embedding

- `assets` embeds `assets/` via `//go:embed`, exposing
  `assets.Image(key string) *ebiten.Image`.
- Conventions are documented in [assets.md](assets.md).

## 8. UI (`uis`)

- Trinity bars (People / Military / Government) for both factions at the top.
- Selection panel: the selected unit's kind, HP and live `Dynamics` values.
- Status text (tick counter, paused indicator).
- Text uses `golang.org/x/image/font/gofont/goregular` via
  `ebiten/v2/text/v2` (scalable, embedded — no font asset needed for now).
- **Not yet implemented:** minimap, teaching layer (contextual German quotes).
