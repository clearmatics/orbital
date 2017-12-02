package main

import (
    "math/big"
    "crypto/sha256"
	"testing"
)


func TestCurvepointGenerate (t *testing.T) {
	Ap, As, err := generateKeyPair()
    if err != nil {
        t.Fatal(err)
    }
    if ! Ap.IsOnCurve() {
    	t.Fatal("Generated invalid public key")
    }
    if ! isValidSecretKey(As) {
    	t.Fatal("Generated invalid secret key")
    }
}


func TestHashToCurve (t *testing.T) {
    test_x, _ := new(big.Int).SetString("18149469767584732552991861025120904666601524803017597654373315627649680264678", 10)
    test_y, _ := new(big.Int).SetString("18593544354303197021588991433499968191850988132424885073381608163097237734820", 10)

    test_p := new(CurvePoint).SetFromXY(test_x, test_y)
    if ! test_p.IsOnCurve() {
        t.Fatal("Test vector not on curve")
    }

    h := sha256.Sum256([]byte("hello world"))
    p := NewCurvePointFromHash(h)

    x, y := p.GetXY()
    if x.Cmp(test_x) != 0 || y.Cmp(test_y) != 0 {
        t.Fatal("HashToCurve(sha256('1')), got", x, y, "expected", test_x, test_y)
    }
}


// This verifies that the Generator for the G1 curve is correct
// See: https://github.com/ethereum/go-ethereum/pull/15591
func TestCurvepointBaseMult (t *testing.T) {
    h := sha256.Sum256([]byte("1"))
    a := new(big.Int).SetBytes(h[:])

    // Verify sha256 hash of "1" is correct
    // sha256("1").hexdigest()
    a_test, _ := new(big.Int).SetString("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", 16)
    if a.Cmp(a_test) != 0 {
        t.Fatal("sha('1') test vector incorrect")
    }

    // Create 'public key' from private key of sha256("1")
    // Python equivalent: bn128.multiply(bn128.G1, int(sha256("1").hexdigest(), 16))
    b := CurvePoint{}.ScalarBaseMult(a)
    b_x, b_y := b.GetXY()

    // Verify test vector from solidity and py_ecc.bn128
    test_x, _ := new(big.Int).SetString("18402258484067100825836416533206638046709953333460439275068607944552700874793", 10)
    test_y, _ := new(big.Int).SetString("3216486158313018618592493241388793958480998389453172132732084762339402552220", 10)
    if b_x.Cmp(test_x) != 0 || b_y.Cmp(test_y) != 0 {
        t.Fatal("Test vector incorrect, ", b, "should be", test_x, test_y)
    }

    // Verify CurvePoint can be unserialized from the outputs it gives you
    c := new(CurvePoint).SetFromXY(b_x, b_y);
    if c == nil {
        t.Fatal("CurvePoint unserialize failed, presumably given invalid points", b, b_x, b_y, c)
    }
    if ! c.Equals(&b) {
        t.Fatal("Points not equal after serialize > unserialize", b, c)
    }
}