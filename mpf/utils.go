package main

import (
	"bytes"
	"fmt"
	"os"
)

func FindPattern(data []byte, tolerance int, pattern []byte) (position int, err error) {
	// Iterate through the whole file
	numsCorrect := 0
	for i := range data {
		for {
			if data[i+numsCorrect] == pattern[numsCorrect] {
				numsCorrect++
				if numsCorrect == len(pattern) {
					if tolerance <= 0 {
						return i, nil
					}
					tolerance--
					numsCorrect = 0
					break
				}
				continue
			}
			numsCorrect = 0
			break
		}
	}
	return 0, fmt.Errorf("could not find pattern")
}

func PrintLoc(buf *bytes.Reader) {
	m, _ := buf.Seek(0, os.SEEK_CUR)
	fmt.Println(m)
}
