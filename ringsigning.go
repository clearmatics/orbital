// Copyright (C) 2017 Clearmatics - All Rights Reserved

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	secp "github.com/btcsuite/btcd/btcec"
	// "echash" the function we want that comes from here is called HashtoBN and takes a byte slice but actually could be edited to take big.Int?!
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
)

var Group *secp.KoblitzCurve

type PubKeyStr struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type RingStr struct {
	PubKeys []PubKeyStr `json:"pubkeys"`
}

type PrivKeysStr struct {
	Keys []string `json:"privkeys"`
}

type PubKey struct {
	CurvePoint
}

type ContractJSON struct {
	keys   []*big.Int
	tau    []*big.Int
	ctlist []*big.Int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {

	Group = secp.S256()
}

func genKeys(n int) ([]PubKeyStr, []*big.Int) {

	var privkeys []*big.Int
	var pubkeys []PubKeyStr

	q := Group.P

	for i := 0; i < n; i++ {
		priv, _ := rand.Int(rand.Reader, q)
		privkeys = append(privkeys, priv)

		p := CurvePoint{}.ScalarBaseMult(priv)

		pubkeys = append(pubkeys, PubKeyStr{p.X.String(), p.Y.String()})
	}

	return pubkeys, privkeys

}

func SignAndVerify(Rn Ring, privBN *big.Int, message []byte) (RingSignature, []byte) {

	pub := CurvePoint{}.ScalarBaseMult(privBN)
	signerNumber := keyCompare(pub, Rn)

	signature := RingSign(Rn, privBN, message, signerNumber)

	verif := RingVerif(Rn, message, signature)
	if verif != true {
		fmt.Println(signature)
		log.Fatal("signature failed to verify")
	}

	var ct []*big.Int

	for i := 0; i < len(signature.Ctlist); i++ {
		ct = append(ct, signature.Ctlist[i])
	}

	ctJS, _ := json.MarshalIndent(ct, "", "\t")

	return signature, ctJS
}

func RingSign(R Ring, ski *big.Int, m []byte, signer int) RingSignature {
	N := Group.N

	mR := RingToBytes(R)
	byteslist := append(mR, m...)
	hashp, _ := HashToCurve(byteslist)
	ski.Mod(ski, N)
	hashSP := hashp.ScalarMult(ski)

	n := len(R.PubKeys)
	var ctlist []*big.Int   //This has to be 2n so here we have n = 4 so 2n = 8 :)
	var hashlist []*big.Int //This has to be 4n but Go won't let it be not const so 16 it is :P
	var a, b CurvePoint
	var ri *big.Int
	var e error
	csum := big.NewInt(0)

	for j := 0; j < n; j++ {

		if j != signer {
			cj, e := rand.Int(rand.Reader, N) // this returns *big.Int
			check(e)
			tj, e := rand.Int(rand.Reader, N) // this returns *big.Int tooo
			check(e)

			a = R.PubKeys[j].ParameterPointAdd(tj, cj)

			b = hashp.HashPointAdd(hashSP, tj, cj)
			ctlist = append(ctlist, cj)
			ctlist = append(ctlist, tj)
			csum.Add(csum, cj)
		}

		if j == signer {
			dummy := big.NewInt(0)
			ctlist = append(ctlist, dummy)
			ctlist = append(ctlist, dummy)
			ri, e = rand.Int(rand.Reader, N)
			check(e)
			a = CurvePoint{}.ScalarBaseMult(ri)
			b = hashp.ScalarMult(ri)
		}
		hashlist = addtolist(hashlist, a.X, a.Y, b.X, b.Y, j)
	}
	for _, v := range hashlist {
		xx := v.Bytes()
		byteslist = append(byteslist, xx[:]...)
	}

	hasha := sha256.Sum256(byteslist)
	hashb := Convert(hasha[:])
	hashb.Mod(hashb, N)
	csum.Mod(csum, N)
	c := new(big.Int).Sub(hashb, csum)
	c.Mod(c, N)

	cx := new(big.Int).Mul(c, ski)
	cx.Mod(cx, N)
	ti := new(big.Int).Sub(ri, cx)
	ti.Mod(ti, N)
	ctlist[2*signer] = c
	ctlist[2*signer+1] = ti

	return RingSignature{hashSP, ctlist}
}

func RingVerif(R Ring, m []byte, sigma RingSignature) bool {
	// ring verification
	// assumes R = pk1, pk2, ..., pkn
	// sigma = H(m||R)^x_i, c1, t1, ..., cn, tn = taux, tauy, c1, t1, ..., cn, tn
	tau := sigma.Tau
	ctlist := sigma.Ctlist
	n := len(R.PubKeys)
	N := Group.N
	var hashlist []*big.Int

	mR := RingToBytes(R)
	byteslist := append(mR, m...)
	hashp, _ := HashToCurve(byteslist)
	csum := big.NewInt(0)

	for j := 0; j < n; j++ {
		cj := ctlist[2*j]
		tj := ctlist[2*j+1]
		cj.Mod(cj, N)
		tj.Mod(tj, N)
		H := hashp.ScalarMult(tj)             //H(m||R)^t
		gt := CurvePoint{}.ScalarBaseMult(tj) //g^t
		yc := R.PubKeys[j].ScalarMult(cj)     // y^c = g^(xc)
		tauc := tau.ScalarMult(cj)            //H(m||R)^(xc)
		gt = gt.Add(yc)
		H = H.Add(tauc) // fieldJacobianToBigAffine `normalizes' values before returning so yes - normalize uses fast reduction using specialised form of secp256k1's prime! :D
		hashlist = addtolist(hashlist, gt.X, gt.Y, H.X, H.Y, j)
		csum.Add(csum, cj)
	}
	for _, v := range hashlist {
		xx := v.Bytes()
		byteslist = append(byteslist, xx[:]...)
	}

	hash := sha256.Sum256(byteslist)
	hashhash := Convert(hash[:])

	hashhash.Mod(hashhash, N)
	csum.Mod(csum, N)
	if csum.Cmp(hashhash) == 0 {
		return true
	}
	return false
}

func convertPubKeys(rn RingStr) Ring {

	rl := len(rn.PubKeys)
	var ring Ring

	for i := 0; i < rl; i++ {
		var bytesx []byte
		var bytesy []byte
		bytesx, _ = hex.DecodeString(string(rn.PubKeys[i].X))
		bytesy, _ = hex.DecodeString(string(rn.PubKeys[i].Y))
		pubkeyx := new(big.Int).SetBytes(bytesx) // This makes big int
		pubkeyy := new(big.Int).SetBytes(bytesy) // So we can do EC arithmetic
		ring.PubKeys = append(ring.PubKeys, PubKey{CurvePoint{pubkeyx, pubkeyy}})
	}
	return ring
}

func RingToBytes(rn Ring) []byte {
	var rbytes []byte
	for i := 0; i < len(rn.PubKeys); i++ {
		rbytes = append(rbytes, rn.PubKeys[i].X.Bytes()...)
	}
	for i := 0; i < len(rn.PubKeys); i++ {
		rbytes = append(rbytes, rn.PubKeys[i].Y.Bytes()...)
	}

	return rbytes
}

func addtolist(list []*big.Int, a *big.Int, b *big.Int, c *big.Int, d *big.Int, j int) []*big.Int {
	list = append(list, a)
	list = append(list, b)
	list = append(list, c)
	list = append(list, d)
	return list
}

func keyCompare(pub CurvePoint, R Ring) int {
	j := 0
	for i := 0; i < len(R.PubKeys); i++ {
		if pub.X.Cmp(R.PubKeys[i].X) == 0 && pub.Y.Cmp(R.PubKeys[i].Y) == 0 {
			j = i
		}
	}
	return j
}

func Convert(data []byte) *big.Int {
	z := new(big.Int)
	z.SetBytes(data)
	return z
}

func HashToCurve(s []byte) (CurvePoint, error) {
	q := Group.P

	x := big.NewInt(0)
	y := big.NewInt(0)
	z := big.NewInt(0)
	z.SetString("57896044618658097711785492504343953926634992332820282019728792003954417335832", 10)

	array := sha256.Sum256(s) // Sum outputs an array of 32 bytes :)
	x = Convert(array[:])
	for true {
		xcube := new(big.Int).Exp(x, big.NewInt(3), q)
		xcube7 := new(big.Int).Add(xcube, big.NewInt(7))
		y.ModSqrt(xcube7, q)
		y.Set(q)
		y.Add(y, big.NewInt(1))
		y.Rsh(y, 2)
		y.Exp(xcube7, y, q)
		z = z.Exp(y, big.NewInt(2), q)
		curveout := Group.IsOnCurve(x, y)
		if curveout == true {
			return CurvePoint{x, y}, nil
		}
		x.Add(x, big.NewInt(1))
	}
	return CurvePoint{}, errors.New("no curve point found")
}

func create(privateKeysFile string, publicKeysFile string, outputFilename string, rawMessage string) error {
	outputBuffer, err := Process(privateKeysFile, publicKeysFile, rawMessage)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputFilename, outputBuffer, 0777)
}

