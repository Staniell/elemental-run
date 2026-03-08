package game

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"

	assetfs "go-2d/internal/assets"
)

type gameAssets struct {
	backgroundFar  *ebiten.Image
	backgroundMid  *ebiten.Image
	backgroundNear *ebiten.Image

	playerFrames []*ebiten.Image
	enemyFrames  map[enemyKind][]*ebiten.Image

	playerProjectiles map[element]*ebiten.Image
	enemyProjectiles  map[element]*ebiten.Image
	attackFX          map[element]*ebiten.Image
	tiles             map[tileKind]*ebiten.Image
	icons             map[element]*ebiten.Image
	heart             *ebiten.Image
}

func loadAssets() (*gameAssets, error) {
	playerSheet, err := loadPNG("player.png")
	if err != nil {
		return nil, err
	}
	enemySheet, err := loadPNG("enemies.png")
	if err != nil {
		return nil, err
	}
	projectileSheet, err := loadPNG("projectiles.png")
	if err != nil {
		return nil, err
	}
	tileSheet, err := loadPNG("tiles.png")
	if err != nil {
		return nil, err
	}
	fxSheet, err := loadPNG("fx.png")
	if err != nil {
		return nil, err
	}
	uiSheet, err := loadPNG("ui.png")
	if err != nil {
		return nil, err
	}
	backgroundFar, err := loadPNG("background_far.png")
	if err != nil {
		return nil, err
	}
	backgroundMid, err := loadPNG("background_mid.png")
	if err != nil {
		return nil, err
	}
	backgroundNear, err := loadPNG("background_near.png")
	if err != nil {
		return nil, err
	}

	playerFrames := sliceHorizontal(playerSheet, 48, 48)
	enemyFrames := sliceHorizontal(enemySheet, 32, 32)
	projectileFrames := sliceHorizontal(projectileSheet, 16, 16)
	tileFrames := sliceHorizontal(tileSheet, 32, 32)
	fxFrames := sliceHorizontal(fxSheet, 32, 32)
	uiFrames := sliceHorizontal(uiSheet, 16, 16)

	if len(playerFrames) < 8 {
		return nil, fmt.Errorf("player spritesheet has %d frames, need at least 8", len(playerFrames))
	}
	if len(enemyFrames) < 9 {
		return nil, fmt.Errorf("enemy spritesheet has %d frames, need at least 9", len(enemyFrames))
	}
	if len(projectileFrames) < 6 {
		return nil, fmt.Errorf("projectile spritesheet has %d frames, need at least 6", len(projectileFrames))
	}
	if len(tileFrames) < 4 {
		return nil, fmt.Errorf("tile spritesheet has %d frames, need at least 4", len(tileFrames))
	}
	if len(fxFrames) < 3 {
		return nil, fmt.Errorf("fx spritesheet has %d frames, need at least 3", len(fxFrames))
	}
	if len(uiFrames) < 4 {
		return nil, fmt.Errorf("ui spritesheet has %d frames, need at least 4", len(uiFrames))
	}

	return &gameAssets{
		backgroundFar:  backgroundFar,
		backgroundMid:  backgroundMid,
		backgroundNear: backgroundNear,
		playerFrames:   playerFrames,
		enemyFrames: map[enemyKind][]*ebiten.Image{
			enemyGround: enemyFrames[0:3],
			enemyTurret: enemyFrames[3:6],
			enemyFlyer:  enemyFrames[6:9],
		},
		playerProjectiles: map[element]*ebiten.Image{
			fire:    projectileFrames[0],
			ice:     projectileFrames[1],
			thunder: projectileFrames[2],
		},
		enemyProjectiles: map[element]*ebiten.Image{
			fire:    projectileFrames[3],
			ice:     projectileFrames[4],
			thunder: projectileFrames[5],
		},
		attackFX: map[element]*ebiten.Image{
			fire:    fxFrames[0],
			ice:     fxFrames[1],
			thunder: fxFrames[2],
		},
		tiles: map[tileKind]*ebiten.Image{
			tileGround:   tileFrames[0],
			tilePlatform: tileFrames[1],
			tileCrate:    tileFrames[2],
			tileSpike:    tileFrames[3],
		},
		icons: map[element]*ebiten.Image{
			fire:    uiFrames[1],
			ice:     uiFrames[2],
			thunder: uiFrames[3],
		},
		heart: uiFrames[0],
	}, nil
}

func loadPNG(name string) (*ebiten.Image, error) {
	bs, err := assetfs.FS.ReadFile("generated/" + name)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func sliceHorizontal(img *ebiten.Image, frameW, frameH int) []*ebiten.Image {
	w, h := img.Size()
	if frameW <= 0 || frameH <= 0 || h < frameH || w < frameW {
		return nil
	}
	count := w / frameW
	frames := make([]*ebiten.Image, 0, count)
	for i := 0; i < count; i++ {
		r := image.Rect(i*frameW, 0, (i+1)*frameW, frameH)
		frames = append(frames, img.SubImage(r).(*ebiten.Image))
	}
	return frames
}
