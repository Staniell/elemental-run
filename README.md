# Element Rush

`Element Rush` is a desktop-first 2D infinite runner platformer written in Go with Ebitengine.

## Mechanics

- auto-running platformer movement
- jump timing with coyote time and jump buffering
- three elemental attacks: fire, ice, thunder
- weakness cycle: fire beats ice, ice beats thunder, thunder beats fire
- enemies with elemental affinities
- hostile projectiles that can be deflected back at enemies
- layered parallax backgrounds and original in-repo pixel art

## Planned Controls

- `Up` / `Down` or `W` / `S`: menu navigation
- `Enter` or `Space`: confirm menu choice
- `Space`, `W`, or `Up`: jump
- `J` or `1`: fire
- `K` or `2`: ice
- `L` or `3`: thunder
- `Esc`: return to menu during a run, or quit from the main menu
- `R`: restart after defeat

## Run

```bash
go run ./cmd/game
```

## Assets

The generated sprite sheets and backgrounds live in `internal/assets/generated`.
The asset pipeline notes live in `docs/assets.md`.
