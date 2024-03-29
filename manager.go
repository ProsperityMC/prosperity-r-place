package prosperity_r_place

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"image"
	"image/png"
	"io"
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
	cMap    map[string]clientInfo
	cLock   *sync.RWMutex
}

type clientInfo struct {
	conn *websocket.Conn
	info utils.DiscordInfo
	send chan []byte
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
		cMap:    make(map[string]clientInfo),
		cLock:   &sync.RWMutex{},
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
			m.broadcastImage()
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
		log.Println("[Manager::backgroundIO] Failed to encode PNG image:", err)
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
	_, err := m.file.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("[Manager::backgroundIO] Failed to seek to the start of the image:", err)
		return
	}
	err = m.file.Truncate(0)
	if err != nil {
		log.Println("[Manager::backgroundIO] Failed to truncate the image:", err)
		return
	}
	m.cacheS.RLock()
	a := m.cache
	m.cacheS.RUnlock()
	_, err = m.file.Write(a)
	if err != nil {
		log.Println("[Manager::backgroundIO] Failed to save image:", err)
	}
}

func (m *Manager) AddClient(uStr string, conn *websocket.Conn, info utils.DiscordInfo, sender chan []byte) {
	m.cLock.Lock()
	m.cMap[uStr] = clientInfo{
		conn: conn,
		info: info,
		send: sender,
	}
	m.cLock.Unlock()
	m.Broadcast([]byte(fmt.Sprintf("names %s=%s", uStr, info.Name)))
}

func (m *Manager) RemoveClient(uStr string) {
	m.cLock.Lock()
	delete(m.cMap, uStr)
	m.cLock.Unlock()
	m.Broadcast([]byte(fmt.Sprintf("quit %s", uStr)))
}

func (m *Manager) Broadcast(data []byte) {
	m.cLock.RLock()
	for _, v := range m.cMap {
		v.send <- data
	}
	m.cLock.RUnlock()
}

func (m *Manager) broadcastImage() {
	m.cacheS.RLock()
	b := m.cacheB
	m.cacheS.RUnlock()
	m.Broadcast(append([]byte("refresh "), b...))
}
