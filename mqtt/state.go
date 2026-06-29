package mqtt

import (
	"sync"
	"time"
)

type deviceState struct {
	mu          sync.Mutex
	timer       *time.Timer
	publishFunc func(state string)
}

func newDeviceState(publishFunc func(state string)) *deviceState {
	return &deviceState{publishFunc: publishFunc}
}

func (ds *deviceState) OnOpen(travelSecs int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.cancelTimer()
	ds.publishFunc("opening")
	ds.timer = time.AfterFunc(time.Duration(travelSecs)*time.Second, func() {
		ds.publishFunc("open")
	})
}

func (ds *deviceState) OnClose(travelSecs int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.cancelTimer()
	ds.publishFunc("closing")
	ds.timer = time.AfterFunc(time.Duration(travelSecs)*time.Second, func() {
		ds.publishFunc("closed")
	})
}

func (ds *deviceState) OnStop() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.cancelTimer()
	ds.publishFunc("stopped")
}

// cancelTimer stops the current timer. Must be called with ds.mu held.
func (ds *deviceState) cancelTimer() {
	if ds.timer != nil {
		ds.timer.Stop()
		ds.timer = nil
	}
}
