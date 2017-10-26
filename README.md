# Orbital

Orbital is a command-line tool to generate off-chain data required by [Möbius][3], a smart contract that offers trustless tumbling for transaction privacy.

## Installation

    go get github.com/clearmatics/orbital

## Usage

When deployed a Möbius contract will emit a `Message` that is an arbitrary hex encoded string. This message is signed to make withdrawals from the contract. 

Orbital can be used to generate all data needed to deposit and withdraw from a Möbius smart contract. Providing you have the `Message` value data can be generated as follows. In this example the hex encoded string is given as `666f6f62617262617a`. A ring size of 2 is generated.

    orbital inputs -n 2 -m 666f6f62617262617a

This generates JSON containing complete data required to deposit and withdraw. If you are just evaluating Möbius this is all you need to deposit into a contract and then make withdrawals once the ring is full. 

``` JSON
{
  "pubkeys": [
  {
    "x": 70742237615164982931155265691258765833107687230355458841421693916924953784687,
      "y": 30840380996293202136210045282507572367400788071239347648907713335359761829138
  },
  {
    "x": 27087078830110323521665066978195701528751557301745005184936110685042440855975,
    "y": 59727407121757349990184094045181827269744685570035889585265768058760849316377
  }
  ],
  "signatures": [
  {
    "tau": {
      "x": 115291306959110534762331116217359678632223049313671872799618131253150795715011,
      "y": 76080870794898937526792396323808041305747583416466311515213280049789496276229
    },
    "ctlist": [
      35677189825894557156716153963502736039665348954321031305574694716619151189549,
    65174561444261965848365705550964568453721181219861371050511999790449291429430,
    68918314092231193430341549175302250051979507143988817957533282992996303297168,
    112489839489816450458560166033668548097711559543822342428171065661517939963774
    ]
  },
  {
    "tau": {
      "x": 115006347313112976079885905821054725945159504615599662386088608699176267829093,
      "y": 69044543452700772992241492155130578191353062867454014235414028518071012328630
    },
    "ctlist": [
      36089834007966644604238667848008859557570133615472349059737948082862755319880,
    112132419790538882202963829436143925417567529362411156317555836886017363595398,
    74804663813747299689957267313237784270831420442219423802333920445997909078632,
    84943970980345001069419227022334735117723462889367512917231257466664229790991
    ]
  }
  ]
}
```

Beyond evaluation it is likely that you will want to first generate public/private keypairs and distribute this to allow signatures to be created separately. To create a public/private keypairs for a ring size of `2` use 
    
    orbital generate -n 2

To verify signatures pass a file with public keys and signatures to the `verify` subcommand.

    orbital verify -f signatures.json -m 666f6f62617262617a

Examples 
```
    $ orbital generate -n 4 > keys.json
    $ orbital verify -f ringSignature.json -m 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc
    Signatures verified
```
or
```
    $ orbital inputs -n 4 -m 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc > ringSignature.json
    $ orbital verify -f ringSignature.json -m 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc
    Signatures verified
```

## Development

Dependencies are managed via [dep][1]. Dependencies are checked into this repository in the `vendor` folder. Documentation for managing dependencies is available in the [dep README][2].

The project follows standard Go conventions using `gofmt`. If you wish to contribute to the project please follow standard Go conventions. The CI server automatically runs these checks.

[1]: https://github.com/golang/dep
[2]: https://github.com/golang/dep/blob/master/README.md
[3]: https://gitlab.clearmatics.com/oss/mobius
