# Orbital

Generate off-chain data for a Mobius smart contract

## Installation

    go get github.com/clearmatics/orbital

## Usage

Create a signature

    orbital -create privkeys.json pubkeys.json output.json HexEncodedString

Verify a signature

    orbital -verify privkeys.json pubkeys.json output.json HexEncodedString`

Print Ring smart contract random inputs (creates a random keypair and a signature from those)

    orbital -geninputs n HexEncodedString`

# Building

- Clean the build artifacts: make clean
- Build the binary: make
- Build and run the unittests: make test
- Run gofmt over the whole tree: make format

