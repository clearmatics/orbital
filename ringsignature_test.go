package main

import (
	"testing"
)

func TestRingSignatureString(t *testing.T) {
	r := RingSignature{}
	actual := r.String()
	expected := "tau: X: <nil>, Y: <nil>\nctlist: [\n]\n"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
