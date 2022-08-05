package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func MxfToGltf(filename string, destination string) {

	// Read the File
	if !strings.HasSuffix(filename, ".mxf") {
		filename += ".mxf"
	}

	f, err := os.Open(filename)
	if err != nil {
		panic("Failed to read .mxf file")
	}

	// Turn the file's data to a struct then to gltf
	theFile := FileToStruct(f)
	StructToGltf(theFile, filename)
}

func StructToGltf(theFile TheFile, filename string) {
	doc := gltf.NewDocument()

	// For every model
	for i := 0; i < int(theFile.Header.ModelCount); i++ {
		modelHeader := &theFile.ModelHeaderList[i]
		modelData := &theFile.ModelDataList[i]

		// Vertices
		var Verts [][3]float32
		for _, v := range modelData.VertexDataList1 {
			Verts = append(Verts, [3]float32{v.LocationX, v.LocationY, v.LocationZ})
		}

		// Normals
		var Norms [][3]float32
		for _, v := range modelData.VertexDataList1 {
			Norms = append(Norms, [3]float32{v.NormalX, v.NormalY, v.NormalZ})
		}

		// UVs
		var UVs [][2]float32
		for _, v := range modelData.VertexDataList1 {
			UVs = append(UVs, [2]float32{v.UVMapU, v.UVMapV})
		}

		// Materials
		for _, v := range modelData.MaterialList {
			name := Clean(string(v.MainTextureName[:]))

			f, err := LoadTexture(name, filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			Idx, err := modeler.WriteImage(doc, name, "image/png", f)
			if err != nil {
				panic(err)
			}
			doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(Idx)})
			doc.Materials = append(doc.Materials, &gltf.Material{
				Name: "SussyMat",
				PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
					BaseColorTexture: &gltf.TextureInfo{Index: Idx},
					MetallicFactor:   gltf.Float(0),
				},
			})
		}

		// Make primitive meshes for each Tristrip container
		var primitives []*gltf.Primitive
		for j := 0; j < len(modelData.TristripDataContainerList); j++ {
			Idx := uint32(modelData.TristripHeaderList[j].MaterialID[0])
			p := &gltf.Primitive{
				Mode:    gltf.PrimitiveTriangleStrip,
				Indices: gltf.Index(modeler.WriteIndices(doc, modelData.TristripDataContainerList[j].Data)),
				Attributes: map[string]uint32{
					gltf.POSITION:   modeler.WritePosition(doc, Verts),
					gltf.NORMAL:     modeler.WriteColor(doc, Norms),
					gltf.TEXCOORD_0: modeler.WriteTextureCoord(doc, UVs),
				},
				Material: gltf.Index(Idx),
			}
			primitives = append(primitives, p)
		}

		// Make the mesh and append it to the document & Nodes
		mesh := &gltf.Mesh{
			Name:       string(modelHeader.Name[:]),
			Primitives: primitives,
		}

		if len(modelData.TristripDataContainerList) == 0 {
			continue
		}
		doc.Meshes = append(doc.Meshes, mesh)
		doc.Nodes = append(doc.Nodes, &gltf.Node{Name: "Root" + fmt.Sprint(i), Mesh: gltf.Index(uint32(i))})
	}

	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	gltf.Save(doc, "./Model.gltf")
}

func LoadTexture(TexName, filename string) (*os.File, error) {
	// Function checks if the directory has the required textures.

	// Find the textures in the directory
	p := path.Dir(filename)
	files, err := os.ReadDir(p)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.Contains(f.Name(), TexName) {
			t, err := os.Open(p + "/" + f.Name())
			if err != nil {
				panic(err)
			}
			return t, nil
		}
	}
	return nil, fmt.Errorf("ERROR: Cant find texture image " + TexName)
}

func Clean(name string) string {
	return strings.ReplaceAll(name, " ", "")
}
