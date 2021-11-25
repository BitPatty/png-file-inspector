package pngscanner

import "hash/crc32"

type ChunkData []uint8

type Chunk struct {
	Length    uint32
	RawHeader []uint8
	Header    HeaderValue
	Data      ChunkData
	CRC       uint32
}

type ChunkInspectionResult struct {
	Header           ChunkHeaderInspectionResult
	Length           int
	HasValidChecksum bool
	Report           interface{}
}

func (c *Chunk) inspect() ChunkInspectionResult {
	res := ChunkInspectionResult{
		Header:           c.Header.inspect(),
		Length:           len(c.Data),
		HasValidChecksum: c.calculateCrc32() == c.CRC,
	}

	switch c.Header {
	case H_IHDR:
		res.Report = inspectIHDRData(&c.Data)
	case H_PLTE:
		res.Report = inspectPLTEData(&c.Data)
	case H_IEND:
		res.Report = inspectIENDData(&c.Data)
	default:
		res.Report = nil
	}

	return res
}

func (c *Chunk) calculateCrc32() uint32 {
	return crc32.ChecksumIEEE(append(c.RawHeader[0:4], []uint8(c.Data)...))
}
