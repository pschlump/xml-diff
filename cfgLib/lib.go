package cfgLib

//
// Copyright (C) Philip Schlump, 2015-2018.
//

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pschlump/Go-FTL/server/sizlib"
)

func ReadConfigFile(fn string, cfg interface{}) {
	// t := "r:listen"
	// cfg.ReplyTo = &t
	// data, err := ioutil.ReadFile(fn)
	data, err := sizlib.ReadJSONDataWithComments(fn)
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
