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
	X        float32
	Y        float32
	Z        float32
	Unknown1 float32
	NormalX  float32
	NormalY  float32
	NormalZ  float32
	Unknown2 uint32
	Unknown3 [4]float32
	UVMapU   float32
	UVMapV   float32
	FFFFFFFF uint32
	Unknown4 uint32
}
