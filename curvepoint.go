// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"bytes"
	"fmt"
	"errors"
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


func (c *CurvePoint) MarshalJSON() ([]byte, error) {
	x, y := c.GetXY()
	return json.Marshal(&struct {
		X *hexBig `json:"x"`
		Y *hexBig `json:"y"`
	}{
		X: (*hexBig)(x),
		Y: (*hexBig)(y),
	})
}


func (c *CurvePoint) UnmarshalJSON(data []byte) error {
	var aux struct {
		X *hexBig `json:"x"`
		Y *hexBig `json:"y"`
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux.X == nil || aux.Y == nil {
		return errors.New("Invalid Point, no X or Y specified")
	}

	if c.SetFromXY((*big.Int)(aux.X), (*big.Int)(aux.Y)) == nil {
		return errors.New("Failed to deserialize CurvePoint")
	}

	return nil
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
	// Each value is a 256-bit number.	
	const numBytes = 256 / 8
	if c.z != nil {
		m := c.z.Marshal()
		x := new(big.Int).SetBytes(m[0*numBytes : 1*numBytes])
        y := new(big.Int).SetBytes(m[1*numBytes : 2*numBytes])
		return x, y
	}
	return nil, nil
}

func (c *CurvePoint) SetFromXY (x *big.Int, y *big.Int) *CurvePoint {
	const numBytes = 256 / 8

	// XXX: there's no equivalent to SetCurvePoints, other than Unmarshal
	xBytes := new(big.Int).Mod(x, bn256.P).Bytes()
	yBytes := new(big.Int).Mod(y, bn256.P).Bytes()

	m := make([]byte, numBytes*2)
	copy(m[1*numBytes-len(xBytes):], xBytes)
	copy(m[2*numBytes-len(yBytes):], yBytes)

	z, isOk := new(bn256.G1).Unmarshal(m)
	if isOk {
		c.z = z
		return c
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
	return c.z.IsOnCurve()
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
func NewCurvePointFromHash(h [sha256.Size]byte) *CurvePoint {
	P := CurvePoint{}.Prime()
	N := CurvePoint{}.Order()

	// (p+1) / 1
	A, _ := new(big.Int).SetString("c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52", 16)

	x := new(big.Int).SetBytes(h[:])
	x.Mod(x, N)

	// TODO: limit number of iterations?
	// y² = x³ + B
	for {
		xx := new(big.Int).Mul(x, x)		// x²
		xx.Mod(xx, P)

		xxx := xx.Mul(xx, x)				// x³
		xxx.Mod(xxx, P)

		beta := new(big.Int).Add(xxx, curveB)	// x³ + B
		beta.Mod(beta, P)						

		//y := new(big.Int).ModSqrt(t, P)		// y = √(x³+B)
		y := new(big.Int).Exp(beta, A, P)

		if y != nil {
			// Then verify (√(x³+B)%P)² == (x³+B)%P
			z := new(big.Int).Mul(y, y)
			z.Mod(z, P)
			if z.Cmp(beta) == 0 {
				curveout := new(CurvePoint).SetFromXY(x, y)
				if curveout != nil {
					return curveout
				}
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
