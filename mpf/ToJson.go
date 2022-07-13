package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func MpfToJson(filename string, destination string) {

	// Read the File
	if !strings.HasSuffix(filename, ".mpf") {
		filename += ".mpf"
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		panic("Failed to read .mpf file")
	}

	// Turn the file's data to a struct
	theFile := BytesToStruct(b)

	// Create the json file
	empJSON, err := json.MarshalIndent(theFile, "", "  ")
	if err != nil {
		panic(err)
	}
	f, _ := os.Create("out.json")
	fmt.Fprintln(f, string(empJSON))
}

func BytesToStruct(b []byte) (file TheFile) {
	// Setup the buffer for reading in little endian
	buf := bytes.NewReader(b)
	Read := func(p any) { //Takes in a Pointer/Address
		err := binary.Read(buf, binary.LittleEndian, p)
		if err != nil {
			panic(err)
		}
	}

	// Start Reading TheFile

	// File Header
	Read(&file.Header)

	// Model headers
	file.ModelHeaderList = make([]ModelHeader, file.Header.ModelCount)
	for i := 0; i < int(file.Header.ModelCount); i++ { //
		Read(&file.ModelHeaderList[i])
	}

	// Model data
	file.ModelDataList = make([]ModelData, file.Header.ModelCount)
	for i := 0; i < int(file.Header.ModelCount); i++ {
		model := &file.ModelDataList[i]

		Read(&model.FirstName)
		Read(&model.Unknown1)
		Read(&model.Unknown2)
		Read(&model.LastName)
		Read(&model.Unknown3)

		// Meshes
		offset := file.Header.OffsetOfModelDataList
		offset += file.ModelHeaderList[i].RelativeOffset
		offset += file.ModelHeaderList[i].OffsetOfMeshData
		buf.Seek(int64(offset), io.SeekStart)
		model.Meshes = make([]Mesh, 0, 5)

		for {
			// Mesh
			model.Meshes = append(model.Meshes, Mesh{})
			LastMesh := &model.Meshes[len(model.Meshes)-1]

			Read(&LastMesh.CountOfTotalRows)
			Read(&LastMesh.Unknown1)
			Read(&LastMesh.Unknown2)
			Read(&LastMesh.Unknown3)
			Read(&LastMesh.PrefixCount)
			Read(&LastMesh.SumOfRows)
			Read(&LastMesh.SuffixOfRows)
			Read(&LastMesh.InfoRows)
			LastMesh.TriStripRows = make([]StripRow, LastMesh.SumOfRows-2)
			Read(&LastMesh.TriStripRows)
			Read(&LastMesh.Unknown4)
			Read(&LastMesh.Unknown5)
			Read(&LastMesh.Unknown6)
			Read(&LastMesh.ElementHeader1)
			{
				// UvBlock
				Read(&LastMesh.UvBlock.Header)
				Read(&LastMesh.UvBlock.Unknown1)
				Read(&LastMesh.UvBlock.Unknown2)
				Read(&LastMesh.UvBlock.UVCountPrefix)
				Read(&LastMesh.UvBlock.CountOfUVs)
				Read(&LastMesh.UvBlock.UVCountSuffix)
				LastMesh.UvBlock.UVs = make([]UV, LastMesh.UvBlock.CountOfUVs)
				Read(&LastMesh.UvBlock.UVs)
				LastMesh.UvBlock.Filler = make([]byte, (LastMesh.UvBlock.CountOfUVs*8)%16)
				Read(&LastMesh.UvBlock.Filler)
			}

			Read(&LastMesh.ElementHeader2)
			{
				// NormalBlock
				Read(&LastMesh.NormBlock.Header)
				Read(&LastMesh.NormBlock.Unknown1)
				Read(&LastMesh.NormBlock.Unknown2)
				Read(&LastMesh.NormBlock.NormalCountPrefix)
				Read(&LastMesh.NormBlock.CountOfNormals)
				Read(&LastMesh.NormBlock.NormalCountSuffix)
				LastMesh.NormBlock.Normals = make([]Normal, LastMesh.NormBlock.CountOfNormals)
				Read(&LastMesh.NormBlock.Normals)
				LastMesh.NormBlock.Filler = make([]byte, (LastMesh.NormBlock.CountOfNormals*6)%16)
				Read(&LastMesh.NormBlock.Filler)
			}
			Read(&LastMesh.ElementHeader3)
			{
				// VertexBlock
				Read(&LastMesh.VertBlock.Header)
				Read(&LastMesh.VertBlock.Unknown1)
				Read(&LastMesh.VertBlock.Unknown2)
				Read(&LastMesh.VertBlock.VertexCountPrefix)
				Read(&LastMesh.VertBlock.CountOfVertices)
				Read(&LastMesh.VertBlock.VertexCountSuffix)
				LastMesh.VertBlock.Vertices = make([]Vertex, LastMesh.VertBlock.CountOfVertices)
				Read(&LastMesh.VertBlock.Vertices)
				LastMesh.VertBlock.Filler = make([]byte, (LastMesh.VertBlock.CountOfVertices*12)%16)
				Read(&LastMesh.VertBlock.Filler)

			}
			Read(&LastMesh.Footer)
			Read(&LastMesh.Unknown7)

			// Check if this is the last mesh
			// Use pattern search to see if there is a model footer close by

		}

	}

	return
}

func PrintLoc(buf *bytes.Reader) {
	m, _ := buf.Seek(0, os.SEEK_CUR)
	fmt.Println(m)
}
