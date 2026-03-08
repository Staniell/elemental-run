package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"go-2d/internal/game"
)

func main() {
	g, err := game.New()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(game.ScreenWidth*3, game.ScreenHeight*3)
	ebiten.SetWindowTitle("Element Rush")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
