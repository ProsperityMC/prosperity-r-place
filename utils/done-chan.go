package utils

import "sync"

type DoneChan struct {
	C chan struct{}
	m *sync.RWMutex
	e bool
}

func NewDoneChan() *DoneChan {
	return &DoneChan{
		C: make(chan struct{}, 0),
		m: &sync.RWMutex{},
	}
}

func (d *DoneChan) Close() {
	d.m.Lock()
	if !d.e {
		close(d.C)
	}
	d.m.Unlock()
}

func (d *DoneChan) Running() bool {
	d.m.RLock()
	defer d.m.RUnlock()
	return !d.e
}
