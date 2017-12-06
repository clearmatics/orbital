# Orbital

Orbital is a command-line tool to generate off-chain data required by [Möbius][3], a smart contract that offers trustless tumbling for transaction privacy.

## Prerequisites

A version of Go >= 1.8 is required. The [dep][1] tool is used for dependency management. 

## Installation

    go get github.com/clearmatics/orbital

## Usage

When deployed a Möbius contract will emit a `RingMessage` that is an arbitrary hex encoded string. This message is signed to make withdrawals from the contract. 

Orbital can be used to generate all data needed to deposit and withdraw from a Möbius smart contract. Providing you have the `RingMessage` value data can be generated as follows. In this example the hex encoded string is given as `291a6780850827fcd8621...`. A ring size of 2 is generated.

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

### Stealth Addresses

Integration of Stealth Addresses into the Möbius contract is still in-progress, however they can be generated using the `mobius stealth` utility.

First generate a pair of key pairs:

```
$ orbital generate -n 2
{
  "pubkeys": [
    {
      "x": 53202990254242129984821116292342958982032538732251891028681618557466605103267,
      "y": 113350357918605008175431316781526832336003746709402773019184703717803218368823
    },
    {
      "x": 88424329917342262816711482350788738042527708970466043100387261704248855107491,
      "y": 59260658595316282581973744342597950668623851099072116196891936195050801350643
    }
  ],
  "privkeys": [
    22349825998821701797378606099676179909897898231961005770162217921273442111598,
    102648465399097654117920514091549396436770184270947888717131433613128480538084
  ]
}
```

Then derive a stealth address for the other party using the first secret key and the second public key, the JSON output displays the shared secret, your public key and their stealth addresses.

```
$ orbital stealth -s 22349825998821701797378606099676179909897898231961005770162217921273442111598 -x 88424329917342262816711482350788738042527708970466043100387261704248855107491 -y 59260658595316282581973744342597950668623851099072116196891936195050801350643
{
  "myPublic": {
    "x": 53202990254242129984821116292342958982032538732251891028681618557466605103267,
    "y": 113350357918605008175431316781526832336003746709402773019184703717803218368823
  },
  "theirPublic": {
    "x": 88424329917342262816711482350788738042527708970466043100387261704248855107491,
    "y": 59260658595316282581973744342597950668623851099072116196891936195050801350643
  },
  "sharedSecret": "5cLLuIzdNwPgDlgcp1l4QjcAKi3lBKuBLA8D3RtWLQI=",
  "theirStealthAddresses": [
    {
      "public": {
        "x": 69380850297621879107185275122279026017384660398818062589861370260278309399018,
        "y": 26871018971989864411633210899039205781865578748955351442237322594146409159475
      },
      "nonce": 0
    }
  ],
  "myStealthAddresses": [
    {
      "public": {
        "x": 74390537890405246099775105804528312454210249389630428389559455977430334747312,
        "y": 53655787374135976950002031543616192689347650013901897026704630974400993904672
      },
      "nonce": 0,
      "private": 111623309102215725136656670222624246002377157739441594457567804582028137434860
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
