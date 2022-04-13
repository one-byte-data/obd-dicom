package media

import (
	"encoding/xml"
	"io/ioutil"
	"strconv"

	"git.onebytedata.com/odb/go-dicom/tags"
)

type dictionary struct {
	XMLName xml.Name `xml:"dictionary"`
	Tags    []xmlTag `xml:"tag"`
}

type xmlTag struct {
	Group       string `xml:"group,attr"`
	Element     string `xml:"element,attr"`
	Name        string `xml:"keyword,attr"`
	VR          string `xml:"vr,attr"`
	VM          string `xml:"vm,attr"`
	Description string `xml:",chardata"`
}

var codes []*tags.Tag

// FillTag - Populates with data from dictionary
func FillTag(tag *DcmTag) {
	dt := GetDictionaryTag(tag.Group, tag.Element)
	tag.Name = dt.Name
	tag.Description = dt.Description
	tag.VR = dt.VR
	tag.VM = dt.VM
}

// GetDictionaryTag - get tag from Dictionary
func GetDictionaryTag(group uint16, element uint16) *tags.Tag {
	if codes == nil {
		InitDict()
	}
	for i := 0; i < len(codes); i++ {
		if (group == codes[i].Group) && (element == codes[i].Element) {
			return codes[i]
		}
	}
	return &tags.Tag{
		Group:       0,
		Element:     0,
		VR:          "UN",
		VM:          "",
		Name:        "Unknown",
		Description: "Unknown",
	}
}

// GetDictionaryVR - get info from Dictionary
func GetDictionaryVR(group uint16, element uint16) string {
	if codes == nil {
		InitDict()
	}
	for i := 0; i < len(codes); i++ {
		if (group == codes[i].Group) && (element == codes[i].Element) {
			return codes[i].VR
		}
	}
	return "UN"
}

func loadPrivateDictionary() {
	privateDictionaryFile := "./private.xml"
	data, err := ioutil.ReadFile(privateDictionaryFile)
	if err != nil {
		return
	}

	dict := new(dictionary)
	err = xml.Unmarshal(data, dict)
	if err != nil {
		return
	}

	for _, t := range dict.Tags {
		g, err := strconv.Atoi(t.Group)
		if err != nil {
			continue
		}
		e, err := strconv.Atoi(t.Element)
		if err != nil {
			continue
		}

		codes = append(codes, &tags.Tag{
			Group:       uint16(g),
			Element:     uint16(e),
			Name:        t.Name,
			Description: t.Description,
			VR:          t.VR,
			VM:          t.VM,
		})
	}
}

// InitDict Initialize Dictionary
func InitDict() {
	codes = tags.GetTags()
	loadPrivateDictionary()
}
