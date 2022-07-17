package main

/*
	-> means holds.
	Hirarchy: Mpf File -> Models -> Meshes -> TriStrips.

	Model = ModelHeader + ModelData.
	i.e: ModelHeader 34 and ModelData 34 make up Model 34.
	The word Model is just a concept, not an actual struct.

*/

type TheFile struct {
	Header          FileHeader
	ModelHeaderList []ModelHeader // Size is Header.ModelCount
	Filler1         [4]byte
	ModelDataList   []ModelData // Size is Header.ModelCount
}

type FileHeader struct {
	/*
		Always stores 8.
		Might be the version/type of the format or decoder for PS2.
	*/
	Version uint32

	/*
		The amount of Models.
		TheFile.ModelHeaderList has ModelCount amount of headers.
		TheFile.ModelDataList has ModelCount amount of Data.
	*/
	ModelCount uint16

	/*
		Points to the first ModelHeader (from TheFile.ModelHeaderList)
	*/
	ModelHeaderListOffset uint16

	/*
		Points to the first ModelData (from TheFile.ModelDataList)
	*/
	ModelDataListOffset uint32
}

type ModelHeader struct {
	/*
		Name of model (ASCII string of bytes)
	*/
	Name [16]byte

	/*
		It points to the current model from TheFile.ModelDataList
		Its relative to FileHeader.ModelDataListOffset.

		The decoder adds this to FileHeader.ModelDataListOffset.
		The sum is where the Model starts.

		ModelStart is the alias for RelativeOffset + FileHeader.ModelDataListOffset.
	*/
	RelativeOffset uint32

	/*
		Represented in bytes.
	*/
	Size uint32

	/*
		Points to the first Bone in the list
	*/
	BoneListOffset uint32 // Relative to ModelStart

	BoneWeightsOffset1 uint32 // Relative to ModelStart
	Unknown1           uint32

	/*
		Points to first Mesh in the Mesh group list
	*/
	MeshDataOffset uint32 // Relative to ModelStart.

	Unknown2           uint32
	BoneWeightsOffset2 uint32 // Relative to ModelStart
	NumListRefOffset   uint32 // Relative to ModelStart
	BoneWeightsOffset3 uint32 // Relative to ModelStart
	Unknown3           uint32
	Unknown4           uint16
	Unknown5           uint16 // Count?
	Unknown6           uint16 // Bone Count?
	BoneDataCount      uint16
	MaterialCount      uint16
	Unknown7           uint16
	Unknown8           uint16
	Unknown9           uint16
	FillerPadding      uint32
}

type ModelData struct {
	MaterialList []Material // Size is ModelHeader.MaterialCount
	BoneList     []Bone     // Size is ModelHeader.BoneDataCount
	Unknown1     []byte     // Size is ModelHeader.MeshDataOffset

	/*
		Starts at ModelHeader.ModelStart + Modelheader.OffsetOfMeshData
	*/
	MeshGroupList []MeshGroup // Unknown Size
}

type Material struct {
	MainTextureName [4]byte
	TextureType1    [4]byte
	TextureType2    [4]byte
	TextureType3    [4]byte
	TextureType4    [4]byte
	Unknown1        [3]float32
}

type Bone struct {
	Name         [16]byte
	Unknown1     uint16
	ParentBoneID uint16
	Unknown2     uint16
	ID           uint16
	X            float32
	Y            float32
	Z            float32
	RotRadianX1  float32
	RotRadianY1  float32
	RotRadianZ1  float32
	RotRadianX2  float32
	RotRadianY2  float32
	RotRadianZ2  float32

	/*
		Contains 6 float values with either -1.0 or 1.0
	*/
	Unknown3 [6]float32
}

type MeshGroup struct {
	MeshList []Mesh // Unknown Size

	// Notes
	// at the end of mesh it can be x10 or x60
	/*
		x10 A 14 0 0 appears in between meshes in a group
		x10 A 14 8 14 6 14 appears in between meshes in a group
		x10 A 14 6 14 6 14 appears in between meshes in a group
		x10 A 14 6 14 6 14 appears at the end of a model
		x10 A 14 8 14 6 14 appears at the end of a model

		x60 1101 appears in between meshes in a group
		x60 00001101 appears in between meshes in a group

		0000 10 0000 14 may be the group footer.
		modelData dont have footers.



				x60 always ends in:
					00000000 01000010 00000000 00000000

				x10 always ends in:
					01010001 0A000014 00000000 00000000
					or
					01010001 0A000014 06000014 06000014
						when its the end of the model.
					or
					01010001 0A000014 08000014 06000014
					or
					any row with 01010001



			first model:
				group list {
					{mesh1 A 14 0 0,
					 mesh2 A 14 6 14 6 14,
					},
					x60 holds 01010001 in 2nd 32byte
					x60 holds 01010001 in 4th 32byte
				}


			second model:
				group list {
					{mesh1 A 14 0 0,
					 mesh2 A 14 6 14 6 14,
					 mesh3 A 14 8 14 6 14,
					},
					x60 holds 01010001 in 2nd 32byte
					x60 holds 01010001 in 4th 32byte
				}

	*/

	/*
		Stores 00000000 00000010 00000000 00000014
	*/
	Footer [16]byte
}

