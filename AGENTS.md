# Element Rush Agent Context

This file is the generic project context for coding agents working in this repository.

## Project

- Name: `Element Rush`
- Stack: Go + Ebitengine (`github.com/hajimehoshi/ebiten/v2`)
- Genre: desktop-first 2D infinite runner platformer
- Current state: playable prototype with original in-repo pixel art, a main menu, chunk spawning, elemental combat, enemies, projectiles, deflects, HUD, and restart flow

## Current Gameplay

- The game starts from a keyboard-driven main menu with `Start` and `Quit`.
- The player auto-runs to the right.
- Jump uses coyote time and jump buffering.
- Attacks are `fire`, `ice`, and `thunder`.
- Weakness triangle: `fire > ice > thunder > fire`.
- Enemy shots can be deflected by any active attack.
- Deflected shots become player-owned and inherit the attack element used for the deflect.

## Controls

- `Space`, `W`, `Up`: jump
- `Up`, `Down`, `W`, `S`: navigate the main menu
- `Enter` or `Space`: confirm a menu selection
- `J` or `1`: fire
- `K` or `2`: ice
- `L` or `3`: thunder
- `Esc`: return to menu during play or quit from the main menu
- `R`: restart after defeat

## Verified Status

- `go build ./...` passes
- `go test ./...` passes
- startup smoke test for `bin/element-rush.exe` passed without an immediate crash

## Important Files

- `cmd/game/main.go`: Ebitengine window/bootstrap
- `internal/game/game.go`: main gameplay implementation; currently holds almost all runtime logic
- `internal/game/assets.go`: asset loading and spritesheet slicing/validation
- `internal/assets/assets.go`: embedded generated art files
- `internal/assets/generated/*.png`: runtime sprites/backgrounds
- `tools/genassets/generate_assets.py`: original in-repo pixel art generator
- `docs/assets.md`: asset pipeline notes
- `README.md`: high-level project overview and run instructions

## Asset Rules

- Keep shipped art original and in-repo.
- Runtime art is generated into `internal/assets/generated`.
- If you change frame counts or sheet dimensions, update `internal/game/assets.go` to match.
- If you change the art style, preserve readability of the three elements:
  - fire = orange/gold
  - ice = cyan/white
  - thunder = yellow/cream

## Architecture Notes

- Current code is intentionally simple and prototype-oriented.
- `internal/game/game.go` contains player logic, collision, enemies, chunk spawning, combat, rendering, HUD, and utility helpers in one file.
- If you do major feature work, it is reasonable to split systems into multiple files, but keep behavior stable while refactoring.
- The logical resolution is `480x270`, scaled up in the window.

## Current Content

- Enemy types:
  - ground enemy (`ice` affinity)
  - turret enemy (`fire` affinity)
  - flyer enemy (`thunder` affinity)
- Terrain/content pieces:
  - ground tiles
  - elevated platforms
  - crates used as step-up obstacles
  - spike hazards
- Chunk patterns currently include breather, crate-step, spike-lane, turret-roost, sky-bridge, and mixed-gauntlet layouts.

## Working Conventions

- Prefer incremental edits over large rewrites.
- Preserve desktop-first pixel-art direction unless explicitly asked to change it.
- Keep gameplay readable before adding complexity.
- Avoid adding external art dependencies when the same result can be kept in the generator script.
- If you add assets, document them in `docs/assets.md`.

## Known Gaps / Good Next Steps

- No real audio/SFX yet.
- No gameplay tests yet; current `go test ./...` only validates package compilation.
- Game logic could be split out of `internal/game/game.go` as the project grows.
- Difficulty, collision feel, and balance still need manual tuning through playtesting.
