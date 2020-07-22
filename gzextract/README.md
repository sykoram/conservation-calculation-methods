# gzextract

This program extracts contents defined by a regular expression from a gzip file and saves it to the output file.

By default, it is configured to extract a MSA column from a gzip file from lines in the format `<column_num>\t<score>\t<MSA>`.



## Setup

[Go](https://golang.org/) has to be installed [[Download]](https://golang.org/dl/); no external libraries are used.

This build command creates a binary file in the working directory:
```sh
go build gzextract.go 
```

## Usage

The flag `-h` or `-help` displays help.

Default configuration extracts a MSA column from a file with lines in the format `<column_num>\t<score>\t<MSA>` where `<column_num>` is an int, `<score>` is a floating point number and `<MSA>` is a sequence of upper letters and a dash character:

```sh
./gzextract -i INPUT_GZIP_FILE -o OUTPUT_FILE 
``` 

To use this program to more general purposes, the regular expression and the number of the capture group can be defined:
```sh
./gzextract -i INPUT_GZIP_FILE -o OUTPUT_FILE -r REGEX_STR -g CAPTURE_GROUP
```
where the `CAPTURE_GROUP` is an integer: 0 is the whole match, 1 is the first capture group in the regular expression.