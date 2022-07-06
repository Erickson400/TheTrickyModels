package main

type FileHeader struct {
	Unknown               uint32 // Always stores 4
	ModelCount            uint16
	ModelHeaderListOffset uint16 // Points to first ModelHeader
	ModelRootOffset       uint32

	/*
		ModelHeaders is not part of the file header, but we put it here
		so that we know it starts right after it.
	*/
	ModelHeaders []ModelHeader // size is ModelCount
}

type ModelHeader struct {
	ModelName [16]byte

	/*
		The decoder adds this to FileHeader.ModelRootOffset.
		The sum is where the Model starts.
		ModelStart is the alias for ModelRelativeOffset + FileHeader.ModelRootOffset.
	*/
	ModelRelativeOffset   uint32
	ModelSize             uint32 // size represented in bytes
	OffsetOfModelData     uint32 // Relative to ModelStart. Points to Model
	OffsetOfBoneWeights1  uint32 // Relative to ModelStart
	Unknown1              uint32
	OffsetOfMeshData      uint32 // Relative to ModelStart. Points to first Mesh in the model's Mesh array
	Unknown2              uint32
	OffsetOfBoneWeights2  uint32 // Relative to ModelStart
	OffsetOfNumListRef    uint32 // Relative to ModelStart
	OffsetOfBoneWeights3  uint32 // Relative to ModelStart
	Unknown3              uint32
	Unknown4              uint32
	Unknown5              uint16
	Unknown6              uint16
	CountOfBoneWeights    uint32 // Relative to ModelStart
	CountOfInternalMeshes uint16 // Ammount of Meshes
	Unknown7              uint16
	CountOfBones          uint16
	Unknown8              [8]byte
}

type Model struct {
	// Missing Stuff...
	MeshData []Mesh // Size is CountOfInternalMeshes from the ModelHeader
}

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
type uint24 struct {
	b [3]byte
}

type TriStripCountRow struct {
	TriStripLength uint32
	padding        [12]byte
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
	Vertices          []VertexData // Size of CountOfNormals
}
type VertexData struct {
	X float32
	Y float32
	Z float32
}
