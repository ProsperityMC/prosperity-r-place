package prosperity_r_place

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"prosperity-r-place/utils"
	"strings"
)

func HandleWebsocket(conn *websocket.Conn, manager *Manager, userInfo utils.DiscordInfo) {
	uStr := uuid.NewString()
	sender := make(chan []byte, 5)
	sendClose := make(chan struct{}, 0)
	manager.AddClient(uStr, conn, userInfo, sender)

	defer func(conn *websocket.Conn) {
		close(sendClose)
		manager.RemoveClient(uStr)
		_ = conn.Close()
	}(conn)

	// using a separate channel controlled sender to prevent concurrent writes
	go func() {
		for {
			select {
			case <-sendClose:
				return
			case a := <-sender:
				_ = conn.WriteMessage(websocket.TextMessage, a)
			}
		}
	}()

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
			sender <- []byte("pong")
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
			sender <- []byte("done")
		case "science":
			if len(line) < 2 {
				break outer
			}
			manager.Broadcast([]byte("science " + uStr + " " + strings.Join(line[1:], " ")))
		case "start":
			// send names of current users
			manager.cLock.RLock()
			a := make([]string, 0)
			for k, v := range manager.cMap {
				a = append(a, fmt.Sprintf("%s=%s", k, v.info.Name))
			}
			manager.cLock.RUnlock()
			sender <- []byte("names " + strings.Join(a, " "))

			// send current image
			manager.cacheS.RLock()
			sender <- []byte("refresh " + manager.cacheB)
			manager.cacheS.RUnlock()
		default:
			sender <- []byte("error")
			break outer
		}
	}
}
