package uid

// ImplicitVRLittleEndian - (1.2.840.10008.1.2) Implicit VR Little Endian: Default Transfer Syntax for DICOM
var ImplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2",
	Name:        "ImplicitVRLittleEndian",
	Description: "Implicit VR Little Endian: Default Transfer Syntax for DICOM",
	Type:        "Transfer Syntax",
}

// ExplicitVRLittleEndian - (1.2.840.10008.1.2.1) Explicit VR Little Endian
var ExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1",
	Name:        "ExplicitVRLittleEndian",
	Description: "Explicit VR Little Endian",
	Type:        "Transfer Syntax",
}

// EncapsulatedUncompressedExplicitVRLittleEndian - (1.2.840.10008.1.2.1.98) Encapsulated Uncompressed Explicit VR Little Endian
var EncapsulatedUncompressedExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1.98",
	Name:        "EncapsulatedUncompressedExplicitVRLittleEndian",
	Description: "Encapsulated Uncompressed Explicit VR Little Endian",
	Type:        "Transfer Syntax",
}

// DeflatedExplicitVRLittleEndian - (1.2.840.10008.1.2.1.99) Deflated Explicit VR Little Endian
var DeflatedExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1.99",
	Name:        "DeflatedExplicitVRLittleEndian",
	Description: "Deflated Explicit VR Little Endian",
	Type:        "Transfer Syntax",
}

// ExplicitVRBigEndian - (1.2.840.10008.1.2.2) Explicit VR Big Endian (Retired)
var ExplicitVRBigEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.2",
	Name:        "ExplicitVRBigEndian",
	Description: "Explicit VR Big Endian (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGBaseline8Bit - (1.2.840.10008.1.2.4.50) JPEG Baseline (Process 1): Default Transfer Syntax for Lossy JPEG 8 Bit Image Compression
var JPEGBaseline8Bit = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.50",
	Name:        "JPEGBaseline8Bit",
	Description: "JPEG Baseline (Process 1): Default Transfer Syntax for Lossy JPEG 8 Bit Image Compression",
	Type:        "Transfer Syntax",
}

// JPEGExtended12Bit - (1.2.840.10008.1.2.4.51) JPEG Extended (Process 2 & 4): Default Transfer Syntax for Lossy JPEG 12 Bit Image Compression (Process 4 only)
var JPEGExtended12Bit = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.51",
	Name:        "JPEGExtended12Bit",
	Description: "JPEG Extended (Process 2 & 4): Default Transfer Syntax for Lossy JPEG 12 Bit Image Compression (Process 4 only)",
	Type:        "Transfer Syntax",
}

// JPEGExtended35 - (1.2.840.10008.1.2.4.52) JPEG Extended (Process 3 & 5) (Retired)
var JPEGExtended35 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.52",
	Name:        "JPEGExtended35",
	Description: "JPEG Extended (Process 3 & 5) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGSpectralSelectionNonHierarchical68 - (1.2.840.10008.1.2.4.53) JPEG Spectral Selection, Non-Hierarchical (Process 6 & 8) (Retired)
var JPEGSpectralSelectionNonHierarchical68 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.53",
	Name:        "JPEGSpectralSelectionNonHierarchical68",
	Description: "JPEG Spectral Selection, Non-Hierarchical (Process 6 & 8) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGSpectralSelectionNonHierarchical79 - (1.2.840.10008.1.2.4.54) JPEG Spectral Selection, Non-Hierarchical (Process 7 & 9) (Retired)
var JPEGSpectralSelectionNonHierarchical79 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.54",
	Name:        "JPEGSpectralSelectionNonHierarchical79",
	Description: "JPEG Spectral Selection, Non-Hierarchical (Process 7 & 9) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGFullProgressionNonHierarchical1012 - (1.2.840.10008.1.2.4.55) JPEG Full Progression, Non-Hierarchical (Process 10 & 12) (Retired)
var JPEGFullProgressionNonHierarchical1012 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.55",
	Name:        "JPEGFullProgressionNonHierarchical1012",
	Description: "JPEG Full Progression, Non-Hierarchical (Process 10 & 12) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGFullProgressionNonHierarchical1113 - (1.2.840.10008.1.2.4.56) JPEG Full Progression, Non-Hierarchical (Process 11 & 13) (Retired)
var JPEGFullProgressionNonHierarchical1113 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.56",
	Name:        "JPEGFullProgressionNonHierarchical1113",
	Description: "JPEG Full Progression, Non-Hierarchical (Process 11 & 13) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGLossless - (1.2.840.10008.1.2.4.57) JPEG Lossless, Non-Hierarchical (Process 14)
var JPEGLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.57",
	Name:        "JPEGLossless",
	Description: "JPEG Lossless, Non-Hierarchical (Process 14)",
	Type:        "Transfer Syntax",
}

