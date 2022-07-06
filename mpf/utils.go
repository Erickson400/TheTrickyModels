package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ProcessMpfSystem(filename string, destination string) {
	b := MpfToBytes(filename)
	BytesToMpfStruct(b)

}

func MpfToBytes(filename string) []byte {
	if !strings.HasSuffix(filename, ".mpf") {
		filename += ".mpf"
	}

	data, err := os.ReadFile(filename) // the file is inside the local directory
	if err != nil {
		panic("Failed to read .mpf file")
	}
	return data
}

func BytesToMpfStruct(b []byte) (head FileHeader) {
	// Setup the buffer for reading in little endian
	buf := bytes.NewReader(b)
	Read := func(p any) { //Takes in a Pointer/Address
		err := binary.Read(buf, binary.LittleEndian, p)
		if err != nil {
			panic(err)
		}
	}

	// Start Reading data

	// FileHeader
	buf.Seek(4, os.SEEK_CUR)
	Read(&head.ModelCount)
	Read(&head.ModelHeaderListOffset)
	Read(&head.ModelRootOffset)

	// ModelHeaders

	//fmt.Printf("%+v", head)
	PrettyPrint(head)

	return
}

func PrettyPrint(structure FileHeader) {
	empJSON, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(empJSON))
}
