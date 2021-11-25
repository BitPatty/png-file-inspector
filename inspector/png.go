package pngscanner

// https://www.w3.org/TR/2003/REC-PNG-20031110

import (
	"encoding/binary"
	"fmt"
)

const (
	SIGNATURE         uint64 = 0x89_50_4E_47_0D_0A_1A_0A
	SIGNATURE_SIZE    int    = 8
	CHUNK_HEADER_SIZE int    = 4
	CHUNK_LENGTH_SIZE int    = 4
	CHUNK_CRC_SIZE    int    = 4
)

type InspectionResult struct {
	HasLeadingIHDR    bool
	HasValidSignature bool
	HasDataAfterIEND  bool
	Chunks            []ChunkInspectionResult
}

func Inspect(bytes []uint8) (InspectionResult, error) {
	if !validateSignature(&bytes) {
		return InspectionResult{
			HasValidSignature: false,
		}, fmt.Errorf("invalid PNG signature")
	}

	res := InspectionResult{
		HasValidSignature: true,
		Chunks:            []ChunkInspectionResult{},
	}

	for i := SIGNATURE_SIZE; i < len(bytes); {
		lastChunk, err := res.processChunk(&bytes, i)

		if err != nil {
			return res, err
		}

		nextChunkOffset := i + int(lastChunk.Length) + CHUNK_HEADER_SIZE + CHUNK_LENGTH_SIZE + CHUNK_CRC_SIZE

		if len(res.Chunks) == 1 {
			res.HasLeadingIHDR = lastChunk.Header == H_IHDR
		} else if lastChunk.Header == H_IEND {
			res.HasDataAfterIEND = nextChunkOffset < len(bytes)
			return res, nil
		}

		i = nextChunkOffset
	}

	return res, fmt.Errorf("Missing critical data")
}

func (res *InspectionResult) processChunk(bytes *[]uint8, offset int) (Chunk, error) {
	c, err := loadNextChunk(bytes, offset)

	if err != nil {
		return c, err
	}

	res.Chunks = append(res.Chunks, c.inspect())
	return c, nil
}

func validateSignature(bytes *[]uint8) bool {
	if len(*bytes) < 8 {
		return false
	}

	return binary.BigEndian.Uint64((*bytes)[:SIGNATURE_SIZE]) == SIGNATURE
}

func loadNextChunk(bytes *[]uint8, offset int) (Chunk, error) {
	c := Chunk{}

	if (offset + 12) > len(*bytes) {
		return c, fmt.Errorf("cannot load chunk, metadata length exceeds size")
	}

	c.Length = binary.BigEndian.Uint32((*bytes)[(offset):(offset + CHUNK_LENGTH_SIZE)])
	c.RawHeader = (*bytes)[(offset + CHUNK_LENGTH_SIZE):(offset + CHUNK_LENGTH_SIZE + CHUNK_HEADER_SIZE)]
	c.Header = HeaderValue(binary.BigEndian.Uint32(c.RawHeader[:]))

	if (offset + 12 + int(c.Length)) > len(*bytes) {
		return c, fmt.Errorf("cannot load chunk, data length exceeds size (%v at position %v)", c.Header.ToString(), offset)
	}

	c.Data = ChunkData((*bytes)[(offset + 8):(offset + 8 + int(c.Length))])
	c.CRC = binary.BigEndian.Uint32((*bytes)[(offset + 8 + int(c.Length)):(offset + 8 + int(c.Length) + 4)])

	return c, nil
}
