package prosperity_r_place

import (
	"github.com/gorilla/websocket"
)

func HandleWebsocket(conn *websocket.Conn) {
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			return
		}
		conn.NextWriter()
	}
}
