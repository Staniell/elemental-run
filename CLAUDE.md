# CLAUDE Context

Use `AGENTS.md` as the primary source of truth for this repository.

## Quick Project Summary

- Project: `Element Rush`
- Stack: Go + Ebitengine
- Type: desktop-first infinite runner platformer
- Status: playable prototype with a keyboard-driven main menu, generated original pixel art, and verified build startup

## What Exists

- auto-run movement
- jump with coyote time and jump buffer
- fire / ice / thunder attacks
- weakness cycle: `fire > ice > thunder > fire`
- elemental enemies, hostile projectiles, and projectile deflects
- parallax backgrounds, HUD, score, kill count, main menu, and restart flow

## Where To Work

- `internal/game/game.go`: gameplay and rendering core
- `internal/game/assets.go`: embedded asset loading and frame validation
- `tools/genassets/generate_assets.py`: generated art pipeline

## Commands

```bash
go run ./cmd/game
go build ./...
go test ./...
python tools/genassets/generate_assets.py
```

## Important Constraints

- Keep art original and in-repo.
- Keep the three elements visually distinct.
- If sprite sheet dimensions change, update loaders in `internal/game/assets.go`.
- Prefer focused edits over broad rewrites unless refactoring is necessary.

## Known Missing Pieces

- no real audio yet
- no gameplay tests yet
- main gameplay file is still monolithic by design for the prototype
