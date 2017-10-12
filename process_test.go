package main_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/clearmatics/orbital"
)

func Test_Process_MissingPrivateKeys_ReturnSuccess(t *testing.T) {
	privateKeysFile := "fakeFile.json"
	publicKeysFile := "data/pubkeys.json"
	rawMessage := "1234567890"

	outputBuffer, processError := main.Process(privateKeysFile, publicKeysFile, rawMessage)

	assert.NotEqual(t, nil, processError)
	assert.Equal(t, []byte(nil), outputBuffer)
}

func Test_Process_MissingPublicKeys_ReturnSuccess(t *testing.T) {
	privateKeysFile := "data/privkeys.json"
	publicKeysFile := "fakeFile.json"
	rawMessage := "1234567890"

	outputBuffer, processError := main.Process(privateKeysFile, publicKeysFile, rawMessage)

	assert.NotEqual(t, nil, processError)
	assert.Equal(t, []byte(nil), outputBuffer)
}

func Test_Process_MissingMessage_ReturnSuccess(t *testing.T) {
	privateKeysFile := "data/privkeys.json"
	publicKeysFile := "fakeFile.json"
	rawMessage := ""

	outputBuffer, processError := main.Process(privateKeysFile, publicKeysFile, rawMessage)

	assert.NotEqual(t, nil, processError)
	assert.Equal(t, []byte(nil), outputBuffer)
}
