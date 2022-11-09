package prosperity_r_place

import (
	"image"
	"prosperity-r-place/utils"
	"time"
)

type Manager struct {
	width   int
	height  int
	img     *image.RGBA
	placing chan utils.Pixel
	save    time.Timer
	done    *utils.DoneChan
}

func NewManager(width, height int) *Manager {
	return &Manager{
		width:   width,
		height:  height,
		img:     image.NewRGBA(image.Rect(0, 0, width, height)),
		placing: make(chan utils.Pixel, 32),
		save:    time.NewTimer(),
		done:    utils.NewDoneChan(),
	}
}

// backgroundIO handles all IO operations in a single thread
func (m *Manager) backgroundIO() {
outer:
	for {
		select {
		case <-m.done.C:
			break outer
		case pixel := <-m.placing:
			m.img.SetRGBA(pixel.X, pixel.Y, pixel.Colour)
		}
	}
}
