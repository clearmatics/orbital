// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import "math/big"

// A RingSignature is represented as a curve point and the signature data itself
type RingSignature struct {
	Tau    CurvePoint `json:"tau"`
	Ctlist []*big.Int `json:"ctlist"`
}
