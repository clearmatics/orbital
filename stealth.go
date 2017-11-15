// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main;

import (
    "crypto/sha256"
    "crypto/rand"
    "math/big"
)


// generateKeyPair generates a random secret key, then derives the
// public key from it
//
func generateKeyPair () (*CurvePoint, *big.Int, error) {    
    q := group.P
    priv, err := rand.Int(rand.Reader, q)
    if err != nil {
        return nil, nil, err
    }
    pub := derivePublicKey(priv)
    return &pub, priv, nil
}


/*
According to the paper (IACR 2017/881), the relationship between stealth
addresses is bidirectional:

    spkA = mpkA + g^H(secret||nonce)
    spkB = mpkB + g^H(secret||nonce)

Which means

    spkA / spkB = mpkA / mpkB

IMO, the notation in the paper is not as clear as it could be, but it's
easily translated and has proven to be correct in practice; the notation
used here has been adjusted accordingly.

Anyway, this means that A can know what B's stealth public key will be if
they both agree on a shared secret, and visa versa.
*/


// StealthPubDerive derives another parties Stealth Public Key (ssp) from
// their Master Public Key and an arbitrary shared secret.
//
// From IACR 2017/881 (2.1):
//
//   spk ← mpk + g^H(secret)
//
// Parameters:
//
//   mpk = their Public Key, as CurvePoint
//   secret = arbitrary number known by both parties
//
func StealthPubDerive(mpk *CurvePoint, secret *big.Int) CurvePoint {
    // X ← H(secret||nonce)
    _hashout := sha256.Sum256(secret.Bytes())
    X := new(big.Int).SetBytes(_hashout[:])

    // Y ← g^X
    Y := derivePublicKey(X)

    // spk ← mpk + Y
    spk := mpk.Add(Y)

    return spk
}


// StealthPrivDerive derives a Stealth Secret Key (ssk) from your
// Master Secret Key (msk), using an arbitrary shared secret.
//
// From IACR 2017/881 (2.1):
//
//   ssk ← msk + H(secret)
//
// Parameters:
// 
//   msk = Your secret key
//   secret = arbitrary number known by both parties
func StealthPrivDerive(msk *big.Int, secret *big.Int) *big.Int {
    // X ← H(secret)
    _hashout := sha256.Sum256(secret.Bytes())
    X := new(big.Int).SetBytes(_hashout[:])

    // ssk ← msk + X
    ssk := new(big.Int).Add(msk, X)

    return ssk
}


// derivePublicKey derives from SecretKey using ScalarBaseMult:
//
//    Pa,Pb ← g^S
//
func derivePublicKey (privateKey *big.Int) CurvePoint {
    return CurvePoint{}.ScalarBaseMult(privateKey)
}


// deriveSharedSecret between two key pairs, aka ECDH, with ScalarMult:
//
//    (Ax,_) ← (Bpx,Bpy) · As
//    (Bx,_) ← (Apx,Apy) · Bs
//    Ax ≡ Bx
//
// Where As and Bs are secret keys, (Bpx,Bpy) and (Apx,Apy) are the public
// keys of A and B. (Ax,_) and (Bx,_) are points, and both Ax and Ay are equal.
// The second points of the result are discarded according to RFC5903 (Section 9).
//
func deriveSharedSecret (myPriv *big.Int, theirPub *CurvePoint) *big.Int {
    // See: RFC5903 (Section 9)
    return theirPub.ScalarMult(myPriv).X
}
