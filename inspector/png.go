package pngscanner

// https://www.w3.org/TR/2003/REC-PNG-20031110

import (
	"encoding/binary"
	"fmt"
)

const (
	SIGNATURE uint64 = 0x89_50_4E_47_0D_0A_1A_0A
)

type InspectionResult struct {
	metadata struct {
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
		return InspectionResult{}, fmt.Errorf("invalid PNG signature")
	}

	res := InspectionResult{
		chunks: []ChunkHeaderInspectionResult{},
	}

	for i := 8; i < len(bytes); {
		if (i + 12) > len(bytes) {
			return res, fmt.Errorf("cannot read next chunk, critical length exceeds image size")
		}

		nextChunkLength := binary.BigEndian.Uint32(bytes[i:(i + 4)])

		if (nextChunkLength & 0x80_00_00_00) > 0 {
			return res, fmt.Errorf("cannot read next chunk, length is greater than max allowed length")
		}

		if (i + 12 + int(nextChunkLength)) > len(bytes) {
			return res, fmt.Errorf("cannot read next chunk, data length exceeds image size")
		}

		headerValue := HeaderValue(binary.BigEndian.Uint32(bytes[(i + 4):(i + 8)]))

		chunk := Chunk{
			length: nextChunkLength,
			header: headerValue,
			data:   ChunkData(bytes[(i + 8):(i + 8 + int(nextChunkLength))]),
			crc:    binary.BigEndian.Uint32(bytes[(i + 8 + int(nextChunkLength)):(i + 8 + int(nextChunkLength) + 4)]),
		}

		headerData := chunk.header.inspect()

		if !headerData.isStandardized {
			if headerData.isAncillary && !options.AllowUnknownAncillaryChunks {
				return res, fmt.Errorf("unknown ancillary chunk %v", headerData.header)
			} else if !headerData.isAncillary && !options.AllowUnknownCriticalChunks {
				return res, fmt.Errorf("unknown critical chunk %v", headerData.header)
			}
		}

		res.chunks = append(res.chunks, headerData)

		i += int(nextChunkLength) + 12

		switch chunk.header {
		case H_IHDR:
			res.metadata = struct{ IHDRInspectionResult }{
				chunk.data.inspectIHDRData(),
			}
		case H_IEND:
			if i < len(bytes) {
				return res, fmt.Errorf("found data after IEND")
			}
		}
	}

	return res, nil
}

func validateSignature(bytes *[]uint8) bool {
	if len(*bytes) < 8 {
		return false
	}

	return binary.BigEndian.Uint64((*bytes)[:8]) == SIGNATURE
}
