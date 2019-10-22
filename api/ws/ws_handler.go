package ws

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
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

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Warn().
				Err(err).
				Msg("WebSocket read error")
			break
		}

		log.Debug().
			Str("received via websocket: %s", string(message))

		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Error().
				Err(err).
				Msg("WebSocket write error")
			break
		}
	}
}
