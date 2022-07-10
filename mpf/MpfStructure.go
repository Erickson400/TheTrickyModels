package main

/*
	-> means holds
	Hirarchy: Mpf File -> Models -> Meshes -> TriStrips.

	Model = ModelHeader + ModelData
*/

type TheFile struct {
	Header          FileHeader
	ModelHeaderList []ModelHeader
	ModelDataList   []ModelData
}

type FileHeader struct {
	/*
		Always stores 4.
		Might be the version of the format or decoder.
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

	Unknown2             uint32
	OffsetOfBoneWeights2 uint32 // Relative to ModelStart
	OffsetOfNumListRef   uint32 // Relative to ModelStart
	OffsetOfBoneWeights3 uint32 // Relative to ModelStart
	Unknown3             uint32
	Unknown4             uint16
	Unknown5             uint16 // Count?
	Unknown6             uint16 // Bone Count?
	CountOfBoneWeights   uint16 // Relative to ModelStart

	/*
	 Ammount of Meshes
	*/
	CountOfInternalMeshes uint16
	Unknown7              uint16
	CountOfBones          uint16
	Unknown8              uint16
	FillerPadding         uint32
}

type ModelData struct {
	FirstName [4]byte
	Unknown1  uint32 // stores 0x00202020
	Unknown2  uint32 // stores 0x00202020
	LastName  [4]byte
	Unknown3  uint32 // stores 0x00202020

	// ...Missing stuff
}

/*

 */
type Mesh struct {
	CountOfTotalRows  uint24
	Unknown1          byte
	Unknown2          [12]byte
	Unknown3          [13]byte
	PrefixCount       byte
	CountMeshInfoRows byte // MeshInfoRows + TriStripRows
	SuffixOfRows      byte
	MeshInfoRow1      [16]byte
	MeshInfoRow2      [16]byte
	TriStripCountRows []TriStripCountRow // Size is CountMeshInfoRows - 2
	CountOfVertices   uint24
	Unknown4          byte
	Unknown5          [12]byte
	UvBlock           UVBlock
	NormBlock         NormalBlock
	VertBlock         VertexBlock
}

type TriStripCountRow struct {
	TriStripLength uint32
	Padding        [12]byte
}

type UVBlock struct {
	UVBlockHeader [16]byte
	Unknown1      [12]byte
	Unknown2      byte
	UVCountPrefix byte
	CountOfUVs    byte
	UVCountSuffix byte
	UV            []UVData // Size of CountOfUVs
}

type UVData struct {
	U         uint16
	V         uint16
	UDistance uint16
	VDistance uint16
}

type NormalBlock struct {
	NormalBlockHeader [16]byte
	Unknown1          [12]byte
	Unknown2          byte
	NormalCountPrefix byte
	CountOfNormals    byte
	NormalCountSuffix byte
	Normals           []NormalData // Size of CountOfNormals
}

type NormalData struct {
	X uint16
	Y uint16
	Z uint16
}

type VertexBlock struct {
	VertexBlockHeader [16]byte
	Unknown1          [12]byte
	Unknown2          byte
	VertexCountPrefix byte
	CountOfVertices   byte
	VertexCountSuffix byte
	Vertices          []VertexData // Size of CountOfVertices
}

type VertexData struct {
	X float32
	Y float32
	Z float32
}

type uint24 struct {
	B [3]byte
}
