# Orbital

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
  "message": "KRpngIUIJ/zYYh0OVHE0ODEQm8FBQuwQFSewSLs9F5Q=",
  "ring": [
    {
      "x": 15522664089955551659852289328166700143127393484899632195419246838083707083092,
      "y": 139725407559755444649292554266801350425737376800472574275109410266714217628
    },
    {
      "x": 15769601697165987408284147780311087633600327546404112082093068531966136785797,
      "y": 574193457253025311982882197162167841837660560297023031472921816123174737979
    }
  ],
  "signatures": [
    {
      "tau": {
        "x": 10156463942228021228545876819561826038791281705421256076529438066781624852871,
        "y": 10521028020321534223595269087268496225446461286869992881505824752024716692663
      },
      "ctlist": [
        21368655831583064736758400887160467297479232369129086125811187621405591279735,
        12597762440685916653847888008183831634464726993537178408507664656131668212830,
        6990298203776631120661116765313744556075488813405658404819872199156195642387,
        14268213418230864809482216147768002586757463386138671647979779695664999133350
      ]
    },
    {
      "tau": {
        "x": 18350120056062457577455309792248846176018677319782425143665947389092379940018,
        "y": 17410042836747417050905448604207541669732600995872545566239761559163293972399
      },
      "ctlist": [
        5234961786120978809275036908808283773251072203711221622865855236314635213488,
        13528210618504793075305289291975247736253154818401760963070837217369448096451,
        2788192854693864373131102131755172014022189917825686286639429882366651311778,
        17072584727078276773282114103923714960058637541295589940252538659929574197139
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
$ orbital stealth -s 0x10dfcf223f6dd92163004a34a5dd509e7ad57cc2199eebb0ac2494e6e87a8246 -x 0x304e70ae7fe4d42ba8a3cbb25d2442b09af08522d93eee77110e502e0bf21d3e -y 0x1a913cdcb4296b564e2cfe1d747e9fa5aa770c5c5b990324c69845be67c98a9b

{
  "myPublic": {
    "x": "0xf6f50dec186fa8724c2b7597c3a51abaef04f42dd30628521e41f4379bf1642",
    "y": "0x221935f311652fe5edde5a0098b9d758edc883faef61acfd4cfc7711e010e75"
  },
  "theirPublic": {
    "x": "0x304e70ae7fe4d42ba8a3cbb25d2442b09af08522d93eee77110e502e0bf21d3e",
    "y": "0x1a913cdcb4296b564e2cfe1d747e9fa5aa770c5c5b990324c69845be67c98a9b"
  },
  "sharedSecret": "HnoXTU8p8WGJPtlcIUb9tLCjurAqnJkwNiVMcaUmVhU=",
  "theirStealthAddresses": [
    {
      "public": {
        "x": "0x1e82f9b171c24536f06061a207112af54db5232dc2331f5432b676339a3bbe8e",
        "y": "0x133e9019c4c367505ca5b791dc92244a9513b5b9c96f183729c1b07c1946d25f"
      },
      "nonce": 0
    }
  ],
  "myStealthAddresses": [
    {
      "public": {
        "x": "0x2ecc04e14f8d74b1a08486d0aa4b8be682ca92f4082b2195bd95380312b9ae46",
        "y": "0x17d954284508272403d121f4b9392456aec7938ea3e70844f6db79ea43b89cbe"
      },
      "nonce": 0,
      "private": 20885560622514516919971043866714208233025685278744852588179698445235938879992
    }
  ]
}
```

The other side can derive their stealth addresses using the following command:

```
$ orbital stealth -s 102648465399097654117920514091549396436770184270947888717131433613128480538084 -x 53202990254242129984821116292342958982032538732251891028681618557466605103267 -y 113350357918605008175431316781526832336003746709402773019184703717803218368823
```

Note that the public keys calculated on either side will be the same, but neither side knows the others private key.

## Development

Dependencies are managed via [dep][1]. Dependencies are checked into this repository in the `vendor` folder. Documentation for managing dependencies is available in the [dep README][2].

The project follows standard Go conventions using `gofmt`. If you wish to contribute to the project please follow standard Go conventions. The CI server automatically runs these checks.

[1]: https://github.com/golang/dep
[2]: https://github.com/golang/dep/blob/master/README.md
[3]: https://github.com/clearmatics/mobius
