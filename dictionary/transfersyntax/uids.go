package transfersyntax

type TransferSyntax struct {
	UID         string
	Name        string
	Description string
	Type        string
}

var supportedTransferSyntaxes = []*TransferSyntax{
	ImplicitVRLittleEndian,
	ExplicitVRLittleEndian,
	DeflatedExplicitVRLittleEndian,
	ExplicitVRBigEndian,
	JPEGLosslessSV1,
	JPEGBaseline8Bit,
	JPEGExtended12Bit,
	JPEG2000Lossless,
	JPEG2000,
}

func GetTransferSyntaxFromName(name string) *TransferSyntax {
	for _, ts := range transferSyntaxes {
		if ts.Name == name {
			return ts
		}
	}
	return nil
}

func GetTransferSyntaxFromUID(uid string) *TransferSyntax {
	for _, ts := range transferSyntaxes {
		if ts.UID == uid {
			return ts
		}
	}
	// Extra loop to fix old bug
	uid = string([]rune(uid)[:len(uid)-1])
	for _, ts := range transferSyntaxes {
		if ts.UID == uid {
			return ts
		}
	}
	return nil
}

func SupportedTransferSyntax(uid string) bool {
	for _, ts := range supportedTransferSyntaxes {
		if ts.UID == uid {
			return true
		}
	}
	return false
}
