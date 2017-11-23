// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
    "testing"
    "bytes"
    "math/big"
)


func generatePairOfTestKeys(t *testing.T) (*CurvePoint, *big.Int, *CurvePoint, *big.Int) {
    Ap, As, err := generateKeyPair()
    if err != nil {
        t.Fatal(err)
    }
    Bp, Bs, err := generateKeyPair()
    if err != nil {
        t.Fatal(err)
    }
    return Ap, As, Bp, Bs
}


// TestSharedSecret verifies that a shared secret can be derived from two curve points
func TestSharedSecret (t *testing.T) {
    Ap, As, Bp, Bs := generatePairOfTestKeys(t)
    if Ap == nil || As == nil {
        t.Fatal("failed to generate pair of test keys")
    }

    sharedSecret := deriveSharedSecret(As, Bp)
    if sharedSecret == nil {
        t.Fatal("nil shared secret")
    }
    sharedSecretCheck := deriveSharedSecret(Bs, Ap)
    if 0 != bytes.Compare(sharedSecret, sharedSecretCheck) {
        t.Fatal("Shared secret incorrect")
    }
}


func TestStealthAddressPrimitives(t *testing.T) {
    Ap, As, Bp, Bs := generatePairOfTestKeys(t)

    // Using ECDH, derive shared secret between two key pairs
    sharedSecret := deriveSharedSecret(As, Bp)
    sharedSecretCheck := deriveSharedSecret(Bs, Ap)
    if 0 != bytes.Compare(sharedSecret, sharedSecretCheck) {
        t.Fatal("Shared secret incorrect")
    }
    
    // stealth address on A side
    spA := StealthPubDerive(Bp, sharedSecret)
    if spA == nil {
        t.Fatal("Failed to derive stealth public key for B from shared secret")
    }
    
    ssA := StealthPrivDerive(As, sharedSecret)
    if ssA == nil {
        t.Fatal("Failed to derive stealth private key for A from shared secret")
    }

    ssAp := derivePublicKey(ssA)

    // stealth address on B side
    spB := StealthPubDerive(Ap, sharedSecret)
    ssB := StealthPrivDerive(Bs, sharedSecret)
    ssBp := derivePublicKey(ssB)

    if false == spA.Equals(&ssBp) {
        t.Fatal("Stealth address deriviation failure A->B")
    }

    if false == ssAp.Equals(spB) {
        t.Fatal("Stealth address deriviation failure B->A")
    }
}


func TestStealthAddressSession(t *testing.T) {
    Ap, As, Bp, Bs := generatePairOfTestKeys(t)

    sessA := NewStealthSession(As, Bp, 0, 2)
    sessB := NewStealthSession(Bs, Ap, 0, 2)

    if ! sessA.TheirPublic.Equals(Bp) {
        t.Fatal("Public Key Mismatch, A.TheirP != Bp")
    }
    if ! sessB.TheirPublic.Equals(Ap) {
        t.Fatal("Public Key Mismatch, B.TheirP != Ap")
    }

    // Verify derived stealth addresses match on either side
    if ! sessA.MyAddresses[0].Public.Equals(&sessB.TheirAddresses[0].Public) {
        t.Fatal("Public Key Mismatch, A.MyA[0].P != B.TheirA[0].P")
    }
    if ! sessA.MyAddresses[1].Public.Equals(&sessB.TheirAddresses[1].Public) {
        t.Fatal("Public Key Mismatch, A.MyA[1].P != B.TheirA[1].P")
    }
}



var testBytes = []byte("test")

// Verify that invalid secret keys cannot be used
// References:
//  - https://crypto.stackexchange.com/a/30272
//
func TestStealthInvalidSecret(t *testing.T) {
    _, _, Bp, _ := generatePairOfTestKeys(t)

    var nPlusOne = new(big.Int).Add(CurvePoint{}.Order(), bigOne)
    var invalidSecretKeys = []*big.Int{bigZero, CurvePoint{}.Order(), nPlusOne}

    for _, secretKey := range invalidSecretKeys {
        if nil != StealthPrivDerive(secretKey, testBytes) {
            t.Fatal(secretKey, "accepted as secret key to StealthPrivDerive")
        }

        if nil != NewStealthSession(secretKey, Bp, 0, 1) {
            t.Fatal(secretKey, "accepted as secret key to NewStealthSession")
        }
    }
}


func TestStealthInvalidPublic(t *testing.T) {
    _, As, _, _ := generatePairOfTestKeys(t)
    var invalidSecretKeys = []*big.Int{bigZero, CurvePoint{}.Order()}

    for _, secretKey := range invalidSecretKeys {
        publicKey := derivePublicKey(secretKey)
        if StealthPubDerive(&publicKey, testBytes) != nil {
            t.Log(publicKey, "(from ", secretKey, ") accepted as public key to StealthPubDerive")
        }
        if NewStealthSession(As, &publicKey, 0, 1) != nil {
            t.Log(publicKey, "(from ", secretKey, ") accepted as public key to NewStealthSession")
        }

        // Deliberately create an invalid curve point from a valid one
        /*
        alteredPublicKey := CurvePoint{}.InitFromXY(new(big.Int).Add(Bp.X, bigOne), new(big.Int).Add(Bp.Y, bigOne))
        if StealthPubDerive(&alteredPublicKey, testBytes) != nil {
            t.Fatal(Bp, " + (1,-1) accepted as public key to StealthPubDerive")
        }
        if NewStealthSession(As, &alteredPublicKey, 0, 1) != nil {
            t.Fatal(Bp, " + (1,-1) accepted as public key to NewStealthSession")
        }
        */
    }
}
