package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"services/types"
)

type WebSocketHandler struct {
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (wsh WebSocketHandler) ServeWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Error().
				Err(err).
				Msg("WebSocket handshake error")
		}
		return
	}

	defer func() { _ = ws.Close() }()

	var message types.WSMessage

	ws.EnableWriteCompression(true)

	for {
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("WebSocket read error")
			break
		}

		log.Debug().
			Str("fileName", message.Filename).
			Str("percentage", fmt.Sprintf("%02.2f", message.Percentage)).
			Int("status", message.Status).
			Msg("recieved message")

		message.Percentage = 100

		err = ws.WriteJSON(&message)
		if err != nil {
			log.Error().
				Err(err).
				Msg("WebSocket write error")
			break
		}
	}
}
