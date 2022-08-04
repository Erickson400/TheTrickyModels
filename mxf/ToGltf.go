package main

import (
	"bytes"
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

		// Read texture files from directory
		MatTextureFiles := LoadTextures(modelData.MaterialList, filename)
		//fmt.Println(MatTextureFiles)

		// Materials
		img_file, err := os.ReadFile("resources/textures_elise/elise1_suit.157.png")
		if err != nil {
			panic(err)
		}

		imageIdx, err := modeler.WriteImage(doc, "suit", "image/png", bytes.NewReader(img_file))
		if err != nil {
			panic(err)
		}

		doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})
		doc.Materials = append(doc.Materials, &gltf.Material{
			Name: "Material_SuS",
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
				BaseColorTexture: &gltf.TextureInfo{Index: uint32(0)},
				MetallicFactor:   gltf.Float(0),
			},
		})

		// Make primitive meshes for each Tristrip container
		var primitives []*gltf.Primitive
		for j := 0; j < len(modelData.TristripDataContainerList); j++ {
			p := &gltf.Primitive{
				Mode:    gltf.PrimitiveTriangleStrip,
				Indices: gltf.Index(modeler.WriteIndices(doc, modelData.TristripDataContainerList[j].Data)),
				Attributes: map[string]uint32{
					gltf.POSITION:   modeler.WritePosition(doc, Verts),
					gltf.NORMAL:     modeler.WriteColor(doc, Norms),
					gltf.TEXCOORD_0: modeler.WriteTextureCoord(doc, UVs),
				},
				Material: gltf.Index(0),
			}
			primitives = append(primitives, p)
		}

		// Make the mesh and append it to the document & Nodes
		mesh := &gltf.Mesh{
			Name:       string(modelHeader.Name[:]),
			Primitives: primitives, //[]*gltf.Primitive{
			// 	{
			// 		Mode:    gltf.PrimitiveTriangleStrip,
			// 		Indices: gltf.Index(modeler.WriteIndices(doc, Tris)),
			// 		Attributes: map[string]uint32{
			// 			gltf.POSITION: modeler.WritePosition(doc, Verts),
			// 			gltf.COLOR_0:  modeler.WriteColor(doc, Norms),
			// 		},
			// 	},
			// },
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

func LoadTextures(MatList []Material, filename string) map[string]*os.File {
	// Function checks if the directory has the required textures.

	// Find required textures and remove duplicates
	var NeededTextures = make(map[string]*os.File)
	for _, v := range MatList {
		// Remove empty spaces from string
		MatName := strings.ReplaceAll(string(v.MainTextureName[:]), " ", "")

		// Check if name already exists
		_, exists := NeededTextures[MatName]
		if exists {
			continue
		}
		NeededTextures[MatName] = &os.File{}
	}

	// Find the textures in the directory
	for k := range NeededTextures {
		p := path.Dir(filename)
		files, err := os.ReadDir(p)
		if err != nil {
			panic(err)
		}
		found := false
		for _, f := range files {
			if strings.Contains(f.Name(), k) {
				t, err := os.Open(p + "/" + f.Name())
				if err != nil {
					panic(err)
				}
				defer t.Close()

				NeededTextures[k] = t
				found = true
			}
		}
		if !found {
			panic("Error: Texture not Found.")
		}
	}
	return NeededTextures
}
