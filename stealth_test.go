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
    ssA := StealthPrivDerive(As, sharedSecret)
    ssAp := derivePublicKey(ssA)

    // stealth address on B side
    spB := StealthPubDerive(Ap, sharedSecret)
    ssB := StealthPrivDerive(Bs, sharedSecret)
    ssBp := derivePublicKey(ssB)

    if false == spA.Equals(&ssBp) {
        t.Fatal("Stealth address deriviation failure A->B")
    }

    if false == ssAp.Equals(&spB) {
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


// Verify that invalid secret keys cannot be used
// References:
//  - https://crypto.stackexchange.com/a/30272
//
func TestStealthInvalidSecret(t *testing.T) {
    _, _, Bp, _ := generatePairOfTestKeys(t)

    bigZero := new(big.Int).SetInt64(int64(0))
    bigOne := new(big.Int).SetInt64(int64(1))
    nPlusOne := new(big.Int).Add(group.N, bigOne)

    testBytes := []byte("test")

    invalidSecretKeys := []*big.Int{bigZero, group.N, nPlusOne}

    for _, secretKey := range invalidSecretKeys {
        if nil != StealthPrivDerive(secretKey, testBytes) {
            t.Log(secretKey, "accepted as secret key to StealthPrivDerive")
        }

        if derivePublicKey(secretKey).IsOnCurve() {
            t.Log(secretKey, "accepted as secret key to derivePublicKey")
        }

        if nil != deriveSharedSecret(secretKey, Bp) {
            t.Log(secretKey, "accepted as secret key to deriveSharedSecret")
        }

        if nil != NewStealthSession(secretKey, Bp, 0, 1) {
            t.Log(secretKey, "accepted as secret key to NewStealthSession")
        }
    }
}
