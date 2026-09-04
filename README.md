# CLAUSRTS

A 2D, tile-based RTS (real-time with pause) that teaches the strategic
philosophy of Carl von Clausewitz's *Vom Kriege* through a modern hybrid-war
setting. Built with [ebitengine](https://github.com/hajimehoshi/ebiten) for Linux.

The game models war as Clausewitz described it — not as an economy, but as a
continuation of politics, shaped by **friction, fog, focus, gravity
(Schwerpunkt) and the culminating point**. The three Trinity bars
(People / Military / Government) are the heart of the interface.

## Run

```sh
make run     # build & run
make build   # binary into ./bin/
make test    # unit tests
make assets  # regenerate placeholder PNGs
```

## Controls

- **Left-click / drag** — select your units
- **Right-click** — move; right-click an enemy to attack it
- **WASD / arrows / middle-drag** — pan, **scroll** — zoom
- **Space** — pause (plan while paused, resolve in real time)
- **Esc** — clear selection

## Documentation

See [docs/](docs/README.md) for the design, architecture, roadmap, Clausewitz
mapping and asset conventions.
