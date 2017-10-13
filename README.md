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

    orbital -geninputs n HexEncodedString

# Building

- Clean the build artifacts: make clean
- Build the binary: make
- Build and run the unittests: make test
- Run gofmt over the whole tree: make format

<<<<<<< HEAD
=======
# Development

Dependencies are managed via [dep][1]. Dependencies are checked into this repository in the `vendor` folder. Documentation for managing dependencies is available on the [dep README][2].

The project follows standard go conventions using `gofmt`. If you wish to contribute to the project please follow standard Go conventions. The CI server automatically runs these checks.

[1]: https://github.com/golang/dep
[2]: https://github.com/golang/dep/blob/master/README.md
>>>>>>> cleanup
