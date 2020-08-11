# conservscore

This program extracts MSA columns from an input file, calculates conservation score for each column using a chosen method, and writes the result into a output file.

**conservscore-dir script** runs conservscore for every file in a directory.


## Setup

[Go](https://golang.org/) has to be installed [[Download]](https://golang.org/dl/) \
(no external libraries are used)

This build command creates a binary file in the working directory:
```sh
go build -o conservscore ./conservscore.go ./methods.go
```

Since conservscore-dir script is a Bash script, you may need to be on Linux, have Git Bash or Windows Subsystem for Linux if you want to use it.


## Usage

Use flag `-h` or `-help` to display help (and also all supported methods).

To calculate conservation scores, specify input file, output file and method. Window is optional.

```sh
./conservscore -i INPUT_FILE -o OUTPUT_FILE -m METHOD [-w WINDOW]
```

If the input file name has a `.gz` extension, it is automatically decompressed. Similarly, the output file is compressed. 

By default, the MSA column is extracted from a line in format `<colNum>\t<score>\t<msaCol>` where `<colNum>` is an integer, `<score>` is a floating point number and `<msaCol>` is a sequence of letters or dashes.

To use this program with different formats, a custom regular expression and a capture group can be defined:
Flag `-r` sets a regular expression. \
Flag `-g` sets a capture group. It is an integer: 0 is the whole match, 1 is the first capture group of the regular expression.

**conservscore-dir script** supports only input directory, output directory and method, but it should be easy to modify it if you want to use the `-r` and `-g` flags.

```sh
./conservscore-dir.sh -i INPUT_DIR -o OUTPUT_DIR -m METHOD [-w WINDOW]
```
