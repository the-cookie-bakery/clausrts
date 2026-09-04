# Roadmap & Milestones — CLAUSRTS

Delivery order is strictly milestone-sequential; each milestone must build,
run and pass tests before the next starts.

**Status:** M0–M3 are **implemented** and playable (a free sandbox skirmish vs.
the AI). M4 (teaching layer + scripted scenarios) and M5 (full campaign) are
planned but not yet built. See [architecture.md](architecture.md) for the
current code layout.

## M0 — Skeleton (engine proof) ✅

**Goal:** prove Linux + ebitengine build and window.

- `cmd/clausrts` opens a window, runs the ebiten `Game` loop.
- `configs`, `helpers` in place.
- `assets` package with `go:embed` loading a generated placeholder PNG.
- `cmd/genassets` regenerates placeholder assets.

**Acceptance:** `make run` opens a window; `make test` and `make build` pass.

## M1 — World, grid, camera, movement ✅

**Goal:** a navigable map.

- `worlds`: grid map with terrain, factions, unit placement, fixed tick.
- `paths`: AStar with terrain cost; 8-direction movement.
- `inputs`: selection (click/drag), right-click move orders, hotkeys, pause.
- `renderers`: camera (pan/zoom), tile rendering, unit rendering with rotation.
- Grid-locked movement with smooth interpolation and facing.

**Acceptance:** select a unit, right-click a tile, unit pathfinds and moves;
camera pans/zooms; `Space` pauses.

## M2 — Clausewitz properties + Trinity UI ✅

**Goal:** the concepts become numeric, observable state.

- `entities.Dynamics` (friction, focus, gravity, fog, moral, supply, culmination).
- `systems/dynamics.go` (friction/focus/gravity/culmination), `systems/fog.go`,
  `systems/trinity.go`.
- `uis`: Trinity bars and a selection panel showing live `Dynamics`.
- Teaching layer scaffolding (quote on concept change) — **pending (M4)**.

**Acceptance:** moving a unit away from HQ/into fog raises its friction;
concentrating units raises focus; the three bars drift and are displayed.

## M3 — Combat & hybrid counters ✅

**Goal:** a playable conflict.

- `systems/combat.go` (damage, range, counters, morale), `systems/trinity.go`
  (supply-driven Government bar), `systems/movement.go` (move/attack orders).
- All hybrid unit kinds + building kinds wired with stats and the counter matrix.
- `ai`: enemy commander advances on the player's Schwerpunkt and engages.
- No projectiles/reinforcement yet (planned).

**Acceptance:** two factions fight; counters apply (air defense 2× vs drones,
cyber 2× vs radar/depot, bunker shields its tile); AI is a functional opponent.

## M4 — Teaching & scenarios

**Goal:** scenarios with scripted lessons and quotes.

- `scenarios`: objectives, triggers, win/lose conditions, intro text.
- `data`: quotes per concept (German, from `references/vomkriege_book`).
- Teaching layer fully wired to concept events.
- Scenarios 1–2 ("Was ist der Krieg?", "Friktion im Kriege") playable.

**Acceptance:** a scenario has an objective, can be won/lost, and surfaces the
correct quote when a concept fires.

## M5 — Campaign & polish

- Scenarios 3–6 (Kulminationspunkt, Schwerpunkt, Ökonomie der Kräfte, Kriegsplan).
- Balance pass over the counter matrix and culmination curve.
- Final docs, `docs/` reference, README, lint clean.

**Acceptance:** full 6-scenario campaign playable end-to-end.
