// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"testing"
)

func TestRingString(t *testing.T) {
	r := Ring{
		PubKeys: make([]PubKey, 1),
	}
	actual := r.String()
	expected := "{ \"x\": \"<nil>\", \"y\": \"<nil>\" }\n"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
