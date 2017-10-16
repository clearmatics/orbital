package main

import (
	"testing"
)

func TestRingString(t *testing.T) {
	r := Ring{
		PubKeys: make([]PubKey, 1),
	}
	actual := r.String()
	expected := "X: <nil>, Y: <nil>\n"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
