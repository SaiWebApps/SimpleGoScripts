# SimpleGoScripts

This package contains a collection of golang scripts that I constructed while familiarizing myself with the language.

## [hellobot.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/hellobot/hellobot.go)

Say "Hello" to the specified subject.
```
go run hellobot/hellobot.go [-i] [-to SUBJECT]
-i: Enter interactive mode to specify subject.
-to: Explicitly/directly specify the subject.
```
If neither flag is specified, then this bot shall say "Hello" to "World."
On the other hand, if both flags are specified, then **-i** shall take precedence, and **-to** shall be ignored.

## [pascalstriangle.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/pascalstriangle/pascalstriangle.go)

Generate Pascal's Triangle for the specified number of levels.

First, run:
```
go install simplegoscripts/intmatrix
```
to compile the custom intmatrix package within this project, which is a dependency for pascalstriangle.go.

Afterwards, invoke pascalstriangle/pascalstriangle.go as follows:
```
go run pascalstriangle/pascalstriangle.go [-n NUM_LEVELS] 
-n: Positive integer representing the number of levels of the target Pascal's Triangle. Default value is 1.
```

## [primenumbers.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/primenumbers/primenumbers.go)

Compute prime numbers up to the specified integer ceiling through the [Sieve-of-Eratosthenes](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes) heuristic.
```
go run primenumbers/primenumbers.go [-n CEILING]
-n: Generate prime numbers up to and including this non-negative integer. Default value is 1.
```

## [wordcounter.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/wordcounter/wordcounter.go)

Generate a word histogram for the specified string, file, and/or URL.
```
go run wordcounter/wordcounter.go [-u URLS] [-f FILE_PATHS] [-s TEXT]
-u: 1 or more URLs separated by whitespace.
-f: 1 or more file-paths separated by whitesapce.
-s: Directly specify the text that needs to be analyzed.
```
For any combination of these flags, this script shall:
1. collate the source texts,
2. remove non-word characters (\W),
3. identify the words, which are presumably separated by whitespaces (\s+),
4. and calculate the number of times that each word occurs in total across all of the source texts.

## [currencyconverter.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/currencyconverter/currencyconverter.go)

Perform currency conversions via command-line by leveraging Google Finance's Currency Converter Tool.
```
go run currencyconverter/currencyconverter.go [-a AMOUNT] [-f FROM_DENOM] [-t TO_DENOM]
-a: Positive amount of money in FROM_DENOM to convert to TO_DENOM. Default is 1.
-f: Denomination of AMOUNT; the denomination that we are converting from.
-to: Denomination that AMOUNT, specified in FROM_DENOM, is being converted to.
-l: List valid short and long names for FROM_DENOM and TO_DENOM denominations.
```
For example
```
go run currencyconverter/currencyconverter.go -f USD -to EUR
```
shall convert 1 US Dollar to Euros.

When **-l** is specified, it shall take precedence over all other parameters, and the script, in turn, shall display an informational message on acceptable values for **-f** and **-t**.

Said informational message shall consist of lines of colon-delimited strings, wherein the left-hand-side of each string is a valid short denomination name, and the right-hand-side is a valid long denomination name. Values specified to **-f** and **-t** shall either exactly match one of the aforementioned short names OR shall be contained within one of the long names.

## [sayhellowebapp.go](https://github.com/SaiWebApps/SimpleGoScripts/blob/master/sayhellowebapp/sayhellowebapp.go)

"Hello World" Golang Web Application
1. go run sayhello/sayhellowebapp.go
2. Open up a web browser, and navigate to http://localhost:8080/. The page should display the default "Hello World" message.
3. To specify a different subject to say hello to, simply add the desired subject's name after the final "/" in http://localhost:8080/. For example, "http://localhost:8080/me" would result in the message "Hello me!".