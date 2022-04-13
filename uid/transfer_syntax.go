package uid

func GetTransferSyntaxFromName(name string) *SOPClass {
	for _, sop := range TransferSyntaxes {
		if sop.Name == name {
			return sop
		}
	}
	return nil
}

func GetTransferSyntaxFromUID(uid string) *SOPClass {
	for _, sop := range TransferSyntaxes {
		if sop.UID == uid {
			return sop
		}
	}
	return nil
}

// ImplicitVRLittleEndian - (1.2.840.10008.1.2) Implicit VR - Little Endian
const ImplicitVRLittleEndian = "1.2.840.10008.1.2"

// ExplicitVRLittleEndian - (1.2.840.10008.1.2.1) Explicit VR - Little Endian
const ExplicitVRLittleEndian = "1.2.840.10008.1.2.1"

// DeflatedExplicitVRLittleEndian - (1.2.840.10008.1.2.1.99) Deflated Explicit VR Little Endian
const DeflatedExplicitVRLittleEndian = "1.2.840.10008.1.2.1.99"

// ExplicitVRBigEndian - (1.2.840.10008.1.2.2) Explicit VR Big Endian (Retired)
const ExplicitVRBigEndian = "1.2.840.10008.1.2.2"

// RLELossless - (1.2.840.10008.1.2.5) RLE (Run Length Encoding) Lossless
const RLELossless = "1.2.840.10008.1.2.5"

// JPEGBaseline1 - (1.2.840.10008.1.2.4.50) JPEG Baseline (Process 1): Default Transfer Syntax for Lossy JPEG 8 Bit Image Compression
const JPEGBaseline1 = "1.2.840.10008.1.2.4.50"

// JPEGExtended24 - (1.2.840.10008.1.2.4.51) JPEG Extended (Process 2 & 4): Default Transfer Syntax for Lossy JPEG 12 Bit Image Compression (Process 4 only)
const JPEGExtended24 = "1.2.840.10008.1.2.4.51"

// JPEGExtended35Retired - (1.2.840.10008.1.2.4.52) JPEG Extended (Process 3 & 5) (Retired)
const JPEGExtended35Retired = "1.2.840.10008.1.2.4.52"

// JPEGLosslessNonHierarchical14 - (1.2.840.10008.1.2.4.57) JPEG Lossless, Non-Hierarchical (Process 14)
const JPEGLosslessNonHierarchical14 = "1.2.840.10008.1.2.4.57"

// JPEGLossless - (1.2.840.10008.1.2.4.70) JPEG Lossless, Non-Hierarchical, First-Order Prediction (Process 14 [Selection Value 1])
const JPEGLossless = "1.2.840.10008.1.2.4.70"

// JPEGLSLossless - (1.2.840.10008.1.2.4.80) JPEG-LS Lossless Image Compression
const JPEGLSLossless = "1.2.840.10008.1.2.4.80"

// JPEGLSLossyNearLossless - (1.2.840.10008.1.2.4.81) JPEG-LS Lossy (Near-Lossless) Image Compression
const JPEGLSLossyNearLossless = "1.2.840.10008.1.2.4.81"

// JPEG2000Lossless - (1.2.840.10008.1.2.4.90) JPEG 2000 Image Compression (Lossless Only)
const JPEG2000Lossless = "1.2.840.10008.1.2.4.90"

// JPEG2000Lossy - (1.2.840.10008.1.2.4.91) JPEG 2000 Image Compression
const JPEG2000Lossy = "1.2.840.10008.1.2.4.91"

// JPEG2000Part2MultiComponentLosslessOnly - (1.2.840.10008.1.2.4.92) JPEG 2000 Part 2 Multi-component Image Compression (Lossless Only)
const JPEG2000Part2MultiComponentLosslessOnly = "1.2.840.10008.1.2.4.92"

// JPEG2000Part2MultiComponent - (1.2.840.10008.1.2.4.93) JPEG 2000 Part 2 Multi-component Image Compression
const JPEG2000Part2MultiComponent = "1.2.840.10008.1.2.4.93"

