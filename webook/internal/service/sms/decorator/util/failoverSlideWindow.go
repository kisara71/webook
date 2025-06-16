package util

import (
	"sync"
	"time"
)

type Result struct {
	ts      int64
	success bool
}

type FailOverSlideWindow struct {
	window        []Result
	failThreshold float64
	timeWindow    int64
	minCount      int
	mu            sync.Mutex
}

func NewFailOverSlideWindow(failThreshold float64, window int64, minCount int) *FailOverSlideWindow {
	return &FailOverSlideWindow{
		window:        make([]Result, 0, minCount*2),
		failThreshold: failThreshold,
		timeWindow:    window,
		minCount:      minCount,
	}
}
func (f *FailOverSlideWindow) Add(success bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.window = append(f.window, Result{
		ts:      time.Now().UnixMilli(),
		success: success,
	})
}

func (f *FailOverSlideWindow) ShouldFailOver() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now().UnixMilli()
	startTime := now - f.timeWindow

	valid := f.window[:0]
	var failed, total int

	for _, res := range f.window {
		if res.ts >= startTime {
			total++
			if !res.success {
				failed++
			}
			valid = append(valid, res)
		}

	}
	f.window = valid
	if total < f.minCount {
		return false
	}
	return float64(failed)/float64(total) > f.failThreshold
}
