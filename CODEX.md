# CODEX Context

Use `AGENTS.md` as the source of truth. This file exists as an explicit Codex/GPT-facing context entry.

## Repository Summary

- `Element Rush` is a Go + Ebitengine prototype.
- It is a desktop-first 2D infinite runner platformer.
- The current prototype already builds and starts successfully.

## Mechanics Already Implemented

- auto-runner movement
- jump with coyote time and jump buffering
- three attacks: fire, ice, thunder
- weakness cycle: `fire > ice > thunder > fire`
- elemental enemies
- hostile projectiles and player deflects
- parallax backgrounds and HUD

## Commands

```bash
go run ./cmd/game
go build ./...
go test ./...
python tools/genassets/generate_assets.py
```

## Files To Know First

- `internal/game/game.go`: most gameplay logic lives here right now
- `internal/game/assets.go`: loader and spritesheet validation
- `tools/genassets/generate_assets.py`: generates the original pixel-art runtime assets
- `docs/assets.md`: asset notes

## Constraints

- Keep assets original and in-repo.
- Maintain visual clarity for fire, ice, and thunder.
- Update asset loader assumptions if generator output changes.
- Prefer small, safe edits unless a refactor is clearly warranted.

## Missing / Future Work

- audio
- more gameplay content
- playtesting and balance tuning
- better internal modularization if the codebase expands
