Copyright (C) 2017 Clearmatics - All Rights Reserved

Clearmatics Transaction Privacy Tools

# Building

- Clean the build artifacts: make clean
- Build the binary: make
- Build and run the unittests: make test
- Run gofmt over the whole tree: make format

# Running

- Create a signature
    ./go -create privkeys.json pubkeys.json output.json HexEncodedString

- Verify a signature
    ./go -verify privkeys.json pubkeys.json output.json HexEncodedString
