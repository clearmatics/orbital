package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type ProcessFlag func(firstarg string, otherargs []string)

type CmdOption struct {
	LongFlag  string
	ShortFlag string
	Process   ProcessFlag
	Usage     string
}

type ProcessOption struct {
	FlagVar *string
	Process ProcessFlag
}

var subCommands = []CmdOption{
	{"create", "c", processCreate, "create privkeys.json pubkeys.json output.json HexEncodedString"},
	{"genkeys", "g", processKeygen, "genkeys n privkeys.json pubkeys.json"},
	{"geninputs", "i", processGenInputs, "geninputs n HexEncodedString"},
}

func createStringFlag(longFlag string, shortFlag string, help string) *string {
	var flagvar string
	flag.StringVar(&flagvar, longFlag, "", help)
	flag.StringVar(&flagvar, shortFlag, "", help)

	return &flagvar
}

func main() {
	var options []ProcessOption

	for i := 0; i < len(subCommands); i++ {
		flagVar := createStringFlag(subCommands[i].LongFlag, subCommands[i].ShortFlag, subCommands[i].Usage)
		newOption := ProcessOption{flagVar, subCommands[i].Process}
		options = append(options, newOption)
	}

	flag.Parse()

	otherArgs := flag.Args()

	for i := 0; i < len(options); i++ {
		currentFlag := *options[i].FlagVar
		if currentFlag != "" {
			options[i].Process(currentFlag, otherArgs)
			return
		}
	}
}

func hexString2Bytes(rawMessage string) []byte {
	var message []byte
	var err error
	if rawMessage != "" {
		message, err = hex.DecodeString(rawMessage)
		if err == nil {
			return message
		}
	}
	// FAIL
	log.Fatal("Failed to parse the message")
	return nil
}

func processGenInputs(firstarg string, otherargs []string) {
	// regex just to put numbers between quotes
	re := regexp.MustCompile("([0-9]+)")

	var sks []*big.Int
	var pks []PubKeyStr

	n, _ := strconv.Atoi(firstarg)

	// generate key ring
	ring, pks, sks := GenerateRandomRing(n)

	// message hexadecimal string to bytes
	rawMessage := otherargs[0]
	message := hexString2Bytes(rawMessage)

	// generate signature and smart contract withdraw and deposit input data
	signature, _ := ProcessSignature(ring, sks, message)

	// verify signature
	verif := true
	for i := 0; i < len(signature); i++ {
		verif = RingVerif(ring, message, signature[i])
		if !verif {
			// FAIL
			log.Fatal("Failed to verify ring signature")
			return
		}
	}

	// print result
	pkJSON, _ := json.MarshalIndent(pks, "  ", "  ")
	signatureJSON, _ := json.MarshalIndent(signature, "  ", "  ")
	signatureJSONStr := re.ReplaceAllString(string(signatureJSON), "\"${1}\"")
	resultStr := "{\n  \"deposit_input\": " + string(pkJSON) + ",\n  \"withdraw_input\": " + signatureJSONStr + "\n}"
	fmt.Printf("%s\n", resultStr)
}

func processKeygen(firstarg string, otherargs []string) {

	var sks []*big.Int
	var pks []PubKeyStr

	n, _ := strconv.Atoi(firstarg)

	pks, sks = genKeys(n)

	sksJSON, _ := json.MarshalIndent(sks, "", "\t")
	pksJSON, _ := json.MarshalIndent(pks, "", "\t")

	ioutil.WriteFile(otherargs[0], []byte(string(sksJSON)), 0777)
	ioutil.WriteFile(otherargs[1], []byte(string(pksJSON)), 0777)

}

func processCreate(firstarg string, otherargs []string) {
	if len(otherargs) == 3 && strings.HasPrefix(otherargs[2], "0x") {
		privateKeysFile := firstarg
		publicKeysFile := otherargs[0]
		outputFilename := otherargs[1]
		rawMessage := otherargs[2]
		stripedPrefixMessage := rawMessage[2:]

		err := create(privateKeysFile, publicKeysFile, outputFilename, stripedPrefixMessage)

		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("syntax error")
	}
}

func processVerify(firstarg string, otherargs []string) {
	if len(otherargs) == 3 && strings.HasPrefix(otherargs[2], "0x") {
		privateKeysFile := firstarg
		publicKeysFile := otherargs[0]
		outputFilename := otherargs[1]
		rawMessage := otherargs[2]
		stripedPrefixMessage := rawMessage[2:]

		match, err := verify(privateKeysFile, publicKeysFile, outputFilename, stripedPrefixMessage)

		if err != nil {
			fmt.Println("Somethnig went wrong")
		} else if match {
			fmt.Println("Signatures match")
		} else {
			fmt.Println("Signatures do not match")
		}
	} else {
		fmt.Println("syntax error")
	}
}
