/*
This program extracts pure MSA columns (by default) from the input gzip file and saves them to the output (not zipped) file. (Everything one line at a time)
But it can be used simply for any extraction from a gzip file by a regular expression.

The default format of the input file: <column>\\t<score>\\t<MSA> - only the <MSA> is extracted by default.
All empty or # commented lines are ignored (removed) by default valid-line regex.
 */

package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var help bool
var ifile string
var ofile string

var validLineRegex = "^(\\d+)\\t(-?\\d+(?:\\.\\d+)?)\\t([A-Z\\-]+)$" // matches line [columnNum]\t[score]\t[MSA]
var msaCaptureGroup = 3  // which group should be captured (first is 1; 0 is the whole match)

func init() {
	flag.StringVar(&ifile, "i", "", "input gzip file path (required)")
	flag.StringVar(&ofile, "o", "", "output file path (required)")

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")

	flag.StringVar(&validLineRegex, "r", validLineRegex, "regex of a valid line where MSA should be extracted from")
	flag.IntVar(&msaCaptureGroup, "g", msaCaptureGroup, "capture group number of the valid-line regex -r (0 is the whole match; 1st group is 1)")
}

func main() {
	flag.Parse()
	handleHelp()
	checkCmd()

	// prepare for reading
	fi := openFile(ifile)
	defer closeFile(fi) // calls closeFile(fi) when the main() function ends
	reader := getGzipReader(fi)
	defer closeGzipReader(reader)
	scanner := bufio.NewScanner(reader)

	// prepare for writing
	fo := createFile(ofile)
	defer closeFile(fo)
	writer := bufio.NewWriter(fo)

	// main loop: read a line, write extracted MSA column to the output file
	for scanner.Scan() {
		line := scanner.Text()
		msa := extractMsaColumn(line, validLineRegex, msaCaptureGroup)
		if msa != "" {
			writeMsaColumn(writer, msa)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

/*
Checks the flags and arguments. If something is not right, fatal error is produced.
-i and -o flags are required, any additional arguments are forbidden.
 */
func checkCmd() {
	fatal := false

	if ifile == "" {
		fmt.Println("[ERROR] The input file path is required: -i path/to/file")
		fatal = true
	}

	if ofile == "" {
		fmt.Println("[ERROR] The output file path is required: -o path/to/file")
		fatal = true
	}

	if flag.NArg() > 0 {
		fmt.Println("[ERROR] Unknown arguments: " + strings.Join(flag.Args(), " "))
		fatal = true
	}

	if fatal {
		fmt.Println("See -h or -help for description and usage")
		os.Exit(1)
	}
}

/*
Handles help flag -h. If the help is requested, prints program description and flags.
 */
func handleHelp() {
	if help {
		fmt.Println("This program extracts contents defined by a regular expression from a gzip file and saves it to the output file.")
		fmt.Println("By default, it is configured to extract a MSA column from a gzip file from lines in the format <column_num>\\t<score>\\t<MSA>")
		fmt.Println("")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

/*
Opens file for reading using os.Open(). Any produced error is handled. The file has to be closed manually!
 */
func openFile(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

/*
Closes file using .Close(). Handles error.
 */
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

/*
Returns new gzip reader for given file. Any produced error is handled. The reader has to be closed manually!
 */
func getGzipReader(file *os.File) *gzip.Reader {
	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	return reader
}

/*
Closes gzip.Reader using .Close(). Handles error.
 */
func closeGzipReader(r *gzip.Reader) {
	err := r.Close()
	if err != nil {
		log.Fatal(err)
	}
}

/*
Creates or truncates file using os.Create(). Potential error is handled. The file has to be closed manually!
 */
func createFile(filepath string) *os.File {
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

/*
Extracts MSA column/string from the line (MSA column matches the regex).
If there is no match, empty string is returned.
 */
func extractMsaColumn(line, validLineRegexStr string, msaCaptureGroup int) string {
	regex, err := regexp.Compile(validLineRegexStr)
	if err != nil {
		log.Fatal(err)
	}
	matches := regex.FindStringSubmatch(line)

	if len(matches) >= msaCaptureGroup+ 1 {
		return matches[msaCaptureGroup]
	}
	return ""
}

/*
Writes the MSA column. Adds \n. Flushes the writer. Handles errors.
 */
func writeMsaColumn(w *bufio.Writer, msa string) {
	_, err := w.WriteString(msa)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.WriteRune('\n')
	if err != nil {
		log.Fatal(err)
	}

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
