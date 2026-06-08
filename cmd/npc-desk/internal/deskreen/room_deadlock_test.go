package deskreen

import (
	"testing"
	"time"
)

// SetPending must not call onChange while holding the write lock (emitState reads PublicState).
func TestSetPendingDoesNotDeadlockOnChange(t *testing.T) {
	r := NewRoom("3131")
	done := make(chan struct{})
	r.SetOnChange(func() {
		_ = r.PublicState()
		close(done)
	})
	if !r.SetPending(DeviceInfo{IP: "192.168.1.10", Browser: "Chrome"}) {
		t.Fatal("SetPending failed")
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("onChange blocked (likely deadlock)")
	}
	if !r.HasPendingViewer() {
		t.Fatal("expected pending viewer")
	}
}
