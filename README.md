# Orbital

Generate off-chain data for a Mobius smart contract

## Installation

    go get github.com/clearmatics/orbital

## Usage

Print Ring smart contract random inputs (creates a random keypair and a signature from those)

    orbital -geninputs n HexEncodedString

Print randomly generated keys

    orbital -genkeys n

Print signature and public keys ring from keys file

    orbital -signature keys.json HexEncodedString

Verify signatures

    orbital -verify signature.json HexEncodedString

Example
```
$ ./orbital -genkeys 4 > keys.json
$ ./orbital -signature keys.json 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc > ringSignature.json
$ ./orbital -verify ringSignature.json 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc
Signatures verified
```
or
```
$ ./orbital -geninputs 4 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc > ringSignature.json
$ ./orbital -verify ringSignature.json 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc
Signatures verified
```

# Building

- Clean the build artifacts: make clean
- Build the binary: make
- Build and run the unittests: make test
- Run gofmt over the whole tree: make format

# Development

Dependencies are managed via [dep][1]. Dependencies are checked into this repository in the `vendor` folder. Documentation for managing dependencies is available on the [dep README][2].

The project follows standard go conventions using `gofmt`. If you wish to contribute to the project please follow standard Go conventions. The CI server automatically runs these checks.

[1]: https://github.com/golang/dep
[2]: https://github.com/golang/dep/blob/master/README.md
