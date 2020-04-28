package scenario

import "github.com/hajimehoshi/ebiten"

type Scenario interface {
	Update() error
	Draw(screen *ebiten.Image)
}
