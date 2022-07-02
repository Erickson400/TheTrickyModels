package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	ssx()
}

func ssx() {
	fileHeader := FileHeader{
		Unknown:         0,
		ModelCount:      10,
		ModelListOffset: 12,
		ModelRootOffset: 3972,
		ModelHeaders: []ModelHeader{
			ModelHeader{
				//ModelName: [16]byte(),
			},
		},
	}

	PrettyPrint(fileHeader)
}

func PrettyPrint(structure FileHeader) {
	empJSON, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(empJSON))
}
