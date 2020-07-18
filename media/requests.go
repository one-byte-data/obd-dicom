package media

// DefaultCFindRequest - Creates a default C-Find request
func DefaultCFindRequest() DcmObj {
	query := NewEmptyDCMObj()
	query.WriteString(0x08, 0x20, "DA", "")
	query.WriteString(0x08, 0x30, "TM", "")
	query.WriteString(0x08, 0x50, "SH", "")
	query.WriteString(0x08, 0x52, "CS", "STUDY")
	query.WriteString(0x08, 0x61, "CS", "")
	query.WriteString(0x08, 0x1030, "LO", "")
	query.WriteString(0x10, 0x10, "PN", "")
	query.WriteString(0x10, 0x20, "LO", "")
	query.WriteString(0x10, 0x30, "DA", "")
	query.WriteString(0x10, 0x40, "CS", "")
	query.WriteString(0x20, 0x0D, "UI", "")
	return query
}

// DefaultCMoveRequest - Creates a default C-Move request
func DefaultCMoveRequest(studyUID string) DcmObj {
	query := NewEmptyDCMObj()
	query.WriteString(0x08, 0x20, "DA", "")
	query.WriteString(0x08, 0x30, "TM", "")
	query.WriteString(0x08, 0x50, "SH", "")
	query.WriteString(0x08, 0x52, "CS", "STUDY")
	query.WriteString(0x08, 0x61, "CS", "")
	query.WriteString(0x08, 0x1030, "LO", "")
	query.WriteString(0x10, 0x10, "PN", "")
	query.WriteString(0x10, 0x20, "LO", "")
	query.WriteString(0x10, 0x30, "DA", "")
	query.WriteString(0x10, 0x40, "CS", "")
	query.WriteString(0x20, 0x0D, "UI", studyUID)
	return query
}
