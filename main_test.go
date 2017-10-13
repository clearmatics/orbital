package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Process_MissingPrivateKeys_ReturnSuccess(t *testing.T) {
	privateKeysFile := "fakeFile.json"
	publicKeysFile := "data/pubkeys.json"
	rawMessage := "1234567890"

	outputBuffer, processError := Process(privateKeysFile, publicKeysFile, rawMessage)

	assert.NotEqual(t, nil, processError)
	assert.Equal(t, []byte(nil), outputBuffer)
}

func Test_Process_MissingPublicKeys_ReturnSuccess(t *testing.T) {
	privateKeysFile := "data/privkeys.json"
	publicKeysFile := "fakeFile.json"
	rawMessage := "1234567890"

	outputBuffer, processError := Process(privateKeysFile, publicKeysFile, rawMessage)

	assert.NotEqual(t, nil, processError)
	assert.Equal(t, []byte(nil), outputBuffer)
}

func Test_Process_MissingMessage_ReturnSuccess(t *testing.T) {
	privateKeysFile := "data/privkeys.json"
	publicKeysFile := "fakeFile.json"
	rawMessage := ""

	outputBuffer, processError := Process(privateKeysFile, publicKeysFile, rawMessage)

	assert.NotEqual(t, nil, processError)
	assert.Equal(t, []byte(nil), outputBuffer)
}

func Test_Process_GenInputs(t *testing.T) {
	// generate key ring
	n := 4
	ring, _, sks := GenerateRandomRing(n)

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
