package prosperity_r_place

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"prosperity-r-place/utils"
	"sync"
	"time"
)

const saveInterval = time.Second * 30

type Manager struct {
	file    *os.File
	width   int
	height  int
	img     *image.RGBA
	placing chan utils.Pixel
	save    *time.Timer
	done    *utils.DoneChan
	wg      *sync.WaitGroup
}

func NewManager(name string, width, height int) (*Manager, error) {
	create, err := os.Create(fmt.Sprintf(".data/images/%s.png", name))
	if err != nil {
		return nil, fmt.Errorf("failed to open '.data/images/%s.png' for writing: %w", err)
	}
	m := &Manager{
		file:    create,
		width:   width,
		height:  height,
		img:     image.NewRGBA(image.Rect(0, 0, width, height)),
		placing: make(chan utils.Pixel, 32),
		save:    time.NewTimer(saveInterval),
		done:    utils.NewDoneChan(),
		wg:      &sync.WaitGroup{},
	}
	// increment the wait group and start backgroundIO goroutine
	m.wg.Add(1)
	go m.backgroundIO()
	return m, nil
}

// Close handles cleaning up the manager data
func (m *Manager) Close() {
	// trigger done
	m.done.Close()
	// wait for backgroundIO to finish
	m.wg.Wait()

	// clear up resources
	m.save.Stop()
	close(m.placing)
	_ = m.file.Close()
}

// backgroundIO handles all IO operations in a single thread
func (m *Manager) backgroundIO() {
	defer m.wg.Done()
	lastSave := false
outer:
	for {
		select {
		case pixel := <-m.placing:
			// set the pixel
			m.img.SetRGBA(pixel.X, pixel.Y, pixel.Colour)
			// if the last operation was not a save then restart the save timer
			if !lastSave {
				m.save.Reset(saveInterval)
			}
			lastSave = false
		case <-m.save.C:
			// if the last operation was also a save then stop the save timer
			if lastSave {
				m.save.Stop()
				continue outer
			}
			// save image and set the lastSave flag
			m.saveImage()
			lastSave = true
		case <-m.done.C:
			// save image before closing
			m.saveImage()
			break outer
		}
	}
}

func (m *Manager) saveImage() {
	// png encode the image to a file and log errors
	err := png.Encode(m.file, m.img)
	if err != nil {
		log.Println("[Manager::backgroundIO] Failed to save PNG image:", err)
		return
	}
}
