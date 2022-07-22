package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
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
	Read(&file.ModelHeaderList)
	Read(&file.Filler1)

	// Model data
	file.ModelDataList = make([]ModelData, file.Header.ModelCount)
	for i := 0; i < len(file.ModelDataList); i++ {
		header := file.ModelHeaderList[i]
		model := &file.ModelDataList[i]

		// Materials
		model.MaterialList = make([]Material, header.MaterialCount)
		Read(&model.MaterialList)

		// Bones
		model.BoneList = make([]Bone, header.BoneDataCount)
		Read(&model.BoneList)

		// Unknown1
		model.Unknown1 = make([]byte, header.MeshDataOffset)
		Read(&model.Unknown1)

		// Mesh Groups
		model.MeshGroupList = make([]MeshGroup, header.GroupCount)
		for j := 0; j < len(model.MeshGroupList); j++ {
			groupList := &model.MeshGroupList[j]
			ReadGroup(Read, buf, header, model, groupList)
			Read(&groupList.Footer)
		}

	}

	return
}

func ReadGroup(Read func(p any), buf *bytes.Reader, header ModelHeader, model *ModelData, groups *MeshGroup) {
	//Check if its a shadow mesh
	if strings.Contains(string(header.Name[:]), "Shdw") || strings.Contains(string(header.Name[:]), "shdw") {
		groups.MeshShadowList = make([]MeshShadow, 1)
		mesh := &groups.MeshShadowList[0]
		ReadMeshShadow(Read, buf, mesh)
		return
	}

}

func ReadMeshShadow(Read func(p any), buf *bytes.Reader, mesh *MeshShadow) {
	Read(&mesh.TriStripRowHeader)
	Read(&mesh.Filler)
	Read(&mesh.CountPrefix)
	Read(&mesh.StripRowCount)
	Read(&mesh.CountSuffix)
	Read(&mesh.InfoRows)
	mesh.StripRowList = make([]StripRow, mesh.StripRowCount)
	Read(&mesh.StripRowList)
	Read(&mesh.BlockRowHeader)
	Read(&mesh.ElementHeader1)
	ReadVertexBlock(Read, &mesh.VertBlock)

	// Read morph data until a group footer is found
	r := Row{}
	groupFooter := Row{
		B: [4]uint32{0x00000000, 0x00000010, 0x00000000, 0x00000014},
	}
	for {
		Read(&r)
		if reflect.DeepEqual(r, groupFooter) {
			buf.Seek(-16, io.SeekCurrent)
			return
		}
	}
}

func ReadVertexBlock(Read func(p any), block *VertexBlock) {
	Read(&block.Header)
	Read(&block.Unknown1)
	Read(&block.Unknown2)
	Read(&block.VertexCountPrefix)
	Read(&block.VertexCount)
	Read(&block.VertexCountSuffix)
	block.Vertices = make([]Vertex, block.VertexCount)
	Read(&block.Vertices)
	block.Filler = make([]byte, 16-(block.VertexCount*12)%16)
	Read(&block.Filler)
}
