// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"fmt"
	"math/big"
)

// CurvePoint represents a point on an elliptic curve
type CurvePoint struct {
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
}

// Equals returns true if X and Y of both curve points are equal
func (c CurvePoint) Equals(d *CurvePoint) bool {
	return 0 == c.X.Cmp(d.X) && 0 == c.Y.Cmp(d.Y)
}

func (c CurvePoint) String() string {
	return fmt.Sprintf("X: %s, Y: %s", c.X, c.Y)
}

// ScalarBaseMult returns the product x where the result and base are the x coordinates of group points, base is the standard generator
func (c CurvePoint) ScalarBaseMult(x *big.Int) CurvePoint {
	px, py := group.ScalarBaseMult(x.Bytes())
	return CurvePoint{px, py}
}

// ScalarMult returns the product c*x where the result and base are the x coordinates of group points 
func (c CurvePoint) ScalarMult(x *big.Int) CurvePoint {
	px, py := group.ScalarMult(c.X, c.Y, x.Bytes())
	return CurvePoint{px, py}
}

// Add performs an addition of two elliptic curve points
func (c CurvePoint) Add(y CurvePoint) CurvePoint {
	px, py := group.Add(c.X, c.Y, y.X, y.Y)
	return CurvePoint{px, py}
}

// ParameterPointAdd returns the addition of c scaled by cj and tj as a curve point
func (c CurvePoint) ParameterPointAdd(tj *big.Int, cj *big.Int) CurvePoint {
	a := CurvePoint{}.ScalarBaseMult(tj)
	pk := c.ScalarMult(cj)

	return a.Add(pk)
}

// HashPointAdd returns the addition of hashSP scaled by cj and c scaled by tj
func (c CurvePoint) HashPointAdd(hashSP CurvePoint, tj *big.Int, cj *big.Int) CurvePoint {
	b := c.ScalarMult(tj)
	bj := hashSP.ScalarMult(cj)

	return b.Add(bj)
}
