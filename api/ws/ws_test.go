package ws

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
	"services/types"
	"testing"
)

var addr = "127.0.0.1:8000"

func TestWS(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("dial error")
		t.FailNow()
	}

	message := types.WSMessage{
		Filename:   "Test file",
		Percentage: 0,
		Status:     0,
	}
	c.EnableWriteCompression(true)

	err = c.WriteJSON(&message)
	if err != nil {
		log.Error().Err(err).Msg("write error")
		t.FailNow()
	}

	err = c.ReadJSON(&message)
	if err != nil {
		log.Error().Err(err).Msg("write error")
		t.FailNow()
	}

	log.Info().
		Str("fileName", message.Filename).
		Int("percentage", message.Percentage).
		Int("status", message.Status).
		Msg("recieved message")

	_ = c.Close()

}