// JPEGLosslessNonHierarchical15 - (1.2.840.10008.1.2.4.58) JPEG Lossless, Non-Hierarchical (Process 15) (Retired)
var JPEGLosslessNonHierarchical15 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.58",
	Name:        "JPEGLosslessNonHierarchical15",
	Description: "JPEG Lossless, Non-Hierarchical (Process 15) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGExtendedHierarchical1618 - (1.2.840.10008.1.2.4.59) JPEG Extended, Hierarchical (Process 16 & 18) (Retired)
var JPEGExtendedHierarchical1618 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.59",
	Name:        "JPEGExtendedHierarchical1618",
	Description: "JPEG Extended, Hierarchical (Process 16 & 18) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGExtendedHierarchical1719 - (1.2.840.10008.1.2.4.60) JPEG Extended, Hierarchical (Process 17 & 19) (Retired)
var JPEGExtendedHierarchical1719 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.60",
	Name:        "JPEGExtendedHierarchical1719",
	Description: "JPEG Extended, Hierarchical (Process 17 & 19) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGSpectralSelectionHierarchical2022 - (1.2.840.10008.1.2.4.61) JPEG Spectral Selection, Hierarchical (Process 20 & 22) (Retired)
var JPEGSpectralSelectionHierarchical2022 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.61",
	Name:        "JPEGSpectralSelectionHierarchical2022",
	Description: "JPEG Spectral Selection, Hierarchical (Process 20 & 22) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGSpectralSelectionHierarchical2123 - (1.2.840.10008.1.2.4.62) JPEG Spectral Selection, Hierarchical (Process 21 & 23) (Retired)
var JPEGSpectralSelectionHierarchical2123 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.62",
	Name:        "JPEGSpectralSelectionHierarchical2123",
	Description: "JPEG Spectral Selection, Hierarchical (Process 21 & 23) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGFullProgressionHierarchical2426 - (1.2.840.10008.1.2.4.63) JPEG Full Progression, Hierarchical (Process 24 & 26) (Retired)
var JPEGFullProgressionHierarchical2426 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.63",
	Name:        "JPEGFullProgressionHierarchical2426",
	Description: "JPEG Full Progression, Hierarchical (Process 24 & 26) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGFullProgressionHierarchical2527 - (1.2.840.10008.1.2.4.64) JPEG Full Progression, Hierarchical (Process 25 & 27) (Retired)
var JPEGFullProgressionHierarchical2527 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.64",
	Name:        "JPEGFullProgressionHierarchical2527",
	Description: "JPEG Full Progression, Hierarchical (Process 25 & 27) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGLosslessHierarchical28 - (1.2.840.10008.1.2.4.65) JPEG Lossless, Hierarchical (Process 28) (Retired)
var JPEGLosslessHierarchical28 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.65",
	Name:        "JPEGLosslessHierarchical28",
	Description: "JPEG Lossless, Hierarchical (Process 28) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGLosslessHierarchical29 - (1.2.840.10008.1.2.4.66) JPEG Lossless, Hierarchical (Process 29) (Retired)
var JPEGLosslessHierarchical29 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.66",
	Name:        "JPEGLosslessHierarchical29",
	Description: "JPEG Lossless, Hierarchical (Process 29) (Retired)",
	Type:        "Transfer Syntax",
}

// JPEGLosslessSV1 - (1.2.840.10008.1.2.4.70) JPEG Lossless, Non-Hierarchical, First-Order Prediction (Process 14 [Selection Value 1]): Default Transfer Syntax for Lossless JPEG Image Compression
var JPEGLosslessSV1 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.70",
	Name:        "JPEGLosslessSV1",
	Description: "JPEG Lossless, Non-Hierarchical, First-Order Prediction (Process 14 [Selection Value 1]): Default Transfer Syntax for Lossless JPEG Image Compression",
	Type:        "Transfer Syntax",
}

// JPEGLSLossless - (1.2.840.10008.1.2.4.80) JPEG-LS Lossless Image Compression
var JPEGLSLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.80",
	Name:        "JPEGLSLossless",
	Description: "JPEG-LS Lossless Image Compression",
	Type:        "Transfer Syntax",
}

// JPEGLSNearLossless - (1.2.840.10008.1.2.4.81) JPEG-LS Lossy (Near-Lossless) Image Compression
var JPEGLSNearLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.81",
	Name:        "JPEGLSNearLossless",
	Description: "JPEG-LS Lossy (Near-Lossless) Image Compression",
	Type:        "Transfer Syntax",
}

// JPEG2000Lossless - (1.2.840.10008.1.2.4.90) JPEG 2000 Image Compression (Lossless Only)
var JPEG2000Lossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.90",
	Name:        "JPEG2000Lossless",
	Description: "JPEG 2000 Image Compression (Lossless Only)",
	Type:        "Transfer Syntax",
}

