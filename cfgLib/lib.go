package cfgLib

//
// Copyright (C) Philip Schlump, 2015-2018.
//

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func ReadConfigFile(fn string, cfg interface{}) {
	// data, err := sizlib.ReadJSONDataWithComments(fn)
	data, err := ReadJSONDataWithComments(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading JSON data %s - %s\n", fn, err)
		return
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		fmt.Printf("Syntax error in JSON file, %s\n", err)
		return
	}

	return
}

// ===============================================================================================================================================
var ln *regexp.Regexp
var fi *regexp.Regexp
var cm *regexp.Regexp
var en *regexp.Regexp

func init() {
	ln = regexp.MustCompile("__LINE__")
	fi = regexp.MustCompile("__FILE__")
	en = regexp.MustCompile("__ENV__:[a-zA-Z][a-zA-Z_0-9]*")
	cm = regexp.MustCompile("////.*$")
}

// ReadJSONDataWithComments read in the file and handle __LINE__, __FILE__ and comments starting with 4 slashes.
func ReadJSONDataWithComments(path string) (file []byte, err error) {
	file, err = ioutil.ReadFile(path)
	if err != nil {
		// fmt.Printf("Error(10014): Error Reading/Opening %v, %s, Config File:%s\n", err, godebug.LF(), path)
		// fmt.Fprintf(os.Stderr, "%sError(10014): Error Reading/Opening %v, %s, Config File:%s%s\n", MiscLib.ColorRed, err, godebug.LF(), path, MiscLib.ColorReset)
		return
	}

	data := strings.Replace(string(file), "\t", " ", -1)
	lines := strings.Split(data, "\n")
	//ln := regexp.MustCompile("__LINE__")
	//fi := regexp.MustCompile("__FILE__")
	//cm := regexp.MustCompile("//.*$")
	for lineNo, aLine := range lines {
		aLine = ln.ReplaceAllString(aLine, fmt.Sprintf("%d", lineNo+1))
		aLine = fi.ReplaceAllString(aLine, path)
		aLine = cm.ReplaceAllString(aLine, "")
		for en.MatchString(aLine) { // pick up and replace environment variables - put passwords in env not in config files
			// fmt.Printf("matched __ENV__:Name, %s\n", godebug.LF())
			ss := en.FindAllString(aLine, 1)
			// fmt.Printf("ss = %s\n", ss)
			s := ss[0] // the matched, no need to check array because inside MatchString already
			// fmt.Printf("s(raw) = %s\n", s)
			s = s[8:] // remove __ENV__:
			// fmt.Printf("env name = [%s]\n", s)
			v := os.Getenv(s)
			// fmt.Printf("v = [%s]\n", v)
			if v == "" {
				fmt.Fprintf(os.Stderr, "Fatal: Invalid environment variable setting: %s - returned empty string - not allowed\n", s)
				os.Exit(1)
			}
			aLine = en.ReplaceAllString(aLine, v)
			// fmt.Printf("final line = [%s]\n", aLine)
		}
		lines[lineNo] = aLine
	}
	file = []byte(strings.Join(lines, "\n"))

	//	fmt.Printf("%s %s\n", godebug.LF(2), godebug.LF(3))
	//	fmt.Printf("After fix of Remove Comments and Set Line Numbers: Results %s >%s<\n", godebug.LF(), file)

	return file, nil
}
