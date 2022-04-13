package uid

var ImplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2",
	Name:        "ImplicitVRLittleEndian",
	Description: "Implicit VR - Little Endian",
}

var ExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1",
	Name:        "ExplicitVRLittleEndian",
	Description: "Explicit VR - Little Endian",
}

var DeflatedExplicitVRLittleEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.1.99",
	Name:        "DeflatedExplicitVRLittleEndian",
	Description: "Deflated Explicit VR Little Endian",
}

var ExplicitVRBigEndian = &SOPClass{
	UID:         "1.2.840.10008.1.2.2",
	Name:        "ExplicitVRBigEndian",
	Description: "Explicit VR Big Endian (Retired)",
}

var RLELossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.5",
	Name:        "RLELossless",
	Description: "RLE (Run Length Encoding) Lossless",
}

var JPEGBaseline1 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.50",
	Name:        "JPEGBaseline1",
	Description: "JPEG Baseline (Process 1): Default Transfer Syntax for Lossy JPEG 8 Bit Image Compression",
}

var JPEGExtended24 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.51",
	Name:        "JPEGExtended24",
	Description: "JPEG Extended (Process 2 & 4): Default Transfer Syntax for Lossy JPEG 12 Bit Image Compression (Process 4 only)",
}

var JPEGExtended35Retired = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.52",
	Name:        "JPEGExtended35Retired",
	Description: "JPEG Extended (Process 3 & 5) (Retired)",
}

var JPEGLosslessNonHierarchical14 = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.57",
	Name:        "JPEGLosslessNonHierarchical14",
	Description: "JPEG Lossless, Non-Hierarchical (Process 14)",
}

var JPEGLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.70",
	Name:        "JPEGLossless",
	Description: "JPEG Lossless, Non-Hierarchical, First-Order Prediction (Process 14 [Selection Value 1])",
}

var JPEGLSLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.80",
	Name:        "JPEGLSLossless",
	Description: "JPEG-LS Lossless Image Compression",
}

var JPEGLSLossyNearLossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.81",
	Name:        "JPEGLSLossyNearLossless",
	Description: "JPEG-LS Lossy (Near-Lossless) Image Compression",
}

var JPEG2000Lossless = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.90",
	Name:        "JPEG2000Lossless",
	Description: "JPEG 2000 Image Compression (Lossless Only)",
}

var JPEG2000Lossy = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.91",
	Name:        "JPEG2000Lossy",
	Description: "JPEG 2000 Image Compression",
}

var JPEG2000Part2MultiComponentLosslessOnly = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.92",
	Name:        "JPEG2000Part2MultiComponentLosslessOnly",
	Description: "JPEG 2000 Part 2 Multi-component Image Compression (Lossless Only)",
}

var JPEG2000Part2MultiComponent = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.93",
	Name:        "JPEG2000Part2MultiComponent",
	Description: "JPEG 2000 Part 2 Multi-component Image Compression",
}

var JPIPReferenced = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.94",
	Name:        "JPIPReferenced",
	Description: "JPIP Referenced",
}

var JPIPReferencedDeflate = &SOPClass{
	UID:         "1.2.840.10008.1.2.4.95",
	Name:        "JPIPReferencedDeflate",
	Description: "JPIP Referenced Deflate",
}

var TransferSyntaxes = []*SOPClass{
	ImplicitVRLittleEndian,
	ExplicitVRLittleEndian,
	DeflatedExplicitVRLittleEndian,
	ExplicitVRBigEndian,
	RLELossless,
	JPEGBaseline1,
	JPEGExtended24,
	JPEGExtended35Retired,
	JPEGLosslessNonHierarchical14,
	JPEGLossless,
	JPEGLSLossless,
	JPEGLSLossyNearLossless,
	JPEG2000Lossless,
	JPEG2000Lossy,
	JPEG2000Part2MultiComponentLosslessOnly,
	JPEG2000Part2MultiComponent,
	JPIPReferenced,
	JPIPReferencedDeflate,
}
