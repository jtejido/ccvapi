# ccvapi
Credit Card Number verifier API

This is an Restful API that verifies credit card numbers (Issuers and Checksum) via Luhn Algorithm.


## Card Types

All Card Types are located in the included card_types.json:

```
[
  {
    "name": "Visa",
    "patterns": [
      4
    ],
    "lengths": [16, 18, 19]
  },
  {
    "name": "MasterCard",
    "patterns": [
      [51, 55],
      [2221, 2229],
      [223, 229],
      [23, 26],
      [270, 271],
      2720
    ],
    "lengths": [16]
  }
]
```

Adding a card Issuer can be done only via this file. The **"patterns"** field is an array of homogenous type that can either be an int (e.g. Visa has one pattern to follow, and that is all numbers would start in 4) or an array of [min, max], (e.g., MasterCard uses multiple of these ranges).

The **lengths** field is where the usual length of digits that is used by the specific Issuer on their Customer's Account Numbers. (This is not range, as you can see well that Visa uses 16, 18 or 19 digits long for their Account Numbers, skipping 17)

## HTTP Config

HTTP configuration is in the included config.toml file. At the very least, you would set the path for both **error.log** and **access.log** in here, the **port** to use, and where the **card types json** file would be fetched from.

All default values are filled in.

## Command Line

Alternatively, you can also set them through commandline flags. All you have to do is run it as follows:

```
go get github.com/jtejido/ccvapi
go build
./ccvapi -h
  -access-log string
        Location of the logfile. (default "error.log")
  -card-path string
        Location of the card types json file. (default "card_types.json")
  -error-log string
        Location of the logfile. (default "access.log")
  -host int
        The local port to listen to. (default 8080)
```

## Sending Request

Once ran, it will listen on the 8080 port of the localhost, (http://localhost:8080/card/api/verify).

Verifying a PAN can be done via POST:

```
curl --header "Content-Type: application/json" --request POST --data '{ "PAN": "2222400050000009" }' http://localhost:8080/card/api/verify
```

And will send the response in json format.
```
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   110  100    81  100    29   4764   1705 --:--:-- --:--:-- --:--:--

{"Valid":true,"Issuer":"MasterCard","Error":"","PatternMatch":4,"LengthMatch":16}
```

## Built-in Card Issuers

See this [json file](https://github.com/jtejido/ccvapi/card_types.json)

| Brand              |
|--------------------|
| `Visa`             |
| `Mastercard`       |
| `American Express` |
| `Diners Club`      |
| `Discover`         |
| `JCB`              |
| `UnionPay`         |
| `Maestro`          |
| `Mir`              |
| `Elo`              |
| `Hiper`            |
| `Hipercard`        |

