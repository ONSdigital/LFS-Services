package ws

import (
	"encoding/json"
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

	defer c.Close()

	m := types.WSMessage{
		Filename:   "Test file",
		Percentage: 0,
		Status:     0,
	}
	b, err := json.Marshal(m)
	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Error().Err(err).Msg("write error")
		t.FailNow()
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("write error")
		t.FailNow()
	}
	log.Printf("recv: %s", message)

}
