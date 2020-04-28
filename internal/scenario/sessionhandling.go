package scenario

import (
	"github.com/atotto/clipboard"
	"github.com/kyeett/multiplayer-game/internal/connection"
	"github.com/pion/webrtc/v2"
)

func (s *SetupGame) pasteOfferFromClipboard() error {
	s.logger.Info("load session description from clipboard")
	d, err := loadSessionDescription()
	if err != nil {
		return err
	}

	c, err := connection.CreateForAnswer(s.messageCallback)
	if err != nil {
		return err
	}
	s.connection = c

	if err = s.connection.SetRemoteDescription(d); err != nil {
		return err
	}

	answer, err := s.connection.CreateAnswer(nil)
	if err != nil {
		return err
	}

	err = s.connection.SetLocalDescription(answer)
	if err != nil {
		return err
	}

	return copySessionDescription(answer)
}

func (s *SetupGame) copyOfferToClipboard() error {
	s.logger.Info("copy session description to clipboard")
	c, err := connection.CreateForOffer(s.messageCallback)
	if err != nil {
		return nil
	}
	s.connection = c

	offer, err := s.connection.CreateOffer(nil)
	if err != nil {
		return err
	}

	if err := s.connection.SetLocalDescription(offer); err != nil {
		return err
	}

	return copySessionDescription(offer)
}


func (s *SetupGame) pasteAnswerFromClipboard() error {
	s.logger.Info("load session description from clipboard")
	d, err := loadSessionDescription()
	if err != nil {
		return err
	}

	if err = s.connection.SetRemoteDescription(d); err != nil {
		return err
	}
	return nil
}


func copySessionDescription(sd webrtc.SessionDescription) error {
	e, err := connection.Encode(sd)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(e); err != nil {
		return err
	}
	return nil
}

func loadSessionDescription() (webrtc.SessionDescription, error) {
	in, err := clipboard.ReadAll()
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	var d webrtc.SessionDescription
	if err := connection.Decode(in, &d); err != nil {
		return webrtc.SessionDescription{}, err
	}
	return d, nil
}