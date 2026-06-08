package client

import "testing"

func TestActiveForwardCount(t *testing.T) {
	ResetForwardCount()
	ClearForwardHooks()
	if got := ActiveForwardCount(); got != 0 {
		t.Fatalf("initial count = %d, want 0", got)
	}
	BeginForward()
	BeginForward()
	if got := ActiveForwardCount(); got != 2 {
		t.Fatalf("count after begin = %d, want 2", got)
	}
	if got := TotalForwardCount(); got != 2 {
		t.Fatalf("total after begin = %d, want 2", got)
	}
	EndForward()
	if got := ActiveForwardCount(); got != 1 {
		t.Fatalf("count after end = %d, want 1", got)
	}
	ResetForwardCount()
	if got := ActiveForwardCount(); got != 0 {
		t.Fatalf("count after reset = %d, want 0", got)
	}
}

func TestForwardHooksOverrideGlobals(t *testing.T) {
	ResetForwardCount()
	ClearForwardHooks()

	var hookActive int
	SetForwardHooks(func() { hookActive++ }, func() { hookActive-- })
	BeginForward()
	BeginForward()
	if hookActive != 2 {
		t.Fatalf("hook active = %d, want 2", hookActive)
	}
	if got := ActiveForwardCount(); got != 0 {
		t.Fatalf("global active should stay 0 with hooks, got %d", got)
	}
	EndForward()
	if hookActive != 1 {
		t.Fatalf("hook active after end = %d, want 1", hookActive)
	}
	ClearForwardHooks()
}
