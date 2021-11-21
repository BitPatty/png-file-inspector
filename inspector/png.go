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
	CHUNK_CRC_SIZE    int    = 32
	CHUNK_MIN_SIZE    int    = CHUNK_LENGTH_SIZE + CHUNK_HEADER_SIZE + CHUNK_CRC_SIZE
)

type InspectionResult struct {
	hasValidSignature bool
	metadata          struct {
		IHDRInspectionResult
	}
	chunks []ChunkHeaderInspectionResult
}

type InspectOptions struct {
	AllowUnknownAncillaryChunks bool
	AllowUnknownCriticalChunks  bool
}

func Inspect(bytes []uint8, options InspectOptions) (InspectionResult, error) {
	if !validateSignature(&bytes) {
		return InspectionResult{
			hasValidSignature: false,
		}, fmt.Errorf("invalid PNG signature")
	}

	res := InspectionResult{
		chunks: []ChunkHeaderInspectionResult{},
	}

	for i := SIGNATURE_SIZE; i < len(bytes); {
		lastChunk, err := res.processChunk(&bytes, i, options)

		if err != nil {
			return res, err
		}

		nextChunkOffset := i + int(lastChunk.length) + CHUNK_MIN_SIZE

		if (len(res.chunks) == 1) && (lastChunk.header != H_IHDR) {
			return res, fmt.Errorf("first chunk must be IHDR")
		} else if lastChunk.header == H_IHDR {
			res.metadata = struct{ IHDRInspectionResult }{
				lastChunk.data.inspectIHDRData(),
			}
		} else if (lastChunk.header == H_IEND) && (nextChunkOffset < len(bytes)) {
			return res, fmt.Errorf("additional data present after IEND")
		}

		i = nextChunkOffset
	}

	return res, nil
}

func (res *InspectionResult) processChunk(bytes *[]uint8, offset int, options InspectOptions) (Chunk, error) {
	c, err := loadNextChunk(bytes, offset)

	if err != nil {
		return c, err
	}

	res.chunks = append(res.chunks, c.header.inspect())
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

	c.length = binary.BigEndian.Uint32((*bytes)[(offset):(offset + CHUNK_LENGTH_SIZE)])
	c.header = HeaderValue(binary.BigEndian.Uint32((*bytes)[(offset + CHUNK_LENGTH_SIZE):(offset + CHUNK_LENGTH_SIZE + CHUNK_HEADER_SIZE)]))

	if (offset + 12 + int(c.length)) > len(*bytes) {
		return c, fmt.Errorf("cannot load chunk, data length exceeds size")
	}

	c.data = ChunkData((*bytes)[(offset + 8):(offset + 8 + int(c.length))])
	c.crc = binary.BigEndian.Uint32((*bytes)[(offset + 8 + int(c.length)):(offset + 8 + int(c.length) + 4)])

	return c, nil
}

// func printHeaderString(h HeaderValue) {
// 	s := ""

// 	for i := 3; i >= 0; i-- {
// 		s += string(uint8((h >> (i * 8)) & 0xFF))
// 	}

// 	fmt.Println(s)
// }

// 	if !headerData.isStandardized {
// 		if headerData.isAncillary && !options.AllowUnknownAncillaryChunks {
// 			return c, fmt.Errorf("unknown ancillary chunk %v", headerData.header)
// 		} else if !headerData.isAncillary && !options.AllowUnknownCriticalChunks {
// 			return c, fmt.Errorf("unknown critical chunk %v", headerData.header)
// 		}
// 	}
