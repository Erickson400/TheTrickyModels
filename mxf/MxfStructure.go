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
	ModelDataListOffset   uint32
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
	RelativeOffset            uint32
	Size                      uint32
	BoneDataOffset1           uint32 // Relative to ModelStart
	BoneDataOffset2           uint32 // Relative to ModelStart
	materialDataOffset        uint32 // Relative to ModelStart
	BoneDataOffset3           uint32 // Relative to ModelStart
	IKDataOffset              uint32 // Relative to ModelStart
	SkinningHeaderListOffset1 uint32 // Relative to ModelStart
	SkinningHeaderListOffset2 uint32 // Relative to ModelStart
	TristripHeaderListOffset  uint32 // Relative to ModelStart
	Unknown1                  uint32
	VertexDataOffset          uint32 // Relative to ModelStart
	Unknown2                  uint32
	Unknown3                  [302]byte
	BoneDataCount             uint16
	Unknown4                  uint16
	MaterialCount             uint16
	IKDataCount               uint16
	SkinningHeaderCount       uint16
	TriStripGroupCount        uint16
	Unknown5                  uint16
	VertexCount               uint16
	Unknown6                  uint16
	Unknown7                  uint16
	Unknown8                  uint16
}

type ModelData struct {
	MaterialList       []Material       // Size is ModelHeader.MaterialCount
	BoneData           []Bone           // Size is ModelHeader.BoneDataCount
	IKData             []IK             // Size is ModelHeader.IKDataCount
	SkinningHeaderList []SkinningHeader // Size is ModelHeader.SkinningHeaderCount

}

// --------------
type TristripData struct {
	OffsetOfTriangleIndicies uint32 // Relative to ModelStart
	TriCount                 uint32
	Unknown1                 [8]byte
}

type TriangleIndicies struct {
	TriangleIndex1 uint16
	TriangleIndex2 uint16
	TriangleIndex3 uint16
}

type VertexData struct {
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

type BoneData struct {
	Name                  [16]byte
	Unknown1              uint16
	ParentBoneID          uint16
	Unknown2              uint16
	BoneID                uint16
	LocationX             float32 // Location and Rotation cord is relative to parent bone.
	LocationY             float32
	LocationZ             float32
	RotationEulerRadianX1 float32
	RotationEulerRadianY1 float32
	RotationEulerRadianZ1 float32
	RotationEulerRadianX2 float32
	RotationEulerRadianY2 float32
	RotationEulerRadianZ2 float32
	Unknown3              [6]float32
}
