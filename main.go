package main

import (
	//"github.com/atotto/clipboard"
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/multiplayer-game/internal/scenario"
	"log"
)

type EbitenGame struct {
	scenario scenario.Scenario
}



func (g *EbitenGame) Update(_ *ebiten.Image) error {
	g.scenario.Update()
	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.scenario.Draw(screen)
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth/2, outsideHeight/2
}

var _ ebiten.Game = &EbitenGame{}

func main() {
	g := EbitenGame{
		scenario: scenario.NewRemoteGame(),
	}
	ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