func verify(privateKeysFile string, publicKeysFile string, outputFilename string, rawMessage string) (bool, error) {
	outputBuffer, err := Process(privateKeysFile, publicKeysFile, rawMessage)

	if err != nil {
		return false, err
	}

	var compareBuffer []byte
	compareBuffer, err = ioutil.ReadFile(outputFilename)
	if err != nil {
		return false, err
	}

	return bytes.Equal(outputBuffer, compareBuffer), nil
}

func GenerateRandomRing(ringSize int) (Ring, []PubKeyStr, []*big.Int) {
	var sks []*big.Int
	var pks []PubKeyStr
	// generate keypair (private and public)
	pks, sks = genKeys(ringSize)
	// populate ring with keypairs
	var ring Ring
	for i := 0; i < len(pks); i++ {
		// type cast to the rign struct
		xPub := new(big.Int)
		xPub.SetString(pks[i].X, 10)
		yPub := new(big.Int)
		yPub.SetString(pks[i].Y, 10)

		// fills the key ring
		ring.PubKeys = append(ring.PubKeys, PubKey{CurvePoint{xPub, yPub}})
	}

	return ring, pks, sks
}

func ProcessSignature(ring Ring, privateKeys []*big.Int, message []byte) ([]RingSignature, error) {

	// generate signature
	var signaturesArr []RingSignature
	for i := 0; i < len(privateKeys); i++ {
		privKey := privateKeys[i]
		// signing function
		signature, _ /*ctlist*/ := SignAndVerify(ring, privKey, message)
		signaturesArr = append(signaturesArr, signature)
	}
	return signaturesArr, nil
}

