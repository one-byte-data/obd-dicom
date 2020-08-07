package media

// DictStruct Dictionary Structure definition
type DictStruct struct {
	Group       uint16
	Element     uint16
	VR          string
	VM          string
	Name        string
	Description string
}

var codes []DictStruct

// FillTag - Populates with data from dictionary
func FillTag(tag *DcmTag) {
	dt := GetDictionaryTag(tag.Group, tag.Element)
	tag.Name = dt.Name
	tag.Description = dt.Description
	tag.VR = dt.VR
	tag.VM = dt.VM
}
// GetDictionaryTag - get tag from Dictionary
func GetDictionaryTag(group uint16, element uint16) *DictStruct {
	if codes == nil {
		return nil
	}
	for i := 0; i < len(codes); i++ {
		if (group == codes[i].Group) && (element == codes[i].Element) {
			return &codes[i]
		}
	}
	return &DictStruct{
		Group: 0,
		Element: 0,
		VR: "UN",
		VM: "",
		Name: "Unknown",
		Description: "Unknown",
	}
}

// GetDictionaryVR - get info from Dictionary
func GetDictionaryVR(group uint16, element uint16) string {
	if codes == nil {
		return "ERROR"
	}
	for i := 0; i < len(codes); i++ {
		if (group == codes[i].Group) && (element == codes[i].Element) {
			return codes[i].VR
		}
	}
	return "UN"
}

// InitDict Initialize Dictionary
func InitDict() {
	codes = Tags
}
