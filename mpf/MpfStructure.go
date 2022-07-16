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
	// .. Missing bone data

	MeshGroupList []MeshGroup // Unknown Size

	/*
		Starts at ModelHeader.ModelStart + Modelheader.OffsetOfMeshData
	*/

	Meshes []Mesh // Unknown Size

	/*
		00000000 00000010 00000000 00000014
	*/
	Footer [32]byte
}

type Material struct {
	MainTextureName [4]byte
	TextureType1    [4]byte
	TextureType2    [4]byte
	TextureType3    [4]byte
	TextureType4    [4]byte
	Unknown1        [3]float32
}

type MeshGroup struct {
	Meshes []Mesh // Unknown Size
}

type RowHeader struct {
	RowCount uint16
	Type     uint16 // 0x10 or 0x60
	Filler   [12]byte
}

type Mesh struct {
	RowHeader1 RowHeader

	/*
		Filler/Padding
	*/
	Unknown3 [13]byte

	/*
		Always 0x80
	*/
	PrefixCount byte

	/*
		Stores an amount of rows:
		InfoRows + TriStripRows
	*/
	SumOfRows byte

	/*
		Always 0x6C
	*/
	SuffixOfRows byte
	InfoRows     [2]Row

	/*
		Stores each tristrip's length
	*/
	TriStripRows []StripRow // Size is SumOfRows - 2

	/*
		Might be a vert count
	*/
	Unknown4 Uint24

	/*
		Always 10
	*/
	Unknown5       byte
	Unknown6       [12]byte
	ElementHeader1 [16]byte // Stores 00000000 00000030 00000000 00000000
	UvBlock        UVBlock
	ElementHeader2 [16]byte // Stores 00000000 00000030 00000000 00000000
	NormBlock      NormalBlock
	ElementHeader3 [16]byte // Stores 00000000 00000030 00000000 00000000
	VertBlock      VertexBlock
	/*
		Stores 01000010 00000000 00000000 00000000 01010001 0A000014
	*/
	Footer   [24]byte
	Unknown7 [8]byte
}

type Row struct {
	B [4]uint32
}

type StripRow struct {
	CountOfVertices uint32
	Padding         [12]byte
}

type TriStripCountRow struct {
	TriStripLength uint32
	Padding        [12]byte
}

type UVBlock struct {

	/*
		Stores 00100000 00100000 00000020 50505050 (aka PPPP)
	*/
	Header        [16]byte
	Unknown1      [12]byte
	Unknown2      byte
	UVCountPrefix byte
	CountOfUVs    byte
	UVCountSuffix byte
	UVs           []UV // Size is CountOfUVs

	/*
		Size is CountOfUVs * 8, then moduled by 16.
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
