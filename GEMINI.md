# GEMINI Context

Use `AGENTS.md` as the canonical project context. This file is a shorter Gemini-friendly mirror.

## Project Snapshot

- Name: `Element Rush`
- Language: Go
- Engine: Ebitengine
- Genre: 2D desktop infinite runner platformer
- Status: working prototype with a main menu, successful `go build ./...`, `go test ./...`, and startup smoke test

## Core Mechanics

- player auto-runs
- jump with coyote time and jump buffering
- attacks: fire, ice, thunder
- weakness loop: `fire > ice > thunder > fire`
- enemies can be struck directly or by deflected projectiles
- hostile projectiles can be deflected by attacks

## Main Files

- `cmd/game/main.go`
- `internal/game/game.go`
- `internal/game/assets.go`
- `tools/genassets/generate_assets.py`
- `docs/assets.md`

## Commands

```bash
go run ./cmd/game
go build ./...
go test ./...
python tools/genassets/generate_assets.py
```

## Guidance

- Preserve desktop-first pixel-art direction.
- Prefer original generated assets over external downloads.
- Keep gameplay readable and responsive.
- If changing spritesheets, update the expected frame counts in `internal/game/assets.go`.
- If the project grows, splitting `internal/game/game.go` by system is reasonable.

## Good Next Steps

- add SFX/music
- add more chunk/enemy variety
- tune difficulty and collision feel
- add gameplay-focused tests or deterministic simulation hooks