// JPEG2000 - (1.2.840.10008.1.2.4.91) JPEG 2000 Image Compression
var JPEG2000 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.91",
	Name:        "JPEG2000",
	Description: "JPEG 2000 Image Compression",
	Type:        "Transfer Syntax",
}

// JPEG2000MCLossless - (1.2.840.10008.1.2.4.92) JPEG 2000 Part 2 Multi-component Image Compression (Lossless Only)
var JPEG2000MCLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.92",
	Name:        "JPEG2000MCLossless",
	Description: "JPEG 2000 Part 2 Multi-component Image Compression (Lossless Only)",
	Type:        "Transfer Syntax",
}

// JPEG2000MC - (1.2.840.10008.1.2.4.93) JPEG 2000 Part 2 Multi-component Image Compression
var JPEG2000MC = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.93",
	Name:        "JPEG2000MC",
	Description: "JPEG 2000 Part 2 Multi-component Image Compression",
	Type:        "Transfer Syntax",
}

// JPIPReferenced - (1.2.840.10008.1.2.4.94) JPIP Referenced
var JPIPReferenced = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.94",
	Name:        "JPIPReferenced",
	Description: "JPIP Referenced",
	Type:        "Transfer Syntax",
}

// JPIPReferencedDeflate - (1.2.840.10008.1.2.4.95) JPIP Referenced Deflate
var JPIPReferencedDeflate = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.95",
	Name:        "JPIPReferencedDeflate",
	Description: "JPIP Referenced Deflate",
	Type:        "Transfer Syntax",
}

// MPEG2MPML - (1.2.840.10008.1.2.4.100) MPEG2 Main Profile / Main Level
var MPEG2MPML = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.100",
	Name:        "MPEG2MPML",
	Description: "MPEG2 Main Profile / Main Level",
	Type:        "Transfer Syntax",
}

// MPEG2MPHL - (1.2.840.10008.1.2.4.101) MPEG2 Main Profile / High Level
var MPEG2MPHL = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.101",
	Name:        "MPEG2MPHL",
	Description: "MPEG2 Main Profile / High Level",
	Type:        "Transfer Syntax",
}

// MPEG4HP41 - (1.2.840.10008.1.2.4.102) MPEG-4 AVC/H.264 High Profile / Level 4.1
var MPEG4HP41 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.102",
	Name:        "MPEG4HP41",
	Description: "MPEG-4 AVC/H.264 High Profile / Level 4.1",
	Type:        "Transfer Syntax",
}

// MPEG4HP41BD - (1.2.840.10008.1.2.4.103) MPEG-4 AVC/H.264 BD-compatible High Profile / Level 4.1
var MPEG4HP41BD = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.103",
	Name:        "MPEG4HP41BD",
	Description: "MPEG-4 AVC/H.264 BD-compatible High Profile / Level 4.1",
	Type:        "Transfer Syntax",
}

// MPEG4HP422D - (1.2.840.10008.1.2.4.104) MPEG-4 AVC/H.264 High Profile / Level 4.2 For 2D Video
var MPEG4HP422D = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.104",
	Name:        "MPEG4HP422D",
	Description: "MPEG-4 AVC/H.264 High Profile / Level 4.2 For 2D Video",
	Type:        "Transfer Syntax",
}

// MPEG4HP423D - (1.2.840.10008.1.2.4.105) MPEG-4 AVC/H.264 High Profile / Level 4.2 For 3D Video
var MPEG4HP423D = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.105",
	Name:        "MPEG4HP423D",
	Description: "MPEG-4 AVC/H.264 High Profile / Level 4.2 For 3D Video",
	Type:        "Transfer Syntax",
}

// MPEG4HP42STEREO - (1.2.840.10008.1.2.4.106) MPEG-4 AVC/H.264 Stereo High Profile / Level 4.2
var MPEG4HP42STEREO = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.106",
	Name:        "MPEG4HP42STEREO",
	Description: "MPEG-4 AVC/H.264 Stereo High Profile / Level 4.2",
	Type:        "Transfer Syntax",
}

// HEVCMP51 - (1.2.840.10008.1.2.4.107) HEVC/H.265 Main Profile / Level 5.1
var HEVCMP51 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.107",
	Name:        "HEVCMP51",
	Description: "HEVC/H.265 Main Profile / Level 5.1",
	Type:        "Transfer Syntax",
}

// HEVCM10P51 - (1.2.840.10008.1.2.4.108) HEVC/H.265 Main 10 Profile / Level 5.1
var HEVCM10P51 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.108",
	Name:        "HEVCM10P51",
	Description: "HEVC/H.265 Main 10 Profile / Level 5.1",
	Type:        "Transfer Syntax",
}

