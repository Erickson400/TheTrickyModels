package main

type FileHeader struct {
	Unknown         uint32 // Always stores 4
	ModelCount      uint16
	ModelListOffset uint16 // Points to first ModelHeader
	ModelRootOffset uint32

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
	ModelRelativeOffset  uint32
	ModelSize            uint32
	OffsetOfBoneData1    uint32 // Relative to ModelStart
	OffsetOfBoneData2    uint32 // Relative to ModelStart
	Unknown1             uint32
	OffsetOfBoneData3    uint32 // Relative to ModelStart
	Unknown2             uint32
	OffsetOfModelData1   uint32 // Relative to ModelStart
	OffsetOfModelData2   uint32 // Relative to ModelStart
	OffsetOfTristripData uint32 // Relative to ModelStart
	Unknown3             uint32
	OffsetOfVertexData   uint32 // Relative to ModelStart
	Unknown4             uint32
	Unknown5             [332]byte
}

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
