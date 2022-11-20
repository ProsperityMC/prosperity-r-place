package prosperity_r_place

import (
	"github.com/gorilla/websocket"
	"image"
	"prosperity-r-place/shapes"
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
		case "ping":
			_ = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
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
		case "circle":
			_ = conn.WriteMessage(websocket.TextMessage, []byte("todo"))
		case "square":
			if len(line) != 5 {
				break outer
			}
			fill, err := utils.ParseFill(line[1])
			if err != nil {
				break outer
			}

			colour, err := utils.ParseColor(line[2])
			if err != nil {
				break outer
			}

			topLeft, err := utils.ParseCoordinate(line[3])
			if err != nil {
				break outer
			}

			bottomRight, err := utils.ParseCoordinate(line[4])
			if err != nil {
				break outer
			}

			manager.placing <- shapes.PixelsInSquare(image.Rectangle{Min: topLeft, Max: bottomRight}, colour, fill)
			_ = conn.WriteMessage(websocket.TextMessage, []byte("done"))
		case "pentagon":
			_ = conn.WriteMessage(websocket.TextMessage, []byte("todo"))
		case "hexagon":
			_ = conn.WriteMessage(websocket.TextMessage, []byte("todo"))
		default:
			_ = conn.WriteMessage(websocket.TextMessage, []byte("error"))
			break outer
		}
	}
}
