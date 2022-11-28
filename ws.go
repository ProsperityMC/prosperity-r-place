package prosperity_r_place

import (
	"github.com/gorilla/websocket"
	"prosperity-r-place/utils"
	"strings"
)

func HandleWebsocket(conn *websocket.Conn, manager *Manager) {
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
outer:
	for {
		messageType, r, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.BinaryMessage {
			// no binary messages
			break
		}

		line := strings.Split(string(r), " ")
		if len(line) < 1 {
			break outer
		}
		switch line[0] {
		case "check":
			if len(line) == 2 {
				manager.cacheS.RLock()
				a := manager.eTag
				if line[1] == a {
					manager.cacheS.RUnlock()
					_ = conn.WriteMessage(websocket.TextMessage, []byte("done"))
				} else {
					b := manager.cacheB
					manager.cacheS.RUnlock()
					_ = conn.WriteMessage(websocket.TextMessage, []byte("refresh "+a+" "+b))
				}
			} else {
				_ = conn.WriteMessage(websocket.TextMessage, []byte("error"))
				break outer
			}
		case "draw":
			if len(line) < 2 {
				break outer
			}
			colour, err := utils.ParseColor(line[1])
			if err != nil {
				break outer
			}

			pixelsRaw := line[2:]
			pixels := make([]utils.Pixel, len(pixelsRaw))

			for i, pixel := range pixelsRaw {
				p, err := utils.ParseCoordinate(pixel)
				if err != nil {
					break outer
				}
				pixels[i] = utils.Pixel{Point: p, Colour: colour}
			}
			manager.placing <- pixels
			_ = conn.WriteMessage(websocket.TextMessage, []byte("done"))
		default:
			_ = conn.WriteMessage(websocket.TextMessage, []byte("error"))
			break outer
		}
	}
}