// JPIPReferenced - (1.2.840.10008.1.2.4.94) JPIP Referenced
const JPIPReferenced = "1.2.840.10008.1.2.4.94"

// JPIPReferencedDeflate - (1.2.840.10008.1.2.4.95) JPIP Referenced Deflate
const JPIPReferencedDeflate = "1.2.840.10008.1.2.4.95"

var TransferSyntaxes = []*SOPClass{
	{
		UID:         "1.2.840.10008.1.2",
		Name:        "ImplicitVRLittleEndian",
		Description: "Implicit VR - Little Endian",
	},
	{
		UID:         "1.2.840.10008.1.2.1",
		Name:        "ExplicitVRLittleEndian",
		Description: "Explicit VR - Little Endian",
	},
	{
		UID:         "1.2.840.10008.1.2.1.99",
		Name:        "DeflatedExplicitVRLittleEndian",
		Description: "Deflated Explicit VR Little Endian",
	},
	{
		UID:         "1.2.840.10008.1.2.2",
		Name:        "ExplicitVRBigEndian",
		Description: "Explicit VR Big Endian (Retired)",
	},
	{
		UID:         "1.2.840.10008.1.2.5",
		Name:        "RLELossless",
		Description: "RLE (Run Length Encoding) Lossless",
	},
	{
		UID:         "1.2.840.10008.1.2.4.50",
		Name:        "JPEGBaseline1",
		Description: "JPEG Baseline (Process 1): Default Transfer Syntax for Lossy JPEG 8 Bit Image Compression",
	},
	{
		UID:         "1.2.840.10008.1.2.4.51",
		Name:        "JPEGExtended24",
		Description: "JPEG Extended (Process 2 & 4): Default Transfer Syntax for Lossy JPEG 12 Bit Image Compression (Process 4 only)",
	},
	{
		UID:         "1.2.840.10008.1.2.4.52",
		Name:        "JPEGExtended35Retired",
		Description: "JPEG Extended (Process 3 & 5) (Retired)",
	},
	{
		UID:         "1.2.840.10008.1.2.4.57",
		Name:        "JPEGLosslessNonHierarchical14",
		Description: "JPEG Lossless, Non-Hierarchical (Process 14)",
	},
	{
		UID:         "1.2.840.10008.1.2.4.70",
		Name:        "JPEGLossless",
		Description: "JPEG Lossless, Non-Hierarchical, First-Order Prediction (Process 14 [Selection Value 1])",
	},
	{
		UID:         "1.2.840.10008.1.2.4.80",
		Name:        "JPEGLSLossless",
		Description: "JPEG-LS Lossless Image Compression",
	},
	{
		UID:         "1.2.840.10008.1.2.4.81",
		Name:        "JPEGLSLossyNearLossless",
		Description: "JPEG-LS Lossy (Near-Lossless) Image Compression",
	},
	{
		UID:         "1.2.840.10008.1.2.4.90",
		Name:        "JPEG2000Lossless",
		Description: "JPEG 2000 Image Compression (Lossless Only)",
	},
	{
		UID:         "1.2.840.10008.1.2.4.91",
		Name:        "JPEG2000Lossy",
		Description: "JPEG 2000 Image Compression",
	},
	{
		UID:         "1.2.840.10008.1.2.4.92",
		Name:        "JPEG2000Part2MultiComponentLosslessOnly",
		Description: "JPEG 2000 Part 2 Multi-component Image Compression (Lossless Only)",
	},
	{
		UID:         "1.2.840.10008.1.2.4.93",
		Name:        "JPEG2000Part2MultiComponent",
		Description: "JPEG 2000 Part 2 Multi-component Image Compression",
	},
	{
		UID:         "1.2.840.10008.1.2.4.94",
		Name:        "JPIPReferenced",
		Description: "JPIP Referenced",
	},
	{
		UID:         "1.2.840.10008.1.2.4.95",
		Name:        "JPIPReferencedDeflate",
		Description: "JPIP Referenced Deflate",
	},
}
