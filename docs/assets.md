# Asset Pipeline

All shipped art for this prototype is generated in-repo and is intended to be original.

## Visual Direction

- desktop-first pixel art at a low internal resolution
- warm dusk sky with layered parallax silhouettes
- readable elemental color language
  - fire: orange and gold
  - ice: cyan and white
  - thunder: yellow and cream
- dark ink outlines with restrained shading for gameplay clarity

## Output Files

- `internal/assets/generated/player.png`
- `internal/assets/generated/enemies.png`
- `internal/assets/generated/projectiles.png`
- `internal/assets/generated/tiles.png`
- `internal/assets/generated/fx.png`
- `internal/assets/generated/ui.png`
- `internal/assets/generated/background_far.png`
- `internal/assets/generated/background_mid.png`
- `internal/assets/generated/background_near.png`

## Source Script

The art generator script is `tools/genassets/generate_assets.py`.

Run it from the repository root:

```bash
python tools/genassets/generate_assets.py
```

This keeps the prototype self-contained and avoids external art licensing issues.
