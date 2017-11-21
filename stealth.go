// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main;

import (
    "crypto/sha256"
    "crypto/rand"
    "math/big"
)


type StealthAddress struct {
    Public CurvePoint `json:"public"`
    Nonce *big.Int `json:"nonce"`
}


type PrivateStealthAddress struct {
    Public CurvePoint `json:"public"`
    Nonce *big.Int `json:"nonce"`
    Private *big.Int `json:"private"`
}


type StealthSession struct {
    MyPublic CurvePoint `json:"myPublic"`
    TheirPublic CurvePoint `json:"theirPublic"`
    SharedSecret []byte `json:"sharedSecret"`
    TheirAddresses []StealthAddress `json:"theirStealthAddresses"`
    MyAddresses []PrivateStealthAddress `json:"myStealthAddresses"`
}


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



var bigZero = new(big.Int).SetInt64(int64(0))
var bigOne = new(big.Int).SetInt64(int64(1))


// isValidSecretKey checks if the secret can be used to derive
// a valid curve point
//
func isValidSecretKey (secret *big.Int) bool {
    // < 1
    if secret.Cmp(bigOne) < 0 {
        return false
    }

    // >= G
    if secret.Cmp(group.N) >= 0 {
        return false
    }

    return true
}


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
func StealthPubDerive(mpk *CurvePoint, secret []byte) *CurvePoint {
    if ! mpk.IsOnCurve() {
        return nil
    }

    // X ← H(secret)
    _hashout := sha256.Sum256(secret)
    X := new(big.Int).SetBytes(_hashout[:])

    // Y ← g^X
    Y := derivePublicKey(X)

    // spk ← mpk + Y
    spk := mpk.Add(Y)

    return &spk
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
func StealthPrivDerive(msk *big.Int, secret []byte) *big.Int {
    if false == isValidSecretKey(msk) {
        return nil
    }

    // X ← H(secret)
    _hashout := sha256.Sum256(secret)
    X := new(big.Int).SetBytes(_hashout[:])

    // ssk ← msk + X    
    Y := new(big.Int).Add(msk, X)

    // XXX: can (msk + X) exceed group.N?
    ssk := new(big.Int).Mod(Y, group.N)
    if ! derivePublicKey(ssk).IsOnCurve() {
        // TODO: return error?
        return nil;
    }

    return ssk
}


// derivePublicKey derives from SecretKey using ScalarBaseMult:
//o
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
func deriveSharedSecret (myPriv *big.Int, theirPub *CurvePoint) []byte {
    // TODO: theirPub.IsOnCurve?
    // TODO: isValidSecretKey(myPriv)
    // See: RFC5903 (Section 9)
    return theirPub.ScalarMult(myPriv).X.Bytes()
}


func NewStealthSession (mySecret *big.Int, theirPublic *CurvePoint, nonceOffset int, addressCount int) *StealthSession {
    var theirAddresses []StealthAddress
    var myAddresses []PrivateStealthAddress

    if false == isValidSecretKey(mySecret) {
        return nil
    }

    sharedSecret := deriveSharedSecret(mySecret, theirPublic)
    for i := 0; i < addressCount; i++ {
        nonce := new(big.Int).SetInt64(int64(nonceOffset + i))
        secret := append(sharedSecret, nonce.Bytes()...)

        theirStealthPub := StealthPubDerive(theirPublic, secret)
        if theirStealthPub == nil {
            // TODO: return error
            return nil
        }
        theirSA := StealthAddress{*theirStealthPub, nonce}
        theirAddresses = append(theirAddresses, theirSA)

        myStealthPriv := StealthPrivDerive(mySecret, secret)
        myStealthPub := derivePublicKey(myStealthPriv)
        mySA := PrivateStealthAddress{myStealthPub, nonce, myStealthPriv}
        myAddresses = append(myAddresses, mySA)
    }

    session := StealthSession{
        MyPublic: derivePublicKey(mySecret),
        TheirPublic: *theirPublic,
        SharedSecret: sharedSecret,
        TheirAddresses: theirAddresses,
        MyAddresses: myAddresses,
    }

    return &session
}
