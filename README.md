# SimpleGoScripts

This package contains a collection of golang scripts that I constructed while familiarizing myself with the language.

## [hellobot.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/hellobot.go)

Say "Hello" to the specified subject.
```
go run hellobot.go [-i] [-to SUBJECT]
-i: Enter interactive mode to specify subject.
-to: Explicitly/directly specify the subject.
```
If neither flag is specified, then this bot shall say "Hello" to "World."
On the other hand, if both flags are specified, then **-i** shall take precedence, and **-to** shall be ignored.

## [pascalstriangle.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/pascalstriangle.go)

Generate Pascal's Triangle for the specified number of levels.
```
go run pascalstriangle.go [-n NUM_LEVELS] 
-n: Positive integer representing the number of levels of the target Pascal's Triangle. Default value is 1.
```

## [primenumbers.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/primenumbers.go)

Compute prime numbers up to the specified integer ceiling through the [Sieve-of-Eratosthenes](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes) heuristic.
```
go run primenumbers.go [-n CEILING]
-n: Generate prime numbers up to and including this non-negative integer. Default value is 1.
```

## [wordcounter.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/wordcounter.go)

Generate a word histogram for the specified string, file, and/or URL.
```
go run wordcounter.go [-u URLS] [-f FILE_PATHS] [-s TEXT]
-u: 1 or more URLs separated by whitespace.
-f: 1 or more file-paths separated by whitesapce.
-s: Directly specify the text that needs to be analyzed.
```
For any combination of these flags, this script shall:
1. collate the source texts,
2. remove non-word characters (\W),
3. identify the words, which are presumably separated by whitespaces (\s+),
4. and calculate the number of times that each word occurs in total across all of the source texts.

## [currencyconverter.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/currencyconverter.go)

Perform currency conversions via command-line by leveraging Google Finance's Currency Converter Tool.
```
go run currencyconverter.go [-a AMOUNT] [-f FROM_DENOM] [-t TO_DENOM]
-a: Positive amount of money in FROM_DENOM to convert to TO_DENOM. Default is 1.
-f: Denomination of AMOUNT; the denomination that we are converting from.
-to: Denomination that AMOUNT, specified in FROM_DENOM, is being converted to.
-l: List all valid denominations.
```
For example
```
go run currencyconverter.go -f USD -to EUR
```
shall convert 1 US Dollar to Euros. For your edification, you can use **-l** to list all possible **-f**/**-t** denominations; in this scenario, all other flags shall be ignored.