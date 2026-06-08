package client

import (
	"sync"
	"sync/atomic"
	"time"
)

const recentForwardWindow = 3 * time.Second

var (
	activeForwards  atomic.Int64
	totalForwards   atomic.Uint64
	lastForwardNano atomic.Int64

	forwardHookMu sync.RWMutex
	onBeginForward func()
	onEndForward   func()
)

// SetForwardHooks registers callbacks for the current embedded client session.
// Pass nil handlers to clear.
func SetForwardHooks(onBegin, onEnd func()) {
	forwardHookMu.Lock()
	onBeginForward = onBegin
	onEndForward = onEnd
	forwardHookMu.Unlock()
}

// ClearForwardHooks removes session callbacks.
func ClearForwardHooks() {
	SetForwardHooks(nil, nil)
}

// BeginForward marks the start of an active forwarding session.
func BeginForward() {
	forwardHookMu.RLock()
	begin := onBeginForward
	forwardHookMu.RUnlock()
	if begin != nil {
		begin()
		return
	}
	activeForwards.Add(1)
	totalForwards.Add(1)
	lastForwardNano.Store(time.Now().UnixNano())
}

// EndForward marks the end of an active forwarding session.
func EndForward() {
	forwardHookMu.RLock()
	end := onEndForward
	forwardHookMu.RUnlock()
	if end != nil {
		end()
		return
	}
	activeForwards.Add(-1)
}

// ActiveForwardCount returns the number of forwarding sessions in progress.
func ActiveForwardCount() int64 {
	n := activeForwards.Load()
	if n < 0 {
		return 0
	}
	return n
}

// TotalForwardCount returns how many forwarding sessions have started since process start.
func TotalForwardCount() uint64 {
	return totalForwards.Load()
}

// RecentlyForwarded reports whether a forwarding session started within the recent window.
func RecentlyForwarded() bool {
	last := lastForwardNano.Load()
	if last == 0 {
		return false
	}
	return time.Since(time.Unix(0, last)) <= recentForwardWindow
}

// ResetForwardCount clears the fallback global forwarding counters.
func ResetForwardCount() {
	activeForwards.Store(0)
	totalForwards.Store(0)
	lastForwardNano.Store(0)
}
