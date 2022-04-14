package uid

var ImplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2",
	Name:        "ImplicitVRLittleEndian",
	Description: "Implicit VR Little Endian: Default Transfer Syntax for DICOM",
	Type:        "Transfer Syntax",
}

var ExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1",
	Name:        "ExplicitVRLittleEndian",
	Description: "Explicit VR Little Endian",
	Type:        "Transfer Syntax",
}

var EncapsulatedUncompressedExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1.98",
	Name:        "EncapsulatedUncompressedExplicitVRLittleEndian",
	Description: "Encapsulated Uncompressed Explicit VR Little Endian",
	Type:        "Transfer Syntax",
}

var DeflatedExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1.99",
	Name:        "DeflatedExplicitVRLittleEndian",
	Description: "Deflated Explicit VR Little Endian",
	Type:        "Transfer Syntax",
}

var ExplicitVRBigEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.2",
	Name:        "ExplicitVRBigEndian",
	Description: "Explicit VR Big Endian (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGBaseline8Bit = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.50",
	Name:        "JPEGBaseline8Bit",
	Description: "JPEG Baseline (Process 1): Default Transfer Syntax for Lossy JPEG 8 Bit Image Compression",
	Type:        "Transfer Syntax",
}

var JPEGExtended12Bit = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.51",
	Name:        "JPEGExtended12Bit",
	Description: "JPEG Extended (Process 2 & 4): Default Transfer Syntax for Lossy JPEG 12 Bit Image Compression (Process 4 only)",
	Type:        "Transfer Syntax",
}

var JPEGExtended35 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.52",
	Name:        "JPEGExtended35",
	Description: "JPEG Extended (Process 3 & 5) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGSpectralSelectionNonHierarchical68 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.53",
	Name:        "JPEGSpectralSelectionNonHierarchical68",
	Description: "JPEG Spectral Selection, Non-Hierarchical (Process 6 & 8) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGSpectralSelectionNonHierarchical79 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.54",
	Name:        "JPEGSpectralSelectionNonHierarchical79",
	Description: "JPEG Spectral Selection, Non-Hierarchical (Process 7 & 9) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGFullProgressionNonHierarchical1012 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.55",
	Name:        "JPEGFullProgressionNonHierarchical1012",
	Description: "JPEG Full Progression, Non-Hierarchical (Process 10 & 12) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGFullProgressionNonHierarchical1113 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.56",
	Name:        "JPEGFullProgressionNonHierarchical1113",
	Description: "JPEG Full Progression, Non-Hierarchical (Process 11 & 13) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.57",
	Name:        "JPEGLossless",
	Description: "JPEG Lossless, Non-Hierarchical (Process 14)",
	Type:        "Transfer Syntax",
}

var JPEGLosslessNonHierarchical15 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.58",
	Name:        "JPEGLosslessNonHierarchical15",
	Description: "JPEG Lossless, Non-Hierarchical (Process 15) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGExtendedHierarchical1618 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.59",
	Name:        "JPEGExtendedHierarchical1618",
	Description: "JPEG Extended, Hierarchical (Process 16 & 18) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGExtendedHierarchical1719 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.60",
	Name:        "JPEGExtendedHierarchical1719",
	Description: "JPEG Extended, Hierarchical (Process 17 & 19) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGSpectralSelectionHierarchical2022 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.61",
	Name:        "JPEGSpectralSelectionHierarchical2022",
	Description: "JPEG Spectral Selection, Hierarchical (Process 20 & 22) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGSpectralSelectionHierarchical2123 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.62",
	Name:        "JPEGSpectralSelectionHierarchical2123",
	Description: "JPEG Spectral Selection, Hierarchical (Process 21 & 23) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGFullProgressionHierarchical2426 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.63",
	Name:        "JPEGFullProgressionHierarchical2426",
	Description: "JPEG Full Progression, Hierarchical (Process 24 & 26) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGFullProgressionHierarchical2527 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.64",
	Name:        "JPEGFullProgressionHierarchical2527",
	Description: "JPEG Full Progression, Hierarchical (Process 25 & 27) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGLosslessHierarchical28 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.65",
	Name:        "JPEGLosslessHierarchical28",
	Description: "JPEG Lossless, Hierarchical (Process 28) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGLosslessHierarchical29 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.66",
	Name:        "JPEGLosslessHierarchical29",
	Description: "JPEG Lossless, Hierarchical (Process 29) (Retired)",
	Type:        "Transfer Syntax",
}

var JPEGLosslessSV1 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.70",
	Name:        "JPEGLosslessSV1",
	Description: "JPEG Lossless, Non-Hierarchical, First-Order Prediction (Process 14 [Selection Value 1]): Default Transfer Syntax for Lossless JPEG Image Compression",
	Type:        "Transfer Syntax",
}

var JPEGLSLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.80",
	Name:        "JPEGLSLossless",
	Description: "JPEG-LS Lossless Image Compression",
	Type:        "Transfer Syntax",
}

var JPEGLSNearLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.81",
	Name:        "JPEGLSNearLossless",
	Description: "JPEG-LS Lossy (Near-Lossless) Image Compression",
	Type:        "Transfer Syntax",
}

