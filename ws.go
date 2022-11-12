package prosperity_r_place

import (
	"bufio"
	"github.com/gorilla/websocket"
	"image/color"
	"strconv"
)

func HandleWebsocket(conn *websocket.Conn) {
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

			}
		}

		// read X coordinate string
		xStr, err := rBuf.ReadString(',')
		if err != nil {
			// bad packet
			break
		}

		// read Y coordinate string (
		yStr, err := rBuf.ReadString(',')
		if err != nil {
			return
		}

		// read colour string (line ends with \n)
		colourStr, err := rBuf.ReadString('\n')
		if err != nil {
			return
		}
	}
}
