// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func flagUsage() {
	usageText := `Orbital generates off-chain data for Möbius contracts 

	Usage:
	orbital command [arguments]
	The commands are:
	generate	Generate public/private key pairs for a contract
	inputs		Generate data inputs for a contract 
	verify		Verify a set of public keys against signatures
	stealth		Generate stealth addresses
	Use "orbital [command] --help" for more information about a command.`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func main() {
	flag.Usage = flagUsage

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	stealthCmd := flag.NewFlagSet("stealth", flag.ExitOnError)
	inputsCmd := flag.NewFlagSet("inputs", flag.ExitOnError)
	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "stealth":
		n := stealthCmd.Int("n", 1, "Number of addresses to generate")
		nonceOffset := stealthCmd.Int("o", 0, "Nonce offset")
		_mySecretKey := stealthCmd.String("s", "", "Your secret key")
		theirPublicKeyX := stealthCmd.String("x", "", "Their public key X point")
		theirPublicKeyY := stealthCmd.String("y", "", "Their public key Y point")

		stealthCmd.Parse(os.Args[2:])
		if *n <= 0 || *_mySecretKey == "" || *theirPublicKeyX == "" || *theirPublicKeyY == "" {
			stealthCmd.Usage()
			return
		}

		mySecretKey, errMSK := ParseBigInt(*_mySecretKey) // new(big.Int).SetString(*_mySecretKey, 0)
		if errMSK != nil || mySecretKey == nil {
			fmt.Fprintf(os.Stderr, "Unable to parse secret key: -s %v: %v\n", *_mySecretKey, errMSK)
			os.Exit(1)
		}

		// TODO: optionally parse their public key as a single string, then derive Y point
		theirPublicKey := ParseCurvePoint(*theirPublicKeyX, *theirPublicKeyY)
		if theirPublicKey == nil {
			fmt.Fprintf(os.Stderr, "Unable to parse public key pair -x %v -y %v\n", *theirPublicKeyX, *theirPublicKeyY)
			os.Exit(1)
		}

		session, err := NewStealthSession(mySecretKey, theirPublicKey, *nonceOffset, *n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate stealth session: %v\n", err)
			os.Exit(1)
		}

		saJSON, err := json.MarshalIndent(session, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(saJSON))

	// Generates a set of key pairs to be used for use with later operations
	case "generate":
		i := generateCmd.Int("n", 0, "Number of key pairs to be generated, e.g. 4")
		generateCmd.Parse(os.Args[2:])

		if *i == 0 {
			generateCmd.Usage()
			return
		}

		ring := &Ring{}
		ring.Generate(*i)

		ringJSON, err := json.MarshalIndent(ring, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(ringJSON))

	case "inputs":
		keys_file := inputsCmd.String("k", "", "Load signing keys from a JSON file")
		n := inputsCmd.Int("n", 0, "The size of the ring to be generated e.g. 4")
		m := inputsCmd.String("m", "", "A Hex encoded string to be used to generate the ring")
		inputsCmd.Parse(os.Args[2:])

		if *n == 0 {
			inputsCmd.Usage()
			return
		}
		if *m == "" {
			inputsCmd.Usage()
			return
		}

		ring := &Ring{}

		var stealthSessionAliceToBob *StealthSession
		var stealthSessionBobToAlice *StealthSession

		if *keys_file != "" {
			// Load keys from the ring
			data, err := ioutil.ReadFile(*keys_file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to read key file '%v': %v\n", *keys_file, err)
				os.Exit(1)
			}

			err = json.Unmarshal(data, &ring)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to parse keys file '%v': %v\n", *keys_file, err)
				os.Exit(1)
			}
		} else {
			// Otherwise, generate a stealth session, as an example
			alicePub, alicePriv, _ := generateKeyPair()
			bobPub, bobPriv, _ := generateKeyPair()
			stealthSessionAliceToBob, err := NewStealthSession(alicePriv, bobPub, 0, 1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to derive stealth session, Alice->Bob: %v\n", err)
				os.Exit(1)
			}
			stealthSessionBobToAlice, err := NewStealthSession(bobPriv, alicePub, 0, 1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to derive stealth session, Alice->Bob: %v\n", err)
				os.Exit(1)
			}

			// Then generate some random key pairs and integrate our stealth session into the ring
			ring.Generate(*n)
			ring.PrivKeys[0] = stealthSessionBobToAlice.MyAddresses[0].Private
			ring.PubKeys[0] = stealthSessionAliceToBob.TheirAddresses[0].Public			
		}


		decoded, err := hex.DecodeString(*m)
		if err != nil {
			panic(err)
		}

		signatures, err := ring.Signatures(decoded)
		if err != nil {
			panic(err)
		}

		inputData := inputData{
			PubKeys:    ring.PubKeys,
			Signatures: signatures,
			Message:    decoded,
			AliceToBob: stealthSessionAliceToBob,
			BobToAlice: stealthSessionBobToAlice,
		}

		ringJSON, err := json.MarshalIndent(inputData, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(ringJSON))
		os.Exit(0)

	case "verify":
		var inputData inputData

		f := verifyCmd.String("f", "", "Path to a JSON file containing public keys and signatures")
		m := verifyCmd.String("m", "", "The Hex encoded message used to generate the ring")
		verifyCmd.Parse(os.Args[2:])

		if *f == "" {
			verifyCmd.Usage()
			return
		}

		if *m == "" {
			verifyCmd.Usage()
			return
		}

		decoded, err := hex.DecodeString(*m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse hex string: %v\n", err)
			os.Exit(1)
		}

		data, err := ioutil.ReadFile(*f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read file: %v\n", err)
			os.Exit(1)
		}

		json.Unmarshal(data, &inputData)

		r := Ring{
			PubKeys: inputData.PubKeys,
		}

		for _, sig := range inputData.Signatures {
			valid := r.VerifySignature(decoded, sig)
			if valid != true {
				fmt.Fprintln(os.Stderr, "Signatures not verified")
				os.Exit(1)
			}
		}
		fmt.Println("Signatures verified")
		os.Exit(0)

	default:
		flag.Usage()
	}
}
