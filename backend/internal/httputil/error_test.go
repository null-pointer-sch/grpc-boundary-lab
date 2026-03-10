package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("test error")
	code := http.StatusBadRequest

	WriteError(w, code, err)

	if w.Code != code {
		t.Errorf("expected status %d, got %d", code, w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Error != http.StatusText(code) {
		t.Errorf("expected error %s, got %s", http.StatusText(code), resp.Error)
	}

	if resp.Message != "test error" {
		t.Errorf("expected message 'test error', got '%s'", resp.Message)
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}

	WriteJSON(w, data)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var got map[string]string
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", got["foo"])
	}
}
