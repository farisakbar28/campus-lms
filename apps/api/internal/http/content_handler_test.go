package http

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOptionalTimeDistinguishesOmittedAndNull(t *testing.T) {
	var omitted struct {
		AvailableFrom optionalTime `json:"available_from"`
	}
	if err := json.Unmarshal([]byte(`{}`), &omitted); err != nil {
		t.Fatalf("decode omitted time: %v", err)
	}
	if omitted.AvailableFrom.set {
		t.Fatal("omitted time marked as set")
	}

	var cleared struct {
		AvailableFrom optionalTime `json:"available_from"`
	}
	if err := json.Unmarshal([]byte(`{"available_from":null}`), &cleared); err != nil {
		t.Fatalf("decode null time: %v", err)
	}
	if !cleared.AvailableFrom.set || cleared.AvailableFrom.value != nil {
		t.Fatalf("null time = %#v, want set with nil value", cleared.AvailableFrom)
	}

	var assigned struct {
		AvailableFrom optionalTime `json:"available_from"`
	}
	if err := json.Unmarshal([]byte(`{"available_from":"2030-01-01T00:00:00Z"}`), &assigned); err != nil {
		t.Fatalf("decode assigned time: %v", err)
	}
	want := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	if !assigned.AvailableFrom.set || assigned.AvailableFrom.value == nil || !assigned.AvailableFrom.value.Equal(want) {
		t.Fatalf("assigned time = %#v, want %s", assigned.AvailableFrom, want)
	}
}
