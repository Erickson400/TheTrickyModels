package main

/*
	-> means holds
	Hirarchy: Mpf File -> Models -> Meshes -> TriStrips.

	Model = ModelHeader + ModelData
*/

type TheFile struct {
	Header          FileHeader
	ModelHeaderList []ModelHeader // Size is Header.ModelCount
	ModelDataList   []ModelData   // Size is Header.ModelCount
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
	OffsetOfModelHeaderList uint16

	/*
		Points to the first ModelData (from TheFile.ModelDataList)
	*/
	OffsetOfModelDataList uint32
}

type ModelHeader struct {
	/*
		Name of model (ASCII string of bytes)
	*/
	Name [16]byte

	/*
		It points to the current model from TheFile.OffsetOfModelDataList.
		Its relative to FileHeader.OffsetOfModelDataList.

		The decoder adds this to FileHeader.OffsetOfModelDataList.
		The sum is where the Model starts.
		ModelStart is the alias for ModelRelativeOffset + FileHeader.OffsetOfModelDataList.
	*/
	RelativeOffset uint32

	/*
		Size represented in bytes.
		Tells how much data the model takes. Good for reading an entire model
		 all based on the ModelHeader.
		This is good because not all models are the same size.
	*/
	Size uint32

	/*
		Not yet clear what this is.
		Model data offset goes to some data above the tristrips.
	*/
	OffsetToDataAboveTristrips uint32 // Relative to ModelStart

	OffsetOfBoneWeights1 uint32 // Relative to ModelStart
	Unknown1             uint32

	/*
		Points to first Mesh in the model's Mesh array
	*/
	OffsetOfMeshData uint32 // Relative to ModelStart.

	Unknown2              uint32
	OffsetOfBoneWeights2  uint32 // Relative to ModelStart
	OffsetOfNumListRef    uint32 // Relative to ModelStart
	OffsetOfBoneWeights3  uint32 // Relative to ModelStart
	Unknown3              uint32
	Unknown4              uint16
	Unknown5              uint16 // Count?
	Unknown6              uint16 // Bone Count?
	CountOfBoneWeights    uint16 // Relative to ModelStart
	CountOfInternalMeshes uint16
	Unknown7              uint16
	Unknown8              uint16
	Unknown9              uint16
	FillerPadding         uint32
}

type ModelData struct {
	FirstName [4]byte
	Unknown1  uint32 // stores 0x00202020
	Unknown2  uint32 // stores 0x00202020
	LastName  [4]byte
	Unknown3  uint32 // stores 0x00202020
	// ...Unknown Missing Data

	/*
		Starts at ModelHeader.ModelStart + Modelheader.OffsetOfMeshData
	*/
	Meshes []Mesh // Unknown amount of meshes
	// ...Unknown Missing Data

	/*
		01000060 00000000 00000000 00000000
		00000000 01010001 00000000 00000000
		or
		00000000 00000000 00000000 01010001
		00000000 00000010 00000000 00000014
	*/
	Footer [32]byte
	// ...Unknown Missing Data
}

type Mesh struct {
	/*
		TriStripCountRow + InfoRows
	*/
	CountOfTotalRows Uint24

	/*
		Always 10
	*/
	Unknown1 byte

	/*
		Filler/Padding
	*/
	Unknown2 [12]byte

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
	UVData        []UV // Size is CountOfUVs

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
			TotalBytes = CountOfNormals * 8
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
