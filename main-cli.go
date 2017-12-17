package main

/*
TODO

2. Implement as pre-process on tree
	combineAllAttrs   bool
	combineAttrs      map[string]map[string]bool // tagName . Attr or tagName . * - for all on tag
	combineAllContent bool
	combineContent    map[string]string

x. Read in config for (2) above
	-lcfg - left config
	-rcfg - right config
	attrs:
		tag-name attr-name -> nested-tag-name
		tag-name * -> nested-tag-name
		* attr-name -> nested-tag-name
		* * -> *

	content:
		tag-name -> attr-name
		* -> attr-name
		* -> *

	nosort:
		tag
		tag






-- done -- -- done -- -- done -- -- done -- -- done -- -- done -- -- done -- -- done -- -- done -- -- done -- -- done -- -- done --

*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pschlump/xml-diff/xmllib"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var Left = flag.String("left", "", "Left Input file name")    // 0
var LeftCfg = flag.String("lcfg", "", "Left Config File")     // 1
var LeftOut = flag.String("lo", "", "Left Output file name")  // 2
var Right = flag.String("right", "", "Right Input file name") // 3
var RightCfg = flag.String("rcfg", "", "Right Config File")   // 4
var RightOut = flag.String("ro", "", "Left Output file name") // 5
// var Debug = flag.String("debug", "", "Debug/Trace Flags")     // 4
// var Output = flag.String("output", "", "Send output to file")     // 4	// and below -o
var ByLine = flag.Bool("byLine", false, "Show changes by entire line") // 5

func init() {
	flag.StringVar(Left, "l", "", "Input file name ")           // 0
	flag.StringVar(Right, "r", "", "Input file to compare to ") // 2
}

func main() {

	flag.Parse()

	fns := flag.Args()
	if len(fns) != 0 {
		fmt.Fprintf(os.Stderr, "Usage: xmlProc -i input-file.xml -o output.file -keepArray\n")
		os.Exit(1)
	}

	cleanXmlLeft := XmlClean(*Left, *LeftOut)
	cleanXmlRight := XmlClean(*Right, *RightOut)

	dmp := diffmatchpatch.New()
	if !*ByLine {
		diffs := dmp.DiffMain(cleanXmlLeft, cleanXmlRight, false)
		fmt.Println(dmp.DiffPrettyText(diffs))
	} else {
		wSrc, wDst, warray := dmp.DiffLinesToChars(cleanXmlLeft, cleanXmlRight)
		diffs := dmp.DiffMain(wSrc, wDst, false)
		diffs = dmp.DiffCharsToLines(diffs, warray)
		fmt.Println(dmp.DiffPrettyText(diffs))
	}
}

func XmlClean(fn, ofn string) string {

	// If file not exits - then fail
	if !xmllib.Exists(fn) {
		fmt.Fprintf(os.Stderr, "Missing input file %s\n", fn)
		os.Exit(1)
	}

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	xml := strings.NewReader(string(buf))
	cleanXmlLeft, err := xmllib.ConvertXML(xml) // returns a []byte, err
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if ofn != "" {
		ofp, err := xmllib.Fopen(ofn, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open output file %s error %s\n", ofn, err)
		} else {
			fmt.Fprintf(ofp, "%s", cleanXmlLeft)
		}
		ofp.Close()
	}

	return cleanXmlLeft.String()
}

const db1 = false

/* vim: set noai ts=4 sw=4: */
