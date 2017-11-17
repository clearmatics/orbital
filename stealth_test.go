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