# conservscore

This program extracts MSA columns from an input file, calculates conservation score for each column using a chosen method, and writes the result into a output file.


## Setup

[Go](https://golang.org/) has to be installed [[Download]](https://golang.org/dl/) \
(no external libraries are used)

This build command creates a binary file in the working directory:
```sh
go build conservscore.go methods.go
```


## Usage

Use flag `-h` or `-help` to display help.

To calculate conservation scores, specify input file, output file and method.

```sh
./conservscore -i INPUT_FILE -o OUTPUT_FILE -m METHOD
```

If the input file name has a `.gz` extension, it is automatically decompressed. Similarly, the output file is compressed. 

Currently available methods are: zero

By default, the MSA column is extracted from a line in format `<colNum>\t<score>\t<msaCol>` where `<colNum>` is an integer, `<score>` is a floating point number and `<msaCol>` is a sequence of letters or dashes.

To use this program with different formats, a custom regular expression and a capture group can be defined:
Flag `-r` sets a regular expression. \
Flag `-g` sets a capture group. It is an integer: 0 is the whole match, 1 is the first capture group of the regular expression.
