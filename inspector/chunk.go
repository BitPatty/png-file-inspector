package pngscanner

import (
	"encoding/binary"
)

type ChunkData []uint8

type Chunk struct {
	length uint32
	header HeaderValue
	data   ChunkData
	crc    uint32
}

type IHDRInspectionResult struct {
	width             uint32
	height            uint32
	bitDepth          uint8
	colourType        uint8
	compressionMethod uint8
	filterMethod      uint8
	interlaceMethod   uint8
}

func (c *ChunkData) inspectIHDRData() IHDRInspectionResult {
	return IHDRInspectionResult{
		width:             binary.BigEndian.Uint32((*c)[0:4]),
		height:            binary.BigEndian.Uint32((*c)[4:8]),
		bitDepth:          (*c)[8],
		colourType:        (*c)[9],
		compressionMethod: (*c)[10],
		filterMethod:      (*c)[11],
		interlaceMethod:   (*c)[12],
	}
}
