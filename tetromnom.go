package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"tetromnom/game"
)

func main() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Tetromnom")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
