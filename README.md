# Orbital

[![Build Status](https://travis-ci.org/clearmatics/bn256.svg?branch=ci)](https://travis-ci.org/clearmatics/bn256)

Orbital is a command-line tool to generate off-chain data required by [Möbius][3], a smart contract that offers trustless tumbling for transaction privacy.

## Prerequisites

A version of Go >= 1.8 is required. The [dep][1] tool is used for dependency management. 

## Installation

    go get github.com/clearmatics/orbital

## Usage

When deployed a Möbius contract will emit a `MixerMessage` that is an arbitrary hex encoded string. This message is signed to make withdrawals from the contract. 

Orbital can be used to generate all data needed to deposit and withdraw from a Möbius smart contract. Providing you have the `MixerMessage` value data can be generated as follows. In this example the hex encoded string is given as `291a6780850827fcd8621...`. A ring size of 2 is generated.

    ./orbital inputs -n 2 -m 291a6780850827fcd8621d0e5471343831109bc14142ec101527b048bb3d1794

This generates JSON containing complete data required to deposit and withdraw. If you are just evaluating Möbius this is all you need to deposit into a contract and then make withdrawals once the ring is full. 

``` JSON
{
  "alice2bob": null,
  "bob2alice": null,
  "message": "KRpngIUIJ/zYYh0OVHE0ODEQm8FBQuwQFSewSLs9F5Q=",
  "ring": [
    {
      "x": "0x20495bff8bfbedb0e9a4a4ec87166340199540b466eb128f42f8e70030d0864",
      "y": "0x12526493c446360fe1c8185730f2df1b2c80b3f107bf41e0c1362c636bd35bf2"
    },
    {
      "x": "0xa612e73ef366816f4354c49ed161659073b1b3a8f559d2c99fecd2766e8876d",
      "y": "0x1236505b328b0074cdd730782c34852a174c9c257d879250cd7690c4220ee1c7"
    }
  ],
  "signatures": [
    {
      "tau": {
        "x": "0x22604e4e03a513070439fd8f33e2f987c6e5f399251769d982c4a697aa397ee1",
        "y": "0x104bc20e2d6922a728b65713cc9783de07969763ad87e1d63873f20ba3e87c36"
      },
      "ctlist": [
        "0x2e0cc2449b13c1f42f874610de1a5bc1275fafe759a1442c3d623de7617d6c6f",
        "0x179ef01a92ca56c03ce64b5c03966e881f9fc91dac228d6f4687855c58520565",
        "0x21812114778a8c80381e747b2b9cbbbc437e8f24d10f8948a937a06538e19941",
        "0x10ac35c6eb4c317e6ef7fcaa4676b374bba694770605fe56d7e03ac96af4eb39"
      ]
    },
    {
      "tau": {
        "x": "0x2b5108123f8a315adc4c1f79cdac8b3bcbbbdf379a2281b2c14abb310a6118ee",
        "y": "0x17a2d1d0e35fcfa79d174689c37f9a6f122205eb75e39bd6f68597d6843031c8"
      },
      "ctlist": [
        "0x1052e41400969dc1a28d77aa6b323292c47f09dd664361440648873a01edaad",
        "0x20f7c07495a2aa100d80b39d30266274b59908e233304b5f18092617be47e8ee",
        "0x20f7c25434304a59197e32f306eb94e1c51364f4dfc7e95914677c0732cba48a",
        "0xb747984ad0ac9b649790fe257b8769ecd38f18d8dd8c500a72d30f2281c4cc9"
      ]
    }
  ]
}
```

Beyond evaluation it is likely that you will want to first generate public/private keypairs and distribute this to allow signatures to be created separately. To create a public/private keypairs for a ring size of `2` use 
    
    orbital generate -n 2

To verify signatures pass a file with public keys and signatures to the `verify` subcommand.

    orbital verify -f signatures.json -m 666f6f62617262617a

Example:

```
    $ orbital generate -n 4 > keys.json
    $ orbital inputs -f keys.json -n 4 -m 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc > ringSignature.json
    $ orbital verify -f ringSignature.json -m 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc
    Signatures verified
```

### Stealth Addresses

Integration of Stealth Addresses into the Möbius contract is still in-progress, however they can be generated using the `mobius stealth` utility.

First generate a pair of key pairs:

```
$ orbital generate -n 2
{
  "pubkeys": [
    {
      "x": "0x27e073fe3485b7ab97de5813342c3ce3dd19eba467a6b5f61b092813e67325e2",
      "y": "0x639bbbc72ef12bc7449c08a4ac2d7d7a9ed0484e4d5d6ee4f4b6fa5844043c2"
    },
    {
      "x": "0xf5165507a4f67c89fd40603a3765e53a3656bf3c247421b6a69a05e0a7fabb",
      "y": "0x6d9e4c9c4054e9f1f5d27cfad56c034f1654bc88ae6916722424c565976ef88"
    }
  ],
  "privkeys": [
    "0x12987c779f1c76fdc01efaa3ce7fe8e730512ee4a12055018a44b8c9044fb860",
    "0x282e1c333c40fb8ffcb97a66d615af2a0c4077c5dc5dd50075cf6abf3e9e9f55"
  ]
}
```

Then derive a stealth address for the other party using the first secret key and the second public key, the JSON output displays the shared secret, your public key and their stealth addresses.

```
$ orbital stealth -s 0x12987c779f1c76fdc01efaa3ce7fe8e730512ee4a12055018a44b8c9044fb860 -x 0xf5165507a4f67c89fd40603a3765e53a3656bf3c247421b6a69a05e0a7fabb -y 0x6d9e4c9c4054e9f1f5d27cfad56c034f1654bc88ae6916722424c565976ef88

{
  "myPublic": {
    "x": "0x27e073fe3485b7ab97de5813342c3ce3dd19eba467a6b5f61b092813e67325e2",
    "y": "0x639bbbc72ef12bc7449c08a4ac2d7d7a9ed0484e4d5d6ee4f4b6fa5844043c2"
  },
  "theirPublic": {
    "x": "0xf5165507a4f67c89fd40603a3765e53a3656bf3c247421b6a69a05e0a7fabb",
    "y": "0x6d9e4c9c4054e9f1f5d27cfad56c034f1654bc88ae6916722424c565976ef88"
  },
  "sharedSecret": "I08U3nLP/qtTwHjR5yGmle3cA5kpNiQDHhLwtzJfvl0=",
  "theirStealthAddresses": [
    {
      "public": {
        "x": "0x1f423980ffbd537b90d8581afb34a91d870d13455928737ee5c60292b7af4752",
        "y": "0x938aab5f627113ef524ce77458e89d0c611ce1837b34c9308a93d35622759de"
      },
      "nonce": 0
    }
  ],
  "myStealthAddresses": [
    {
      "public": {
        "x": "0x15472269a2b1c629d136a10732f7422ac73213d7712198be6e4f0060fa48e964",
        "y": "0x129e34473d8c37f6206a9bc0be2f8fd1bb4999902d10ea35d3d50fb775e9f5e4"
      },
      "nonce": 0,
      "private": 14674696873446288336876702922514381849233625593912482006822333717333299928756
    }
  ]
}
```

The other side can derive their stealth addresses using the following command:

```
$ orbital stealth -s 0x282e1c333c40fb8ffcb97a66d615af2a0c4077c5dc5dd50075cf6abf3e9e9f55 -x 0x27e073fe3485b7ab97de5813342c3ce3dd19eba467a6b5f61b092813e67325e2 -y 0x639bbbc72ef12bc7449c08a4ac2d7d7a9ed0484e4d5d6ee4f4b6fa5844043c2
```

Note that the public keys calculated on either side will be the same, but neither side knows the others private key.

## Development

Dependencies are managed via [dep][1]. Dependencies are checked into this repository in the `vendor` folder. Documentation for managing dependencies is available in the [dep README][2].

The project follows standard Go conventions using `gofmt`. If you wish to contribute to the project please follow standard Go conventions. The CI server automatically runs these checks.

[1]: https://github.com/golang/dep
[2]: https://github.com/golang/dep/blob/master/README.md
[3]: https://github.com/clearmatics/mobius
