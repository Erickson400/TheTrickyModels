package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func MxfToJson(filename string, destination string) {

	// Read the File
	if !strings.HasSuffix(filename, ".mxf") {
		filename += ".mxf"
	}

	f, err := os.Open(filename)
	if err != nil {
		panic("Failed to read .mxf file")
	}

	// Turn the file's data to a struct
	theFile := FileToStruct(f)

	// Create the json file
	empJSON, err := json.MarshalIndent(theFile, "", "  ")
	if err != nil {
		panic(err)
	}
	o, _ := os.Create("out.json")
	fmt.Fprintln(o, string(empJSON))
}
