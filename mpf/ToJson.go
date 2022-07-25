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
	fmt.Println("Reading File Header...")
	Read(&file.Header)
	fmt.Println("SUCCESS: Reading File Header")

	// Model headers
	fmt.Println("Reading Model Headers...")
	file.ModelHeaderList = make([]ModelHeader, file.Header.ModelCount)
	Read(&file.ModelHeaderList)
	Read(&file.Filler1)
	fmt.Println("SUCCESS: Reading Model Headers")

	// Model data
	file.ModelDataList = make([]ModelData, file.Header.ModelCount)
	for i := 0; i < len(file.ModelDataList); i++ {
		header := file.ModelHeaderList[i]
		model := &file.ModelDataList[i]
		modelStart := file.Header.ModelDataListOffset + header.RelativeOffset

		fmt.Println("Reading Model '" + string(header.Name[:]) + "' Data...")

		fmt.Println("Reading Model '" + string(header.Name[:]) + "' Materials, Bones, Unk1...")
		// Materials
		model.MaterialList = make([]Material, header.MaterialCount)
		Read(&model.MaterialList)

		// Bones
		model.BoneList = make([]Bone, header.BoneDataCount)
		Read(&model.BoneList)

		// IK Points
		model.IKPointList = make([]IK, header.IKPointCount)
		Read(&model.IKPointList)

		// Unknown1
		loc, _ := buf.Seek(0, io.SeekCurrent)
		model.Unknown1 = make([]byte, int64(modelStart+header.MeshDataOffset)-loc)
		Read(&model.Unknown1)
		fmt.Println("SUCCESS: Reading Model '" + string(header.Name[:]) + "' Materials, Bones, Unk1")

		// Mesh Groups
		fmt.Println("Reading Mesh Groups...")
		fmt.Println(header.GroupOffsetData + modelStart)

		model.MeshGroupList = make([]MeshGroup, header.GroupCount)
		for j := 0; j < len(model.MeshGroupList); j++ {
			fmt.Println("Reading Mesh Groups")
			group := &model.MeshGroupList[j]
			ReadGroup(Read, buf, header, model, group)
			Read(&group.Unknown1)
			Read(&group.Footer)

		}
	}
	return
}

func ReadGroup(Read func(p any), buf *bytes.Reader, header ModelHeader, model *ModelData, group *MeshGroup) {

	// Check if its a Face model
	if strings.Contains(string(header.Name[:]), "Face") || strings.Contains(string(header.Name[:]), "face") {
		group.MeshFaceList = make([]FaceMesh, 0, 1)
		for {
			group.MeshFaceList = append(group.MeshFaceList, FaceMesh{})
			mesh := &group.MeshFaceList[len(group.MeshFaceList)-1]
			ReadFaceMesh(Read, header, mesh)
			if CheckForFooter(Read, buf) {
				return
			}
		}
	}

	//Check if its a Shadow model
	if strings.Contains(string(header.Name[:]), "Shdw") || strings.Contains(string(header.Name[:]), "shdw") {
		group.MeshShadowList = make([]ShadowMesh, 0, 1)
		for {
			group.MeshShadowList = append(group.MeshShadowList, ShadowMesh{})
			mesh := &group.MeshShadowList[len(group.MeshShadowList)-1]
			ReadMeshShadow(Read, buf, mesh)
			if CheckForFooter(Read, buf) {
				return
			}
		}
	}

	// Read Default Meshes if non above apply

	group.MeshDefaultList = make([]DefaultMesh, 0, 1)
	for {
		fmt.Println("-----Sus")
		group.MeshDefaultList = append(group.MeshDefaultList, DefaultMesh{})
		mesh := &group.MeshDefaultList[len(group.MeshDefaultList)-1]
		ReadMeshDefault(Read, buf, mesh)
		//PrintLoc(buf)

		if CheckForFooter(Read, buf) {
			fmt.Println("AAAAAAAAAAHHHHHHHHHH")

			return
		}
	}
}

func ReadFaceMesh(Read func(p any), header ModelHeader, mesh *FaceMesh) {
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
	ReadUVBlock(Read, &mesh.UvBlock)
	Read(&mesh.ElementHeader2)
	ReadNormalBlock(Read, &mesh.NormBlock)
	Read(&mesh.ElementHeader3)
	ReadVertexBlock(Read, &mesh.VertBlock)

	Read(&mesh.Unknown1)
	Read(&mesh.Unknown2)
	mesh.MorphList = make([]Morph, header.MorphCount)
	for i := 0; i < len(mesh.MorphList); i++ {
		Read(&mesh.MorphList[i].MorphHeader)
		mesh.MorphList[i].R = make([]Row, mesh.MorphList[i].MorphHeader.RowCount)
		Read(&mesh.MorphList[i].R)
	}
}

func ReadMeshDefault(Read func(p any), buf *bytes.Reader, mesh *DefaultMesh) {
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
	ReadUVBlock(Read, &mesh.UvBlock)
	Read(&mesh.ElementHeader2)
	ReadNormalBlock(Read, &mesh.NormBlock)
	Read(&mesh.ElementHeader3)
	ReadVertexBlock(Read, &mesh.VertBlock)
	Read(&mesh.Unknown1)
	Read(&mesh.Unknown2)
}

func ReadMeshShadow(Read func(p any), buf *bytes.Reader, mesh *ShadowMesh) {
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
	Read(&mesh.Unknown1)
	Read(&mesh.Unknown2)
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

func ReadUVBlock(Read func(p any), block *UVBlock) {
	Read(&block.Header)
	Read(&block.Unknown1)
	Read(&block.UVCountPrefix)
	Read(&block.UVCount)
	Read(&block.UVCountSuffix)
	block.UVList = make([]UV, block.UVCount)
	Read(&block.UVList)
	block.Filler = make([]byte, 16-(block.UVCount*8)%16)
	Read(&block.Filler)
}

func ReadNormalBlock(Read func(p any), block *NormalBlock) {
	Read(&block.Header)
	Read(&block.Unknown1)
	Read(&block.Unknown2)
	Read(&block.NormalCountPrefix)
	Read(&block.NormalCount)
	Read(&block.NormalCountSuffix)
	block.Normals = make([]Normal, block.NormalCount)
	Read(&block.Normals)
	block.Filler = make([]byte, 16-(block.NormalCount*6)%16)
	Read(&block.Filler)
}

func CheckForFooter(Read func(p any), buf *bytes.Reader) bool {
	r := [4]uint32{0, 0, 0, 0}
	groupFooter := [4]uint32{0x00000000, 0x00000010, 0x00000000, 0x00000014}

	buf.Seek(int64(-4), io.SeekCurrent)
	i := 0
	for i < 5 {
		Read(&r)
		if r == groupFooter {
			return true
		}
		i++
	}
	buf.Seek(int64(-16*i), io.SeekCurrent)
	return false
}
