// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
    "testing"
    "bytes"
)


func TestPubDerive(t *testing.T) {
    // Generate key pairs for either side
    Ap, As, err := generateKeyPair()
    if err != nil {
        t.Fatal(err)
    }
    Bp, Bs, err := generateKeyPair()
    if err != nil {
        t.Fatal(err)
    }

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
