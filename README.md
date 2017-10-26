# Orbital

Orbital is a command-line tool to generate off-chain data required by [Möbius][3], a smart contract that offers trustless tumbling for transaction privacy.

## Installation

    go get github.com/clearmatics/orbital

## Usage

When deployed a Möbius contract will emit a `Message` that is an arbitrary hex encoded string. This message is signed to make withdrawals from the contract. 

Orbital can be used to generate all data needed to deposit and withdraw from a Möbius smart contract. Providing you have the `Message` value data can be generated as follows. In this example the hex encoded string is given as `666f6f62617262617a`. A ring size of 2 is generated.

    orbital -geninputs 2 666f6f62617262617a

This generates JSON containing complete data required to deposit and withdraw. If you are just evaluating Möbius this is all you need to deposit into a contract and then make withdrawals once the ring is full. 

``` JSON
{
  "ring": [
    {
      "x": "35470038841138887629374608865734874129902573308249245396015702392600471414928",
      "y": "18469706635113601059728665984946917241028154186083727264552207721727178628984"
    },
    {
      "x": "95316358434591987899377038267177201978376874068933785035289629333730618475640",
      "y": "77512950697498367848129712442116806909438279661773376431344468660653131210719"
    }
  ],
  "signatures": [
    {
      "tau": {
        "x": "5553827006078961602311820336891916902356846643372753101575873865552224584041",
        "y": "71485606063063329573362498629693429103219908348180489741983970039582563279828"
      },
      "ctlist": [
        "106869368489868104981767483990118849056954574146050847794438537149023534777936",
        "80060180827376907910188627680811694073152642244719696613779648949475310945440",
        "76632194641985177784785131818097694279922836679214727849391048407032094385287",
        "84576045742035030875699019658185630958421078950082272987171632827975652546317"
      ]
    },
    {
      "tau": {
        "x": "21571541925636444597303138636855820121400431633514116849476600225379338470925",
        "y": "103993752232816935193321298139040608100940624085983889467161354138852167761468"
      },
      "ctlist": [
        "79694312270207959348401881209687603232755996067468568085249462727176736059759",
        "73139344268323607166737892175878352019314998602996403987316276511785289908333",
        "36831617827612302911886522178304054730468273001705980549338498547586736925474",
        "45975040705170998173120548901509912630059801394260160165656428067929922137782"
      ]
    }
  ]
}
```

Beyond evaluation it is likely that you will want to first generate public/private keypairs and distribute this to allow signatures to be created separately. To create a public/private keypairs for a ring size of `2` use 
    
    oribtal -genkeys -2

To print signatures and public keys ring from keys file

    orbital -signature keys.json HexEncodedString

To verify signatures

    orbital -verify signature.json HexEncodedString

Examples 
```
    $ orbital -genkeys 4 > keys.json
    $ orbital -signature keys.json 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc > ringSignature.json
    $ orbital -verify ringSignature.json 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc
    Signatures verified
```
or
```
    $ orbital -geninputs 4 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc > ringSignature.json
    $ orbital -verify ringSignature.json 50b44f86159783db5092ebe77fb4b9cc29e445e54db17f0e8d2bed4eb63126fc
    Signatures verified
```

## Development

Dependencies are managed via [dep][1]. Dependencies are checked into this repository in the `vendor` folder. Documentation for managing dependencies is available on the [dep README][2].

The project follows standard go conventions using `gofmt`. If you wish to contribute to the project please follow standard Go conventions. The CI server automatically runs these checks.

[1]: https://github.com/golang/dep
[2]: https://github.com/golang/dep/blob/master/README.md
[3]: https://gitlab.clearmatics.com/oss/mobius