func Process(privateKeysFile string, publicKeysFile string, rawMessage string) ([]byte, error) {
	privkeyfile, err := ioutil.ReadFile(privateKeysFile)
	pk := PrivKeysStr{}
	if err = json.Unmarshal(privkeyfile, &pk); err != nil {
		return nil, err
	}

	keyfile, _ := ioutil.ReadFile(publicKeysFile)
	rn := RingStr{}
	if err = json.Unmarshal(keyfile, &rn); err != nil {
		return nil, err
	}
	Rn := convertPubKeys(rn)

	var message []byte
	if rawMessage != "" {
		message, err = hex.DecodeString(rawMessage)
		if err != nil {
			return nil, err
		}
	}

	var outputBuffer []byte

	for i := 0; i < len(pk.Keys); i++ {
		privbytes, err := hex.DecodeString(pk.Keys[i])
		if err != nil {
			return nil, err
		}
		privBN := new(big.Int).SetBytes(privbytes)

		signature, ctList := SignAndVerify(Rn, privBN, message)
		outputBuffer = append(outputBuffer, []byte(fmt.Sprintln("["))...)
		outputBuffer = append(outputBuffer, []byte(fmt.Sprintf("\t%s,\n", signature.Tau.X))...)
		outputBuffer = append(outputBuffer, []byte(fmt.Sprintf("\t%s,\n", signature.Tau.Y))...)
		outputBuffer = append(outputBuffer, []byte(fmt.Sprintf("%s\n", string(ctList)))...)
		outputBuffer = append(outputBuffer, []byte(fmt.Sprintln("],"))...)
	}

	return outputBuffer, nil
}
