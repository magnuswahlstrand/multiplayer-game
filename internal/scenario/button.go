package scenario

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/peterhellberg/gfx"
	"image/color"
	"sync"
)

type Button struct {
	label         string
	pos           gfx.Vec
	width, height float64
	color         color.Color
	hidden bool

	onClick       func()

	mu     *sync.Mutex
}

func (b *Button) Draw(screen *ebiten.Image) {
	if b.hidden {
		return
	}

	ebitenutil.DrawRect(screen, b.pos.X, b.pos.Y, b.width, b.height, b.color)
	ebitenutil.DebugPrintAt(screen, b.label, int(b.pos.X+2), int(b.pos.Y))
}

func (b *Button) Hidden(v bool) {
	fmt.Println("update state for ", b.label, " to ", v)
	b.hidden = v
}

func cursorPosition() gfx.Vec {
	x, y := ebiten.CursorPosition()
	return gfx.IV(x, y)
}

func (b *Button) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		bounds := gfx.R(0, 0, b.width, b.height).Moved(b.pos)

		if bounds.Contains(cursorPosition()) {
			b.mu.Lock()
			defer b.mu.Unlock()
			b.onClick()
		}
	}
}

func (b *Button) OnClick(f func()) {
	b.mu.Lock()
	b.onClick = f
	b.mu.Unlock()
}

func NewButton(label string, x, y, width, height float64, c color.Color) Button {
	return Button{
		label:   label,
		pos:     gfx.V(x, y),
		width:   width,
		height:  height,
		color:   c,
		onClick: func() {},
		mu: &sync.Mutex{},
	}
}
