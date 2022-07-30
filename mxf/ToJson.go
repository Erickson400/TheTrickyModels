package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
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

func FileToStruct(f *os.File) (file TheFile) {
	// Setup the buffer for reading in little endian
	Read := func(p any) { //Takes in a Pointer/Address
		err := binary.Read(f, binary.LittleEndian, p)
		if err != nil {
			panic(err)
		}
	}

	// Start Reading TheFile

	// File Header
	fmt.Println("Reading File Header...")
	Read(&file.Header)
	fmt.Println("SUCCESS: Reading File Header")

	// Model headers
	fmt.Println("Reading Model Headers...")
	file.ModelHeaderList = make([]ModelHeader, file.Header.ModelCount)
	Read(&file.ModelHeaderList)
	fmt.Println("SUCCESS: Reading Model Headers")

	// Model data
	fmt.Println("\nReading Model Data...")
	file.ModelDataList = make([]ModelData, file.Header.ModelCount)
	for i := 0; i < len(file.ModelDataList); i++ {
		header := file.ModelHeaderList[i]
		model := &file.ModelDataList[i]
		modelStart := file.Header.ModelDataListOffset + header.RelativeOffset
		fmt.Println("Reading Model '" + string(header.Name[:]) + "' Data...")

		// Materials
		model.MaterialList = make([]Material, header.MaterialCount)
		f.Seek(int64(modelStart+header.MaterialDataOffset), io.SeekStart)
		Read(&model.MaterialList)

		// Bones
		model.BoneList = make([]Bone, header.BoneDataCount)
		f.Seek(int64(modelStart+header.BoneDataOffset1), io.SeekStart)
		Read(&model.BoneList)

		// IK Points
		model.IKList = make([]IK, header.IKDataCount)
		f.Seek(int64(modelStart+header.IKDataOffset), io.SeekStart)
		Read(&model.IKList)

		// Morph Header List
		model.MorphHeaderList = make([]MorphHeader, header.MorphHeaderCount)
		f.Seek(int64(modelStart+header.MorphHeaderListOffset), io.SeekStart)
		Read(&model.MorphHeaderList)
		//fmt.Println(len(model.MorphHeaderList))
		//PrintLoc(f)

		// Morph Data Container List
		model.MorphDataContainerList = make([]MorphDataContainer, header.MorphHeaderCount)
		for j := 0; j < len(model.MorphDataContainerList); j++ {
			model.MorphDataContainerList[j].Data = make([]MorphData, model.MorphHeaderList[j].Count)
			Read(&model.MorphDataContainerList[j].Data)
		}

		// Skinning Header List
		model.SkinningHeaderList = make([]SkinningHeader, header.SkinningHeaderCount)
		f.Seek(int64(modelStart+header.SkinningHeaderListOffset), io.SeekStart)
		Read(&model.SkinningHeaderList)

		// Skinning Data Container List
		model.SkinningDataContainerList = make([]SkinningDataContainer, header.SkinningHeaderCount)
		for j := 0; j < len(model.SkinningDataContainerList); j++ {
			model.SkinningDataContainerList[j].Data = make([]SkinningData, model.SkinningHeaderList[j].Count)
			Read(&model.SkinningDataContainerList[j].Data)
		}

		// Tristrip Header List
		model.TristripHeaderList = make([]TristripHeader, header.TristripGroupCount)
		f.Seek(int64(modelStart+header.TristripHeaderListOffset), io.SeekStart)
		Read(&model.TristripHeaderList)

		// Tristrip Data Container List
		model.TristripDataContainerList = make([]TristripDataContainer, header.TristripGroupCount)
		for j := 0; j < len(model.TristripDataContainerList); j++ {
			triheader := &model.TristripHeaderList[j]
			container := &model.TristripDataContainerList[j]
			container.Data = make([]uint16, triheader.IndexCount)
			Read(&container.Data)
		}

		// Vertex Data
		model.VertexDataList1 = make([]VertexData, header.VertexCount)
		model.VertexDataList2 = make([]VertexData, header.VertexCount)
		f.Seek(int64(modelStart+header.VertexDataOffset), io.SeekStart)
		Read(&model.VertexDataList1)
		Read(&model.VertexDataList2)
		fmt.Println("SUCCESS: Reading Model '" + string(header.Name[:]) + "' Data...")
	}
	fmt.Println("SUCCESS: Reading Model Data...")
	return
}
