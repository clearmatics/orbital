// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"bytes"
	"fmt"
	"math/big"
	"golang.org/x/crypto/bn256"
)

// CurvePoint represents a point on an elliptic curve
type CurvePoint struct {
	Z *bn256.G1 `json:""`
	/*
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
	*/
}

// Equals returns true if X and Y of both curve points are equal
func (c CurvePoint) Equals(d *CurvePoint) bool {
	// 0 == c.X.Cmp(d.X) && 0 == c.Y.Cmp(d.Y)
	return bytes.Compare(c.Marshal(), d.Marshal()) == 0;
}

func (c CurvePoint) InitFromXY (x *big.Int, y *big.Int) (*CurvePoint, bool) {
	// XXX: Pack into format that Unmarshal accepts...
	// TODO: split bn256.G1.Unmarshal into Unmarshal + SetFromXY
	const numBytes = 256/8
	packed := make([]byte, numBytes*2)
	xBytes := x.Bytes()
	yBytes := y.Bytes()
	copy(packed[1*numBytes-len(xBytes):], xBytes)
	copy(packed[2*numBytes-len(yBytes):], yBytes)
	z, isOk := new(bn256.G1).Unmarshal(packed)
	if isOk {
		c.Z = z
	}
	return &c, isOk
}

func (c CurvePoint) Marshal() []byte {
	return c.Z.Marshal()
}

func (c CurvePoint) Unmarshal(m []byte) bool {
	_, ret := c.Z.Unmarshal(m)
	return ret
}

func (c CurvePoint) InitFromSecret (x *big.Int) {
	c.Z = new(bn256.G1).ScalarBaseMult(x)
}

// IsOnCurve returns true if point is on curve
func (c CurvePoint) IsOnCurve() bool {
	return false
	//return group.IsOnCurve(c.X, c.Y)
}

func (c CurvePoint) String() string {
	return fmt.Sprintf("CurvePoint(%v)", c.Z)
}

// ScalarBaseMult returns the product x where the result and base are the x coordinates of group points, base is the standard generator
func (c CurvePoint) ScalarBaseMult(x *big.Int) CurvePoint {
	return CurvePoint{new(bn256.G1).ScalarBaseMult(x)}
	/*
	px, py := group.ScalarBaseMult(x.Bytes())
	return CurvePoint{px, py}
	*/
}

// ScalarMult returns the product c*x where the result and base are the x coordinates of group points 
func (c CurvePoint) ScalarMult(x *big.Int) CurvePoint {
	return CurvePoint{new(bn256.G1).ScalarMult(c.Z, x)}
	/*
	px, py := group.ScalarMult(c.X, c.Y, x.Bytes())
	return CurvePoint{px, py}
	*/
}

// Add performs an addition of two elliptic curve points
func (c CurvePoint) Add(y CurvePoint) CurvePoint {
	return CurvePoint{new(bn256.G1).Add(c.Z, y.Z)}
	/*
	px, py := group.Add(c.X, c.Y, y.X, y.Y)
	return CurvePoint{px, py}
	*/
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

// ParseCurvePoint parses string representations of X and Y points
// these can be hex or base10 encoded
func ParseCurvePoint( pointX string, pointY string ) *CurvePoint {
	/*
	X, errX := new(big.Int).SetString(pointX, 0)
	Y, errY := new(big.Int).SetString(pointY, 0)
	if ! errX || ! errY {
		return nil;
	}

	return &CurvePoint{X, Y}
	*/

	// TODO: G1.Marshal() returns 512bit in bytes
	return nil
}