// RLELossless - (1.2.840.10008.1.2.5) RLE Lossless
var RLELossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.5",
	Name:        "RLELossless",
	Description: "RLE Lossless",
	Type:        "Transfer Syntax",
}

// RFC2557MIMEEncapsulation - (1.2.840.10008.1.2.6.1) RFC 2557 MIME encapsulation (Retired)
var RFC2557MIMEEncapsulation = &SOPClass{
	UID:         "1.2.840.10008.1.2.6.1",
	Name:        "RFC2557MIMEEncapsulation",
	Description: "RFC 2557 MIME encapsulation (Retired)",
	Type:        "Transfer Syntax",
}

// XMLEncoding - (1.2.840.10008.1.2.6.2) XML Encoding (Retired)
var XMLEncoding = &SOPClass{
	UID:         "1.2.840.10008.1.2.6.2",
	Name:        "XMLEncoding",
	Description: "XML Encoding (Retired)",
	Type:        "Transfer Syntax",
}

// SMPTEST211020UncompressedProgressiveActiveVideo - (1.2.840.10008.1.2.7.1) SMPTE ST 2110-20 Uncompressed Progressive Active Video
var SMPTEST211020UncompressedProgressiveActiveVideo = &SOPClass{
	UID:         "1.2.840.10008.1.2.7.1",
	Name:        "SMPTEST211020UncompressedProgressiveActiveVideo",
	Description: "SMPTE ST 2110-20 Uncompressed Progressive Active Video",
	Type:        "Transfer Syntax",
}

// SMPTEST211020UncompressedInterlacedActiveVideo - (1.2.840.10008.1.2.7.2) SMPTE ST 2110-20 Uncompressed Interlaced Active Video
var SMPTEST211020UncompressedInterlacedActiveVideo = &SOPClass{
	UID:         "1.2.840.10008.1.2.7.2",
	Name:        "SMPTEST211020UncompressedInterlacedActiveVideo",
	Description: "SMPTE ST 2110-20 Uncompressed Interlaced Active Video",
	Type:        "Transfer Syntax",
}

// SMPTEST211030PCMDigitalAudio - (1.2.840.10008.1.2.7.3) SMPTE ST 2110-30 PCM Digital Audio
var SMPTEST211030PCMDigitalAudio = &SOPClass{
	UID:         "1.2.840.10008.1.2.7.3",
	Name:        "SMPTEST211030PCMDigitalAudio",
	Description: "SMPTE ST 2110-30 PCM Digital Audio",
	Type:        "Transfer Syntax",
}

// Papyrus3ImplicitVRLittleEndian - (1.2.840.10008.1.20) Papyrus 3 Implicit VR Little Endian (Retired)
var Papyrus3ImplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.20",
	Name:        "Papyrus3ImplicitVRLittleEndian",
	Description: "Papyrus 3 Implicit VR Little Endian (Retired)",
	Type:        "Transfer Syntax",
}

var transferSyntaxes = []*SOPClass{
	ImplicitVRLittleEndian,
	ExplicitVRLittleEndian,
	EncapsulatedUncompressedExplicitVRLittleEndian,
	DeflatedExplicitVRLittleEndian,
	ExplicitVRBigEndian,
	JPEGBaseline8Bit,
	JPEGExtended12Bit,
	JPEGExtended35,
	JPEGSpectralSelectionNonHierarchical68,
	JPEGSpectralSelectionNonHierarchical79,
	JPEGFullProgressionNonHierarchical1012,
	JPEGFullProgressionNonHierarchical1113,
	JPEGLossless,
	JPEGLosslessNonHierarchical15,
	JPEGExtendedHierarchical1618,
	JPEGExtendedHierarchical1719,
	JPEGSpectralSelectionHierarchical2022,
	JPEGSpectralSelectionHierarchical2123,
	JPEGFullProgressionHierarchical2426,
	JPEGFullProgressionHierarchical2527,
	JPEGLosslessHierarchical28,
	JPEGLosslessHierarchical29,
	JPEGLosslessSV1,
	JPEGLSLossless,
	JPEGLSNearLossless,
	JPEG2000Lossless,
	JPEG2000,
	JPEG2000MCLossless,
	JPEG2000MC,
	JPIPReferenced,
	JPIPReferencedDeflate,
	MPEG2MPML,
	MPEG2MPHL,
	MPEG4HP41,
	MPEG4HP41BD,
	MPEG4HP422D,
	MPEG4HP423D,
	MPEG4HP42STEREO,
	HEVCMP51,
	HEVCM10P51,
	RLELossless,
	RFC2557MIMEEncapsulation,
	XMLEncoding,
	SMPTEST211020UncompressedProgressiveActiveVideo,
	SMPTEST211020UncompressedInterlacedActiveVideo,
	SMPTEST211030PCMDigitalAudio,
	Papyrus3ImplicitVRLittleEndian,
}
