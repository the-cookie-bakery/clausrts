# CLAUSRTS — Documentation

A 2D, tile-based RTS (real-time with pause) that teaches the strategic
philosophy of Carl von Clausewitz's *Vom Kriege* (On War) through a modern
hybrid-war setting.

## Document index

| Document | Contents |
|---|---|
| [design.md](design.md) | Vision, learning objectives, gameplay loop, Clausewitz → mechanic mapping, Trinity bars, hybrid-war counter matrix |
| [architecture.md](architecture.md) | Package layout, data flow, coordinate/grid model, determinism, rendering & rotation |
| [roadmap.md](roadmap.md) | Milestones M0–M5 with acceptance criteria and delivery order |
| [clausewitz.md](clausewitz.md) | Concept reference with original German quotes and their game translation |
| [assets.md](assets.md) | Asset conventions, `entities/<name>.go` ↔ `assets/<name>.png` correlation, rotation rules |

## Quickstart

```sh
make run     # build & run the game (Linux)
make build   # build the binary into ./bin/
make test    # run unit tests (AStar, systems, entities)
make lint    # golangci-lint (if installed)
make assets  # regenerate placeholder PNG assets (cmd/genassets)
```

## Key facts

- **Engine:** `github.com/hajimehoshi/ebiten/v2` (see `go.mod`; the `references/ebiten`
  snapshot is for source-reading only and is not part of the module).
- **Platform:** Linux only (for now).
- **Tile size:** 32 px (`configs.TileSize`).
- **Movement:** grid-locked, 8-directional, AStar pathfinding whose cost function
  includes terrain and enemy gravity/threat (`paths` package). Drones fly over
  water.
- **UI language:** English; **scenario/quotes:** original German (from
  `references/vomkriege_book`).

## Controls

- **Left-click / drag** — select your units (Shift-click to add/remove)
- **Right-click** — move; right-click an enemy unit/building to attack it
- **WASD / arrows / middle-drag** — pan; **scroll** — zoom
- **Space** — pause (plan while paused, resolve in real time)
- **Esc** — clear selection

## Current status

M0–M3 are implemented and playable as a free sandbox skirmish against a simple
AI (`scenarios.Sandbox`): grid + camera + AStar movement, fog-of-war, the
Clausewitz `Dynamics` (friction/focus/gravity/fog/moral/culmination), the three
Trinity bars, hybrid-unit combat with counters, health bars, and attack-radius
circles on selection. M4 (scripted scenarios + teaching layer) and M5 (full
campaign) are pending — see [roadmap.md](roadmap.md).

The primary teaching device is that the game's *numeric state* (friction, focus,
gravity, fog, culmination, morale) is the exact set of ideas Clausewitz describes,
and the UI surfaces the relevant quote whenever one of those values changes
meaningfully.
