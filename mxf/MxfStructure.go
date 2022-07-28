package main

type TheFile struct {
	Header          FileHeader
	ModelHeaderList []ModelHeader // Size is Header.ModelCount
	ModelDataList   []ModelData   // Size is Header.ModelCount
}

type FileHeader struct { // Size is 12
	Version               uint32 // Always stores 4
	ModelCount            uint16
	ModelHeaderListOffset uint16 // Points to first ModelHeader
	ModelDataListOffset   uint32 // Points to first ModelData
}

type ModelHeader struct { // Size is 396
	Name [16]byte

	/*
		It points to the current model from TheFile.ModelDataList.
		Its relative to TheFile.Header.ModelDataListOffset.

		The decoder adds this to TheFile.Header.ModelDataListOffset.
		The sum is where the Model starts.

		ModelStart is the alias for RelativeOffset + TheFile.Header.ModelDataListOffset.
	*/
	RelativeOffset           uint32
	Size                     uint32
	BoneDataOffset1          uint32 // Relative to ModelStart
	BoneDataOffset2          uint32 // Relative to ModelStart
	MaterialDataOffset       uint32 // Relative to ModelStart
	BoneDataOffset3          uint32 // Relative to ModelStart
	IKDataOffset             uint32 // Relative to ModelStart
	MorphHeaderListOffset    uint32 // Relative to ModelStart
	SkinningHeaderListOffset uint32 // Relative to ModelStart
	TristripHeaderListOffset uint32 // Relative to ModelStart
	Unknown1                 uint32
	VertexDataOffset         uint32 // Relative to ModelStart
	Unknown2                 uint32
	Unknown3                 [302]byte
	Unknown4                 uint16
	Unknown5                 uint16
	BoneDataCount            uint16
	MorphHeaderCount         uint16
	MaterialCount            uint16
	IKDataCount              uint16
	SkinningHeaderCount      uint16
	TristripGroupCount       uint16
	Unknown6                 uint16
	VertexCount              uint16
	Unknown7                 uint16
	Unknown8                 uint16
	Unknown9                 uint16
}

type ModelData struct {
	MaterialList              []Material // Size is ModelHeader.MaterialCount
	BoneList                  []Bone     // Size is ModelHeader.BoneDataCount
	IKList                    []IK       // Size is ModelHeader.IKDataCount
	Filler1                   []byte
	MorphHeaderList           []MorphHeader           // Size is ModelHeader.MorphHeaderCount
	MorphDataContainerList    []MorphDataContainer    // Size is ModelHeader.MorphHeaderCount
	SkinningHeaderList        []SkinningHeader        // Size is ModelHeader.SkinningHeaderCount
	SkinningDataContainerList []SkinningDataContainer // Size is ModelHeader.SkinningHeaderCount
	TristripHeaderList        []TristripHeader        // Size is ModelHeader.TristripGroupCount
	TristripDataContainerList []TristripDataContainer // Size is ModelHeader.TristripGroupCount
	VertexDataList1           []VertexData            // Size ModelHeader.VertexCount
	VertexDataList2           []VertexData            // Size ModelHeader.VertexCount
}

type Material struct { // Size is 32
	MainTextureName [4]byte
	TextureType1    [4]byte
	TextureType2    [4]byte
	TextureType3    [4]byte
	TextureType4    [4]byte
	Unknown1        [3]float32
}

type Bone struct { // Size is 84
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

type IK struct { // Size is 16
	Location [3]float32
	Unknown1 uint32
}

type MorphHeader struct { // Size is 8
	Count      uint32
	DataOffset uint32 // Relative to ModelStart
}

type MorphDataContainer struct {
	Data []MorphData // Size is MorphHeader.Count
}

type MorphData struct { // Size is 16
	Unknown1 [4]byte
	Unknown2 [2]float32
	Unknown4 [4]byte
}

type SkinningHeader struct { // Size is 12
	Count      uint32
	ListOffset uint32
	Unknown1   uint16
	Unknown2   uint16
}

type SkinningDataContainer struct {
	Data []SkinningData // Size is SkinningHeader.Count
}

type SkinningData struct {
	BoneWeightPercentage uint16
	BoneID               uint16
}

type TristripHeader struct { // Size is 16
	IndexListOffset uint32
	IndexCount      uint16
	MaterialID      [5]uint16
}

type TristripDataContainer struct {
	Data []TristripData // Size is TristripHeader.IndexCount
}

type TristripData struct {
	Data []uint16 // Size is TristripHeader.IndexCount
}

type VertexData struct { // Size is 64
	LocationX float32
	LocationY float32
	LocationZ float32
	Unknown1  float32
	NormalX   float32
	NormalY   float32
	NormalZ   float32
	Unknown2  uint32
	Unknown3  [4]float32
	UVMapU    float32
	UVMapV    float32
	FFFFFFFF  uint32
	Unknown4  uint32
}
