package scenario

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/kyeett/multiplayer-game/internal/connection"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"
	"go.uber.org/zap"
)

var _ Scenario = &SetupGame{}

type SetupGame struct {
	pasteOfferBtn  Button
	createOfferBtn Button
	pasteAnswerBtn Button
	state          string
	connection     *connection.Connection

	logger *zap.Logger

	otherCursor *gfx.Vec

	messageCallback func([]byte)
}

func (s *SetupGame) Update() error {
	x, y := ebiten.CursorPosition()
	pos := gfx.IV(x, y)
	b, err := json.Marshal(pos)
	if err == nil {
		s.connection.Send(b)
	}

	s.pasteOfferBtn.Update()
	s.createOfferBtn.Update()
	s.pasteAnswerBtn.Update()
	return nil
}

func (s *SetupGame) Draw(screen *ebiten.Image) {
	s.pasteOfferBtn.Draw(screen)
	s.createOfferBtn.Draw(screen)
	s.pasteAnswerBtn.Draw(screen)

	ebitenutil.DebugPrintAt(screen, "conn:"+s.connection.State(), 120, 10)
	ebitenutil.DebugPrintAt(screen, "data:"+s.connection.ChannelState(), 120, 30)

	if s.otherCursor != nil {
		gfx.DrawCicleFast(screen, *s.otherCursor, 10, colornames.Darkgoldenrod)
		fmt.Println("cursor drawn")
	} else {
		fmt.Println("message received")
	}
}

func NewRemoteGame() *SetupGame {
	createOfferBtn := NewButton("Create offer", 10, 10, 100, 18, colornames.Green)
	pasteAnswerBtn := NewButton("Paste answer", 10, 30, 100, 18, colornames.Darkviolet)
	pasteAnswerBtn.Hidden(true)
	pasteOfferBtn := NewButton("Paste offer", 10, 70, 100, 18, colornames.Darkred)

	s := &SetupGame{
		pasteOfferBtn:  pasteOfferBtn,
		createOfferBtn: createOfferBtn,
		pasteAnswerBtn: pasteAnswerBtn,
	}

	s.messageCallback = func(msg []byte) {
		fmt.Println("message received")
		var pos gfx.Vec
		if err := json.Unmarshal(msg, &pos); err != nil {
			fmt.Println("cursor not updated", err)
			return
		}
		s.otherCursor = &pos
		fmt.Println("cursor updated")
	}

	s.createOfferBtn.OnClick(func() {
		if err := s.copyOfferToClipboard(); err != nil {
			return
		}
		s.pasteAnswerBtn.Hidden(false)
		s.createOfferBtn.Hidden(true)
		s.pasteOfferBtn.Hidden(true)
	})

	s.pasteAnswerBtn.OnClick(func() {
		if err := s.pasteAnswerFromClipboard(); err != nil {
			return
		}
	})

	s.pasteOfferBtn.OnClick(func() {
		if err := s.pasteOfferFromClipboard(); err != nil {
			return
		}
	})

	s.logger, _ = zap.NewDevelopment()
	return s
}
