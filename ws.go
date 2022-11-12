package prosperity_r_place

import (
	"bufio"
	"github.com/gorilla/websocket"
	"image/color"
	"prosperity-r-place/utils"
	"strconv"
)

func HandleWebsocket(conn *websocket.Conn, manager *Manager) {
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
outer:
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			return
		}
		if messageType == websocket.BinaryMessage {
			// no binary messages
			break
		}

		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanWords)

		var x, y int64
		var colour color.RGBA
		i := 0
		for scanner.Scan() {
			switch i {
			case 0:
				x, err = strconv.ParseInt(scanner.Text(), 10, 64)
				if err != nil {
					break outer
				}
			case 1:
				y, err = strconv.ParseInt(scanner.Text(), 10, 64)
				if err != nil {
					break outer
				}
			case 2:
				text := scanner.Text()
				if text[0] == '#' {
					colourParse, err := strconv.ParseInt(text[1:], 16, 64)
					if err != nil {
						break outer
					}
					b1 := colourParse & 0xff
					g1 := (colourParse >> 8) & 0xff
					r1 := (colourParse >> 16) & 0xff
					colour = color.RGBA{R: uint8(r1), G: uint8(g1), B: uint8(b1)}
				}
			}
		}

		manager.placing <- utils.Pixel{X: x, Y: y, Colour: colour}
	}
}
