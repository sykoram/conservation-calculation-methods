/*
This program extracts MSA columns from a file (may be gzip),
calculates a conservation score using a chosen method,
and writes the result into a output file.
 */

package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
	filepath2 "path/filepath"
	"regexp"
	"strings"
)

const replaceNegScore0 = true

var help bool
var ifile string
var ofile string

var methodId string
var validLineRegex = "^(\\d+)\\t(-?\\d+(?:\\.\\d+)?)\\t([A-Z\\-]+)$" // matches line <colNum>\t<score>\t<msaCol>
var msaCaptureGroup = 3  // which group should be captured (first is 1; 0 is the whole match)

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")

	flag.StringVar(&ifile, "i", "", "input file path (can be .gz) (required)")
	flag.StringVar(&ofile, "o", "", "output file path (gzipped if with .gz extension) (required)")
	flag.StringVar(&methodId, "m", "", fmt.Sprintf("conservation calculation method (required) %s", GetMethodNames()))

	flag.IntVar(&WindowSize, "w", 0, "window size (number of residues on each side included in the window)")
	flag.StringVar(&validLineRegex, "r", validLineRegex, "regex that matches a valid line where a MSA column is")
	flag.IntVar(&msaCaptureGroup, "g", msaCaptureGroup, "capture group of the valid-line regex -r (0 is the whole match; 1st group is 1)")
}

func main() {
	flag.Parse()
	handleHelp()
	checkCmd()

	// prepare the input file for reading
	fi, err := os.Open(ifile)
	fatalIfErr(err)
	defer fatalIfErrF(fi.Close) // defer executes the statement when the main() function ends

	var scanner *bufio.Scanner
	if filepath2.Ext(ifile) == ".gz" {
		gzr, err := gzip.NewReader(fi)
		fatalIfErr(err)
		defer fatalIfErrF(gzr.Close)
		scanner = bufio.NewScanner(gzr)
	} else {
		scanner = bufio.NewScanner(fi)
	}

	// prepare for writing
	fo, err := createFile(ofile)
	fatalIfErr(err)
	defer fatalIfErrF(fo.Close)

	var writer *bufio.Writer
	if filepath2.Ext(ofile) == ".gz" {
		gzw := gzip.NewWriter(fo)
		defer fatalIfErrF(gzw.Close)
		writer = bufio.NewWriter(gzw)
	} else {
		writer = bufio.NewWriter(fo)
	}

	// get method
	method, ok := Methods[methodId]
	if !ok {
		log.Fatal("Unknown method identifier.")
	}

	// read input file and extract all MSA columns
	i := 0
	msaCols := make([]MsaColumn, 0)
	for scanner.Scan() {
		col, err := extractMsaColumn(scanner.Text(), validLineRegex, msaCaptureGroup)
		fatalIfErr(err)
		if col == "" {
			continue
		}
		msaCols = append(msaCols, col)
		i++
	}

	// calculate score, save to output file
	seqWeights := GetSequenceWeights(msaCols)
	scores := make([]float64, len(msaCols))
	for i, col := range msaCols {
		scores[i] = -1000.0
		if GetGapRatio(col) <= MaxGapRatio {
			scores[i] = method(col, Blosum62SimMatrix, Blosum62BgDistr, seqWeights)
		}
	}

	if WindowSize > 0 {
		scores = WindowScores(scores, WindowSize, WindowLam)
	}

	for i, col := range msaCols {
		if replaceNegScore0 && scores[i] < 0 {
			scores[i] = 0
		}
		_, err = writer.WriteString(fmt.Sprintf("%d\t%.5f\t%s\n", i, scores[i], col)) // the format of the output line is <colNum>\t<score>\t<msaCol>
		fatalIfErr(err)
	}

	fatalIfErr(writer.Flush())
	fatalIfErrF(scanner.Err)
}

/*
Checks the flags and arguments. If something is not right, fatal error is produced.
-i, -o and -m flags are required, any additional arguments are forbidden.
*/
func checkCmd() {
	if ifile == "" {
		log.Println("[ERROR] The input file path is required: -i path/to/file")
		defer os.Exit(1)
	}

	if ofile == "" {
		log.Println("[ERROR] The output file path is required: -o path/to/file")
		defer os.Exit(1)
	}

	if methodId == "" {
		log.Println("[ERROR] The method has to be specified: -m METHOD")
		defer os.Exit(1)
	}

	if flag.NArg() > 0 {
		log.Println("[ERROR] Unknown arguments: " + strings.Join(flag.Args(), " "))
		defer os.Exit(1)
	}
}

/*
Handles help flag -h. If the help is requested, prints program description and flags.
*/
func handleHelp() {
	if help {
		fmt.Print(`conservscore extracts MSA columns from an input file (may be gzip),
calculates conservation score for each column using a chosen method,
and writes the result into a output file.

By default, the MSA column is extracted from a line in format
<colNum>\t<score>\t<msaCol>
where <colNum> is an integer, <score> is a floating point number
and <msaCol> is a sequence of letters or dashes.

But the user can use this program with different formats by defining
a custom regular expression and a capture group.

Usage:
`)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

/*
Creates or truncates file using os.Create(). Creates parent directories if required.
Error is returned if present.
The file has to be closed manually!
*/
func createFile(filepath string) (*os.File, error) {
	err := os.MkdirAll(filepath2.Dir(filepath), 666)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

/*
Extracts a MSA column from the line based on the regex (and capture group).
If there is no match, empty string is returned.
*/
func extractMsaColumn(line, validLineRegexStr string, msaCaptureGroup int) (MsaColumn, error) {
	regex, err := regexp.Compile(validLineRegexStr)
	if err != nil {
		return "", err
	}
	matches := regex.FindStringSubmatch(line)

	if len(matches) >= msaCaptureGroup+ 1 {
		return MsaColumn(matches[msaCaptureGroup]), nil
	}
	return "", nil
}

/*
Calls log.Fatal(err) if there is an error.
 */
func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
Calls the function and then fatalIfErr if there is an error.
(Used with defer)
 */
func fatalIfErrF(f func() error) {
	fatalIfErr(f())
}

