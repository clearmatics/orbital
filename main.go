package main

import (
	"encoding/json"
	"flag"
	"fmt"
	//"io/ioutil"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/big"
	"regexp"
	"strconv"
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

type SignatureData struct {
	Ring       []CurvePoint    `json:"ring"`
	Signatures []RingSignature `json:"signatures"`
}

type KeyPair struct {
	Private []*big.Int   `json:"private"`
	Public  []CurvePoint `json:"public"`
}

var subCommands = []CmdOption{
	{"genkeys", "g", processKeygen, "genkeys n"},
	{"geninputs", "i", processGenInputs, "geninputs n HexEncodedString"},
	{"signature", "c", processGenerateSignature, "signature keys.json HexEncodedString"},
	{"verify", "v", processVerifySignature, "verify ringSignature.json HexEncodedString"},
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
	var sks []*big.Int
	var pks []PubKey

	n, err := strconv.Atoi(firstarg)
	if err != nil {
		log.Fatal("Failed to parse amount: %s", err)
	}

	// generate key ring
	ring, pks, sks := GenerateRandomRing(n)

	// message hexadecimal string to bytes
	rawMessage := otherargs[0]
	message := hexString2Bytes(rawMessage)

	// generate signature and smart contract withdraw and deposit input data
	signatureArr, err := ProcessSignature(ring, sks, message)
	if err != nil {
		log.Fatal("Failed to create the signature: %s", err)
	}

	// regex just to put numbers between quotes
	re := regexp.MustCompile("([0-9]+)")
	// print result
	pkJSON, _ := json.MarshalIndent(pks, "  ", "  ")
	signatureJSON, _ := json.MarshalIndent(signatureArr, "  ", "  ")
	signatureJSONStr := re.ReplaceAllString(string(signatureJSON), "\"${1}\"")
	resultStr := "{\n  \"ring\": " + string(pkJSON) + ",\n  \"signatures\": " + signatureJSONStr + "\n}"
	fmt.Printf("%s\n", resultStr)
}

func processVerifySignature(firstarg string, otherargs []string) {
	signatureFile := firstarg
	rawMessage := otherargs[0]

	signatureFileData, err := ioutil.ReadFile(signatureFile)
	if err != nil {
		// FAIL
		log.Fatal("Failed to parse file: %s", err)
		//return
	}
	// Integers can't have " to parse
	re := regexp.MustCompile("(\"([0-9]+)\")")
	signatureFileDataStr := re.ReplaceAllString(string(signatureFileData), "${2}")
	// json to golang object
	var signatureData SignatureData
	err = json.Unmarshal([]byte(signatureFileDataStr), &signatureData)
	if err != nil {
		// FAIL
		log.Fatal("Failed to parse file: %s", err)
		//return
	}

	// create ring
	var ring Ring
	for _, pubKey := range signatureData.Ring {
		ring.PubKeys = append(ring.PubKeys, PubKey{pubKey})
	}

	// message hexadecimal string to bytes
	message := hexString2Bytes(rawMessage)

	// verify signature
	for i, signature := range signatureData.Signatures {
		if !RingVerif(ring, message, signature) {
			// FAIL
			log.Fatal("Failed to verify ring signature number " + string(i))
			return
		}
	}
	fmt.Println("Signatures verified")
}

func processGenerateSignature(firstarg string, otherargs []string) {
	keysFile := firstarg
	rawMessage := otherargs[0]

	keysFileData, err := ioutil.ReadFile(keysFile)
	// Integers can't have " to parse
	re := regexp.MustCompile("(\"([0-9]+)\")")
	keysFileDataStr := re.ReplaceAllString(string(keysFileData), "${2}")

	var keyPair KeyPair
	err = json.Unmarshal([]byte(keysFileDataStr), &keyPair)
	if err != nil {
		// FAIL
		log.Fatal("Failed to parse file: %s", err)
		//return
	}

	// message hexadecimal string to bytes
	message := hexString2Bytes(rawMessage)

	sks := keyPair.Private
	// create ring
	var ring Ring
	for _, pubKey := range keyPair.Public {
		ring.PubKeys = append(ring.PubKeys, PubKey{pubKey})
	}

	// generate signature
	signature, err := ProcessSignature(ring, sks, message)
	if err != nil {
		log.Fatal("Failed to create the signature: %s", err)
	}

	// regex just to put numbers between quotes
	re = regexp.MustCompile("([0-9]+)")
	// print result
	pkJSON, _ := json.MarshalIndent(keyPair.Public, "  ", "  ")
	pksJSONStr := re.ReplaceAllString(string(pkJSON), "\"${1}\"")
	signatureJSON, _ := json.MarshalIndent(signature, "  ", "  ")
	signatureJSONStr := re.ReplaceAllString(string(signatureJSON), "\"${1}\"")
	resultStr := "{\n  \"ring\": " + string(pksJSONStr) + ",\n  \"signatures\": " + signatureJSONStr + "\n}"
	fmt.Printf("%s\n", resultStr)
}

func processKeygen(firstarg string, otherargs []string) {
	var sks []*big.Int
	var pks []PubKey

	n, err := strconv.Atoi(firstarg)
	if err != nil {
		log.Fatal("Failed to parse amount: %s", err)
	}

	// generate key ring
	_, pks, sks = GenerateRandomRing(n)

	// print keys
	var sksStrArr []string
	for _, privateKey := range sks {
		sksStrArr = append(sksStrArr, privateKey.String())
	}
	// regex just to put numbers between quotes
    re := regexp.MustCompile("([0-9]+)")
	sksJSON, _ := json.MarshalIndent(sksStrArr, "  ", "  ")
	pksJSON, _ := json.MarshalIndent(pks, "  ", "  ")
	pksJSONStr := re.ReplaceAllString(string(pksJSON), "\"${1}\"")
	fmt.Printf("{\n  \"private\": %s,\n  \"public\": %s\n}\n", sksJSON, pksJSONStr)
}
