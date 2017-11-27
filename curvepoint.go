// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"crypto/rand"
	"math/big"
	"github.com/clearmatics/bn256"
	"crypto/sha256"
)


// CurvePoint represents a point on an elliptic curve
type CurvePoint struct {
	z *bn256.G1
}


func (c CurvePoint) MarshalJSON() ([]byte, error) {
	x, y := c.GetXY()
	return json.Marshal(&struct {
		X *big.Int `json:"x"`
		Y *big.Int `json:"y"`
	}{
		X: x,
		Y: y,
	})
}


// Equals returns true if X and Y of both curve points are equal
func (c CurvePoint) Equals(d *CurvePoint) bool {
	return bytes.Compare(c.Marshal(), d.Marshal()) == 0;
}

func (c CurvePoint) Prime() *big.Int {
	return bn256.P
}

func (c CurvePoint) Order() *big.Int {
	return bn256.Order
}


// isBetween checks number is within range of (lower,upper)
// e.g. number > lower && number < upper
func isBetween (number *big.Int, lower *big.Int, upper *big.Int) bool {
	return false == (number.Cmp(lower) <= 0 || number.Cmp(upper) >= 0)
}


// randomPositiveBelow generates a uniformly random number between 1 and `below` 
func randomPositiveBelow (below *big.Int) *big.Int {
	for {
        number, err := rand.Int(rand.Reader, bn256.Order)
        if err != nil {
            return nil
        }

        // x > 0 && x < below
        if isBetween(number, bigZero, below) {
        	return number
        }
    }
}

// RandomN returns a uniformly random integer between 1 and P-1
func (c CurvePoint) RandomN() *big.Int {
	return randomPositiveBelow(c.Order())
}

// RandomP returns a uniformly random integer between 1 and P-1
func (c CurvePoint) RandomP() *big.Int {
	return randomPositiveBelow(c.Prime())
}

func (c CurvePoint) GetXY() (*big.Int, *big.Int) {
	if c.z != nil {
		// TODO: c.z.MakeAffine(nil) instead of marshal! Less byte buffer copying
		buffer := c.z.Marshal()
		x := new(big.Int).SetBytes(buffer[:32])
		y := new(big.Int).SetBytes(buffer[32:])
		return x, y
	}
	return nil, nil
}

func (c CurvePoint) SetFromXY (x *big.Int, y *big.Int) *CurvePoint {
	z, isOk := new(bn256.G1).SetFromXY(x, y)
	if isOk {
		c.z = z
		return &c
	}
	return nil
}

func (c CurvePoint) Marshal() []byte {
	return c.z.Marshal()
}

func (c CurvePoint) Unmarshal(m []byte) bool {
	_, ret := c.z.Unmarshal(m)
	return ret
}

// IsOnCurve returns true if point is on curve
func (c CurvePoint) IsOnCurve() bool {
	p := c.z.Point()
	p.MakeAffine(nil)
	return p.IsOnCurve()
}

func (c CurvePoint) String() string {
	return fmt.Sprintf("CurvePoint(%v)", c.z)
}

func NewCurvePointFromString(s []byte) *CurvePoint {
	return NewCurvePointFromHash(sha256.Sum256(s))
}

// NewCurvePointFromHash implements the 'try-and-increment' method of
// hashing into a curve which preserves random oracle proofs of security
//
// See: https://www.normalesup.org/~tibouchi/papers/bnhash-scis.pdf
//
func NewCurvePointFromHash(h [sha256.Size]byte) *CurvePoint {
	P := CurvePoint{}.Prime()
	N := CurvePoint{}.Order()

	x := new(big.Int).SetBytes(h[:])
	x.Mod(x, N)

	// TODO: limit number of iterations?
	for {
		xxx := new(big.Int).Mul(x, x)
		xxx.Mul(xxx, x)
		t := new(big.Int).Add(xxx, curveB)

		y := new(big.Int).ModSqrt(t, P)
		if y != nil {
			curveout := CurvePoint{}.SetFromXY(x, y)
			if curveout != nil {
				return curveout
			}
		}
		x.Add(x, bigOne)
	}
}

// ScalarBaseMult returns the product x where the result and base are the x coordinates of group points, base is the standard generator
func (c CurvePoint) ScalarBaseMult(x *big.Int) CurvePoint {
	return CurvePoint{new(bn256.G1).ScalarBaseMult(x)}
}

// ScalarMult returns the product c*x where the result and base are the x coordinates of group points 
func (c CurvePoint) ScalarMult(x *big.Int) CurvePoint {
	return CurvePoint{new(bn256.G1).ScalarMult(c.z, x)}
}

// Add performs an addition of two elliptic curve points
func (c CurvePoint) Add(y CurvePoint) CurvePoint {
	return CurvePoint{new(bn256.G1).Add(c.z, y.z)}
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
	x, errX := new(big.Int).SetString(pointX, 0)
	y, errY := new(big.Int).SetString(pointY, 0)
	if ! errX || ! errY {
		return nil;
	}

	c := CurvePoint{}
	c.SetFromXY(x, y)
	return &c
}
