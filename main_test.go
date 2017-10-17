package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Process_GenInputs(t *testing.T) {
	// generate key ring
	n := 4
    pks, sks := GenKeys(n)
	var ring Ring
	ring.PubKeys = pks

	// message hexadecimal string to bytes
	rawMessage := "50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc"
	message := hexString2Bytes(rawMessage)

	// generate signature and smart contract withdraw and deposit input data
	signature, _ := ProcessSignature(ring, sks, message)

	// verify signature
	verif := true
	for i := 0; i < len(signature); i++ {
		verif = verif && RingVerif(ring, message, signature[i])
	}
	assert.True(t, verif)
}
