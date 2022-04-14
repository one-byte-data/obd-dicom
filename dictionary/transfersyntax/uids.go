package transfersyntax

type TransferSyntax struct {
	UID         string
	Name        string
	Description string
	Type        string
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
	return nil
}
