package uid

type SOPClass struct {
	UID         string
	Name        string
	Description string
	Type        string
}

func GetSOPClassFromName(name string) *SOPClass {
	for _, sop := range SOPClasses {
		if sop.Name == name {
			return sop
		}
	}
	for _, sop := range TransferSyntaxes {
		if sop.Name == name {
			return sop
		}
	}
	return nil
}

func GetSOPClassFromUID(uid string) *SOPClass {
	for _, sop := range SOPClasses {
		if sop.UID == uid {
			return sop
		}
	}
	for _, sop := range TransferSyntaxes {
		if sop.UID == uid {
			return sop
		}
	}
	return nil
}

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