type Mesh struct {
	TriStripRowHeader RowHeader

	/*
		Filler/Padding/Footer
	*/
	Filler [13]byte

	/*
		Always 0x80
	*/
	CountPrefix byte

	/*
		Stores an amount of rows: InfoRows + TriStripRows.

		The first row from the count is below this row.
	*/
	StripRowCount byte

	/*
		Always 0x6C
	*/
	CountSuffix byte
	InfoRows    [2]Row

	/*
		Stores each strip's length
	*/
	StripRowList []StripRow // Size is TriStripDataRowCount - 2

	BlockRowHeader RowHeader
	ElementHeader1 [16]byte // Stores 00000000 00000030 00000000 00000000
	UvBlock        UVBlock
	ElementHeader2 [16]byte // Stores 00000000 00000030 00000000 00000000
	NormBlock      NormalBlock
	ElementHeader3 [16]byte // Stores 00000000 00000030 00000000 00000000
	VertBlock      VertexBlock
}

type RowHeader struct {
	/*
		RowHeader stores the amount of rows below the RowHeader.
		It can tell us the section's amount of data in advance.
		e.g:
			The Mesh.BlockRowHeader is how many rows the Mesh's UVBlock, NormalBlock, UVBlock
			and element headers take up.
	*/
	RowCount uint16

	/*
		0x10 = holds tristrip row counts, or block data
		0x60 = holds a mesh and model footer, and maybe even some unknown data.
	*/
	Type   uint16 // 0x10 or 0x60
	Filler [12]byte
}

type Row struct {
	B [4]uint32
}

type StripRow struct {
	VerticesCount uint32
	Padding       [12]byte
}

type UVBlock struct {

	/*
		Stores 00100000 00100000 00000020 50505050 (aka PPPP)
	*/
	Header        [16]byte
	Unknown1      [13]byte
	UVCountPrefix byte
	UVCount       byte
	UVCountSuffix byte
	UVList        []UV // Size is UVCount

	/*
		Size is CountOfUVs * 8, then moduled by 16.
		This is so that the UVs always take up a row. If it doesnt fill the whole row
		then the remaining bytes will be set to 0s.

		e.g.:
			CountOfUVs = 58
			TotalBytes = CountOfUVs * 8
			FillerSize = TotalBytes % 16
	*/
	Filler []byte
}

type UV struct {
	U         uint16
	V         uint16
	UDistance uint16
	VDistance uint16
}

type NormalBlock struct {
	/*
		Stores 00000000 00800000 00000020 40404040 (aka @@@@)
	*/
	Header [16]byte

	Unknown1          [12]byte
	Unknown2          byte
	NormalCountPrefix byte
	CountOfNormals    byte
	NormalCountSuffix byte
	Normals           []Normal // Size is CountOfNormals
	/*
		Size is CountOfNormals * 6, then moduled by 16.
		This is so that the Normalss always take up a row. If it doesnt fill the whole row
		then the remaining bytes will be set to 0s.
		e.g.:
			CountOfNormals = 58
			TotalBytes = CountOfNormals * 6
			FillerSize = TotalBytes % 16
			FillerSize = 16 - FillerSize
	*/
	Filler []byte
}

type Normal struct {
	X uint16
	Y uint16
	Z uint16
}

type VertexBlock struct {
	/*
		Stores 00000000 0000803F 00000020 40404040 (aka @@@@)
	*/
	Header [16]byte

	Unknown1          [12]byte
	Unknown2          byte
	VertexCountPrefix byte
	CountOfVertices   byte
	VertexCountSuffix byte
	Vertices          []Vertex // Size of CountOfVertices

	/*
		Size is CountOfVertices * 12, then moduled by 16.
		This is so that the Vertices always take up a row. If it doesnt fill the whole row
		then the remaining bytes will be set to 0s.
		e.g.:
			CountOfVertices = 58
			TotalBytes = CountOfVertices * 12
			FillerSize = TotalBytes % 16
			FillerSize = 16 - FillerSize
	*/
	Filler []byte
}

type Vertex struct {
	X float32
	Y float32
	Z float32
}

type Uint24 struct {
	B [3]byte
}
