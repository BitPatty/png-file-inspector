package pngscanner

type HeaderValue uint32

// Critical
const H_IHDR HeaderValue = 0x49_48_44_52 // IHDR Image Header
const H_PLTE HeaderValue = 0x50_4C_54_45 // PLTE Palette
const H_IDAT HeaderValue = 0x49_44_41_54 // IDAT Image Data
const H_IEND HeaderValue = 0x49_45_4E_44 // IEND Image Trailer

// Transparency Information
const H_TRNS HeaderValue = 0x74_52_4E_53 // tRNS Transparency Information

// Colour Space Information
const H_CHRM HeaderValue = 0x63_48_52_4D // cHRM Primary chromaticities and white point
const H_GAMA HeaderValue = 0x67_41_4D_41 // gAMA Image Gamma
const H_ICCP HeaderValue = 0x69_43_43_50 // iCCP Embedded ICC Profile
const H_SBIT HeaderValue = 0x73_42_49_54 // sBIT Significant Bits
const H_SRGB HeaderValue = 0x73_52_47_42 // sRGB Standard RGB colour space

// Textual Information
const H_TEXT HeaderValue = 0x74_45_58_74 // tEXt Textual Data
const H_ITXT HeaderValue = 0x69_54_58_74 // iTXt International Textual Data
const H_ZTXT HeaderValue = 0x7A_54_58_74 // zTXt Compressed Textual Data

// Miscellaneous Information
const H_BKGD HeaderValue = 0x62_4B_47_44 // bKGD Background Color
const H_HIST HeaderValue = 0x68_49_53_54 // hIST Image Histogram
const H_PHYS HeaderValue = 0x70_48_59_73 // pHYs Physical Pixel Dimensions
const H_SPLT HeaderValue = 0x73_50_4C_54 // sPLT Suggested Palette

// Time Stamp Information
const H_TIME HeaderValue = 0x74_49_4D_45 // Last Modification Time

var criticalHeaders = [4]HeaderValue{
	H_IHDR,
	H_PLTE,
	H_IDAT,
	H_IEND,
}

var ancillaryHeaders = [13]HeaderValue{
	H_TRNS,
	H_CHRM,
	H_GAMA,
	H_ICCP,
	H_SBIT,
	H_SRGB,
	H_TEXT,
	H_ITXT,
	H_ZTXT,
	H_BKGD,
	H_HIST,
	H_PHYS,
	H_SPLT,
}

type ChunkHeaderInspectionResult struct {
	header               *HeaderValue
	isAncillary          bool
	isPrivate            bool
	isReservedBitSet     bool
	isSafeToCopy         bool
	isStandardized       bool
	hasInvalidCharacters bool
}

func (h *HeaderValue) inspect() ChunkHeaderInspectionResult {
	return ChunkHeaderInspectionResult{
		header:               h,
		isAncillary:          h.isAncillaryBitSet(),
		isPrivate:            h.isPrivateBitSet(),
		isReservedBitSet:     h.isReservedBitSet(),
		isSafeToCopy:         h.isSafeToCopyBitSet(),
		isStandardized:       h.isStandardized(),
		hasInvalidCharacters: h.hasInvalidCharacters(),
	}
}

func (h *HeaderValue) isAncillaryBitSet() bool {
	return (*h & 0x20_00_00_00) > 0
}

func (h *HeaderValue) isPrivateBitSet() bool {
	return (*h & 0x00_20_00_00) > 0
}

func (h *HeaderValue) isReservedBitSet() bool {
	return (*h & 0x00_00_20_00) > 0
}

func (h *HeaderValue) isSafeToCopyBitSet() bool {
	return (*h & 0x00_00_00_20) > 0
}

func (h *HeaderValue) isStandardized() bool {
	isAncillary := h.isAncillaryBitSet()

	if !isAncillary {
		for _, v := range criticalHeaders {
			if v == *h {
				return true
			}
		}

		return false
	}

	for _, v := range ancillaryHeaders {
		if v == *h {
			return true
		}
	}

	return false
}

func (h *HeaderValue) hasInvalidCharacters() bool {
	for i := 0; i < 4; i++ {
		b := uint32(*h >> (i * 8) & 0xFF)

		if b < 0x41 || (b > 0x5A && b < 0x61) || b > 0x7A {
			return true
		}
	}

	return false
}
