package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	//-------------------------------------------
	// xmxf
	// arg1: [mxf file] make sure the textures are on the same dir if its not a board
	// arg2: [gltf destination path]

	// go run . ./resources/board.mxf  ./
	//-------------------------------------------

	// Help
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "h" || os.Args[1] == "-help" || os.Args[1] == "help" {
			fmt.Printf(" Command: xmxf [mxf file] [gltf destination path] \n ")
			return
		}
	}
	if len(os.Args) < 3 {
		fmt.Println("You must give 2 arguments: xmxf [mxf file] [gltf destination path]")
		return
	}

	// Check if file/path exist
	_, err := os.Stat(filepath.ToSlash(os.Args[1]))
	if os.IsNotExist(err) {
		fmt.Println("the .mxf file directory is not valid: ", os.Args[1])
		return
	}
	_, err = os.Stat(filepath.ToSlash(os.Args[2]))
	if os.IsNotExist(err) {
		fmt.Println("gltf destination path does not exist", os.Args[2])
		return
	}

	// If all is gucci then procced.
	MxfToGltf(filepath.ToSlash(os.Args[1]), filepath.ToSlash(os.Args[2]))
}
