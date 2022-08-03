package main

import (
	"fmt"
	"os"
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
	StructToGltf(theFile)
}

func StructToGltf(theFile TheFile) {
	doc := gltf.NewDocument()

	// For every model
	for i := 0; i < int(theFile.Header.ModelCount); i++ {
		modelHeader := theFile.ModelHeaderList[i]
		modelData := theFile.ModelDataList[i]

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

		// Tristrips

		//var Tris []uint16
		//for _, v := range modelData.TristripDataContainerList {
		//	Tris = append(Tris, v.Data...)
		//}

		// Make primitive meshes for each tristrip container
		var primitives []*gltf.Primitive

		for j := 0; j < len(modelData.TristripDataContainerList); j++ {
			//fmt.Println(modelData.TristripDataContainerList[j].Data)
			fmt.Println(len(modelData.VertexDataList1))
			p := &gltf.Primitive{
				Mode:    gltf.PrimitiveTriangleStrip,
				Indices: gltf.Index(modeler.WriteIndices(doc, modelData.TristripDataContainerList[j].Data)),
				Attributes: map[string]uint32{
					gltf.POSITION: modeler.WritePosition(doc, Verts),
					gltf.COLOR_0:  modeler.WriteColor(doc, Norms),
				},
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
		doc.Meshes = append(doc.Meshes, mesh)
		doc.Nodes = append(doc.Nodes, &gltf.Node{Name: "Root" + fmt.Sprint(i), Mesh: gltf.Index(uint32(i))})
	}

	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	gltf.Save(doc, "./Model.gltf")
}
