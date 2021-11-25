package pngscanner

// +-----------------------+-------------+--------------------+---------------------------------------------------------------+
// |    PNG image type     | Colour type | Allowed bit depths |                        Interpretation                         |
// +-----------------------+-------------+--------------------+---------------------------------------------------------------+
// | Greyscale             |           0 |     1, 2, 4, 8, 16 | Each pixel is a greyscale sample                              |
// | Truecolour            |           2 |              8, 16 | Each pixel is an R,G,B triple                                 |
// | Indexed-colour        |           3 |         1, 2, 4, 8 | Each pixel is a palette index; a PLTE chunk shall appear.     |
// | Greyscale with alpha  |           4 |              8, 16 | Each pixel is a greyscale sample followed by an alpha sample. |
// | Truecolour with alpha |           6 |              8, 16 | Each pixel is an R,G,B triple followed by an alpha sample.    |
// +-----------------------+-------------+--------------------+---------------------------------------------------------------+

import (
	"encoding/binary"
)

const (
	COL_GREYSCALE       = 0b0000 // Each pixel is a greyscale sample
	COL_TRUECOLOUR      = 0b0010 // Each pixel is an RGB triple
	COL_INDEXED         = 0b0011 // Each pixel is a palette index
	COL_GREYSCALE_ALPHA = 0b0100 // Each pixel is a greyscale sample followed by an alpha sample
	COL_TRUECOLOR_ALPHA = 0b0110 // Each pixel is an RGB triple followed by an alpha sample
)

type IHDRValidationResult struct {
	HasValidWidth             bool
	HasValidHeight            bool
	HasValidColourType        bool
	HasValidBitDepth          bool
	HasValidCompressionMethod bool
	HasValidFilterMethod      bool
	HasValidInterlaceMethod   bool
}

type IHDRInspectionResult struct {
	Width             uint32
	Height            uint32
	BitDepth          uint8
	ColourType        uint8
	CompressionMethod uint8
	FilterMethod      uint8
	InterlaceMethod   uint8
	Validation        IHDRValidationResult
}

func inspectIHDRData(c *ChunkData) IHDRInspectionResult {
	res := IHDRInspectionResult{
		Width:             binary.BigEndian.Uint32((*c)[0:4]),
		Height:            binary.BigEndian.Uint32((*c)[4:8]),
		BitDepth:          (*c)[8],
		ColourType:        (*c)[9],
		CompressionMethod: (*c)[10],
		FilterMethod:      (*c)[11],
		InterlaceMethod:   (*c)[12],
	}

	res.Validation = IHDRValidationResult{
		HasValidWidth:             validateWidth(res.Width),
		HasValidHeight:            validateHeight(res.Height),
		HasValidColourType:        validateColourType(res.ColourType),
		HasValidBitDepth:          validateBitDepth(res.ColourType, res.BitDepth),
		HasValidCompressionMethod: validateCompressionMethod(res.CompressionMethod),
		HasValidFilterMethod:      validateFilterMethod(res.FilterMethod),
		HasValidInterlaceMethod:   validateInterlaceMethod(res.InterlaceMethod),
	}

	return res
}

func validateWidth(value uint32) bool {
	return value > 0
}

func validateHeight(value uint32) bool {
	return value > 0
}

func validateColourType(value uint8) bool {
	return ((value & 0b1100) == 0) || ((value & 0b1001) == 0)
}

func validateBitDepth(colourType uint8, value uint8) bool {
	switch value {
	case 0x1:
		fallthrough
	case 0x2:
		fallthrough
	case 0x4:
		return colourType == COL_GREYSCALE || colourType == COL_INDEXED
	case 0x8:
		return validateColourType(colourType)
	case 0x10:
		return colourType != COL_INDEXED && validateColourType(colourType)
	default:
		return false
	}
}

func validateCompressionMethod(compressionMethod uint8) bool {
	return compressionMethod == 0
}

func validateFilterMethod(filterMethod uint8) bool {
	return filterMethod == 0
}

func validateInterlaceMethod(interlaceMethod uint8) bool {
	return interlaceMethod < 2
}
