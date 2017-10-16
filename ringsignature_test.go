// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"math/big"
	"testing"
)

func TestRingSignatureString(t *testing.T) {
	r := RingSignature{
		Ctlist: make([]*big.Int, 1),
	}
	actual := r.String()
	expected := "tau: X: <nil>, Y: <nil>\nctlist: [\n\t<nil>\n]\n"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
