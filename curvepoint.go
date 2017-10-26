// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"fmt"
	"math/big"
)

type CurvePoint struct {
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
}

func (c CurvePoint) String() string {
	return fmt.Sprintf("X: %s, Y: %s", c.X, c.Y)
}

func (c CurvePoint) ScalarBaseMult(x *big.Int) CurvePoint {
	px, py := Group.ScalarBaseMult(x.Bytes())
	return CurvePoint{px, py}
}

func (c CurvePoint) ScalarMult(x *big.Int) CurvePoint {
	px, py := Group.ScalarMult(c.X, c.Y, x.Bytes())
	return CurvePoint{px, py}
}

func (c CurvePoint) Add(y CurvePoint) CurvePoint {
	px, py := Group.Add(c.X, c.Y, y.X, y.Y)
	return CurvePoint{px, py}
}

func (c CurvePoint) ParameterPointAdd(tj *big.Int, cj *big.Int) CurvePoint {
	a := CurvePoint{}.ScalarBaseMult(tj)
	pk := c.ScalarMult(cj)

	return a.Add(pk)
}

func (c CurvePoint) HashPointAdd(hashSP CurvePoint, tj *big.Int, cj *big.Int) CurvePoint {
	b := c.ScalarMult(tj)
	bj := hashSP.ScalarMult(cj)

	return b.Add(bj)
}
