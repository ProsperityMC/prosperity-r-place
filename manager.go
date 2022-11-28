package prosperity_r_place

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"prosperity-r-place/utils"
	"sync"
	"time"
)

const saveInterval = time.Second * 10

type Manager struct {
	Name    string `json:"name"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	file    *os.File
	img     *image.RGBA
	placing chan []utils.Pixel
	save    *time.Timer
	done    *utils.DoneChan
	wg      *sync.WaitGroup
	cacheS  *sync.RWMutex
	cacheB  string
	cache   []byte
	eTag    string
}

func NewManager(name string, width, height int) (*Manager, error) {
	fmt.Printf("[Manager] %s\n", name)
	fPath := fmt.Sprintf(".data/images/%s.png", name)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	f, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open '.data/images/%s.png' for saving: %w", name, err)
	} else {
		im, err := png.Decode(f)
		if err == nil {
			img = utils.ImageToRGBA(im)
		}
	}
	m := &Manager{
		Name:    name,
		Width:   width,
		Height:  height,
		file:    f,
		img:     img,
		placing: make(chan []utils.Pixel, 32),
		save:    time.NewTimer(saveInterval),
		done:    utils.NewDoneChan(),
		wg:      &sync.WaitGroup{},
		cacheS:  &sync.RWMutex{},
		cacheB:  "",
		cache:   nil,
		eTag:    "",
	}
	// generate the first cache of the image
	m.encodeImage()
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

// Image returns the current cached image
func (m *Manager) Image() ([]byte, string) {
	m.cacheS.RLock()
	defer m.cacheS.RUnlock()
	return m.cache, m.eTag
}

// backgroundIO handles all IO operations in a single thread
func (m *Manager) backgroundIO() {
	defer m.wg.Done()
	lastSave := false
outer:
	for {
		select {
		case pixels := <-m.placing:
			// set the pixels
			for _, pixel := range pixels {
				m.img.SetRGBA(pixel.Point.X, pixel.Point.Y, pixel.Colour)
			}
			m.encodeImage()
			// if the last operation was not a save then restart the save timer
			if lastSave {
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

func (m *Manager) encodeImage() {
	// png encode the image to a file and log errors
	buf := new(bytes.Buffer)
	err := png.Encode(buf, m.img)
	if err != nil {
		log.Println("[Manager::backgroundIO] Failed to save PNG image:", err)
		return
	}
	sum := sha1.Sum(buf.Bytes())
	hex1 := hex.EncodeToString(sum[:])
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	m.cacheS.Lock()
	m.cache = buf.Bytes()
	m.cacheB = b64
	m.eTag = hex1
	m.cacheS.Unlock()
}

func (m *Manager) saveImage() {
	m.cacheS.RLock()
	a := m.cache
	m.cacheS.RUnlock()
	_, _ = m.file.Write(a)
}
