package capture

import "testing"

func TestUseWindowCapture(t *testing.T) {
	if !UseWindowCapture(Source{ID: "12345678", Type: "", Name: "x"}, "") {
		t.Fatal("large id should be treated as HWND")
	}
	if UseWindowCapture(Source{ID: "0", Type: "", Name: "x"}, "") {
		t.Fatal("0 should be display index")
	}
	if !UseWindowCapture(Source{ID: "0", Type: "app"}, "app") {
		t.Fatal("explicit app mode")
	}
	if UseWindowCapture(Source{ID: "12345678", Type: "", Name: "x"}, "screen") {
		t.Fatal("explicit screen mode")
	}
}

func TestIsDisplayIndexID(t *testing.T) {
	if !IsDisplayIndexID("0") {
		t.Fatal("0 is display")
	}
	if IsDisplayIndexID("999999") {
		t.Fatal("large number is not display index")
	}
}
