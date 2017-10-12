// Copyright (C) 2017 Clearmatics - All Rights Reserved

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
    "regexp"
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

func processGenInputs(firstarg string, otherargs []string) {
	var sks []*big.Int
	var pks []PubKeyStr

	n, _ := strconv.Atoi(firstarg)
    message := otherargs[0]

    // generate keypair (private and public)
	pks, sks = genKeys(n)

    // generate signature and smart contract withdraw and deposit input data
    signature, _ := ProcessSignature(pks,sks,message)

    signatureJson, _ := json.MarshalIndent(signature, "", "  ")
    // regex just to put numbers between quotes
    re := regexp.MustCompile("([0-9]+)")
    signatureJsonStr := re.ReplaceAllString(string(signatureJson),"\"${1}\"")
    fmt.Printf("%s\n",signatureJsonStr)
}

func processKeygen(firstarg string, otherargs []string) {

	var sks []*big.Int
	var pks []PubKeyStr

	n, _ := strconv.Atoi(firstarg)

	pks, sks = genKeys(n)

	sksJson, _ := json.MarshalIndent(sks, "", "\t")
	pksJson, _ := json.MarshalIndent(pks, "", "\t")

	ioutil.WriteFile(otherargs[0], []byte(string(sksJson)), 0777)
	ioutil.WriteFile(otherargs[1], []byte(string(pksJson)), 0777)

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