var JPEG2000Lossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.90",
	Name:        "JPEG2000Lossless",
	Description: "JPEG 2000 Image Compression (Lossless Only)",
	Type:        "Transfer Syntax",
}

var JPEG2000 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.91",
	Name:        "JPEG2000",
	Description: "JPEG 2000 Image Compression",
	Type:        "Transfer Syntax",
}

var JPEG2000MCLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.92",
	Name:        "JPEG2000MCLossless",
	Description: "JPEG 2000 Part 2 Multi-component Image Compression (Lossless Only)",
	Type:        "Transfer Syntax",
}

var JPEG2000MC = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.93",
	Name:        "JPEG2000MC",
	Description: "JPEG 2000 Part 2 Multi-component Image Compression",
	Type:        "Transfer Syntax",
}

var JPIPReferenced = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.94",
	Name:        "JPIPReferenced",
	Description: "JPIP Referenced",
	Type:        "Transfer Syntax",
}

var JPIPReferencedDeflate = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.95",
	Name:        "JPIPReferencedDeflate",
	Description: "JPIP Referenced Deflate",
	Type:        "Transfer Syntax",
}

var MPEG2MPML = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.100",
	Name:        "MPEG2MPML",
	Description: "MPEG2 Main Profile / Main Level",
	Type:        "Transfer Syntax",
}

var MPEG2MPHL = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.101",
	Name:        "MPEG2MPHL",
	Description: "MPEG2 Main Profile / High Level",
	Type:        "Transfer Syntax",
}

var MPEG4HP41 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.102",
	Name:        "MPEG4HP41",
	Description: "MPEG-4 AVC/H.264 High Profile / Level 4.1",
	Type:        "Transfer Syntax",
}

var MPEG4HP41BD = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.103",
	Name:        "MPEG4HP41BD",
	Description: "MPEG-4 AVC/H.264 BD-compatible High Profile / Level 4.1",
	Type:        "Transfer Syntax",
}

var MPEG4HP422D = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.104",
	Name:        "MPEG4HP422D",
	Description: "MPEG-4 AVC/H.264 High Profile / Level 4.2 For 2D Video",
	Type:        "Transfer Syntax",
}

var MPEG4HP423D = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.105",
	Name:        "MPEG4HP423D",
	Description: "MPEG-4 AVC/H.264 High Profile / Level 4.2 For 3D Video",
	Type:        "Transfer Syntax",
}

var MPEG4HP42STEREO = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.106",
	Name:        "MPEG4HP42STEREO",
	Description: "MPEG-4 AVC/H.264 Stereo High Profile / Level 4.2",
	Type:        "Transfer Syntax",
}

var HEVCMP51 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.107",
	Name:        "HEVCMP51",
	Description: "HEVC/H.265 Main Profile / Level 5.1",
	Type:        "Transfer Syntax",
}

var HEVCM10P51 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.108",
	Name:        "HEVCM10P51",
	Description: "HEVC/H.265 Main 10 Profile / Level 5.1",
	Type:        "Transfer Syntax",
}

var RLELossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.5",
	Name:        "RLELossless",
	Description: "RLE Lossless",
	Type:        "Transfer Syntax",
}

var RFC2557MIMEEncapsulation = &SOPClass{
	UID:         "1.2.840.10008.1.2.6.1",
	Name:        "RFC2557MIMEEncapsulation",
	Description: "RFC 2557 MIME encapsulation (Retired)",
	Type:        "Transfer Syntax",
}

var XMLEncoding = &SOPClass{
	UID:         "1.2.840.10008.1.2.6.2",
	Name:        "XMLEncoding",
	Description: "XML Encoding (Retired)",
	Type:        "Transfer Syntax",
}

var SMPTEST211020UncompressedProgressiveActiveVideo = &SOPClass{
	UID:         "1.2.840.10008.1.2.7.1",
	Name:        "SMPTEST211020UncompressedProgressiveActiveVideo",
	Description: "SMPTE ST 2110-20 Uncompressed Progressive Active Video",
	Type:        "Transfer Syntax",
}

var SMPTEST211020UncompressedInterlacedActiveVideo = &SOPClass{
	UID:         "1.2.840.10008.1.2.7.2",
	Name:        "SMPTEST211020UncompressedInterlacedActiveVideo",
	Description: "SMPTE ST 2110-20 Uncompressed Interlaced Active Video",
	Type:        "Transfer Syntax",
}

var SMPTEST211030PCMDigitalAudio = &SOPClass{
	UID:         "1.2.840.10008.1.2.7.3",
	Name:        "SMPTEST211030PCMDigitalAudio",
	Description: "SMPTE ST 2110-30 PCM Digital Audio",
	Type:        "Transfer Syntax",
}

var Papyrus3ImplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.20",
	Name:        "Papyrus3ImplicitVRLittleEndian",
	Description: "Papyrus 3 Implicit VR Little Endian (Retired)",
	Type:        "Transfer Syntax",
}

var TransferSyntaxes = []*SOPClass{
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
