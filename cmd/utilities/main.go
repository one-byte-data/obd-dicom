package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	httpclient "git.onebytedata.com/odb/go-dicom/clients/http_client"
)

const dictionaryURL string = "https://raw.githubusercontent.com/fo-dicom/fo-dicom/development/FO-DICOM.Core/Dictionaries/DICOM%20Dictionary.xml"

const dictionaryTagsFile string = "../../tags/dicom-tags.go"

const tagsFileName string = "../../tags/tags.go"

type dictionary struct {
	XMLName xml.Name `xml:"dictionary"`
	Tags    []tag    `xml:"tag"`
}

type tag struct {
	Group   string `xml:"group,attr"`
	Element string `xml:"element,attr"`
	Keyword string `xml:"keyword,attr"`
	VR      string `xml:"vr,attr"`
	VM      string `xml:"vm,attr"`
	Name    string `xml:",chardata"`
}

func main() {
	tags := downloadDictionary()
	writeTagsFile(tags)
	writeDictionaryTags(tags)
}

func downloadDictionary() []tag {
	params := httpclient.HTTPParams{
		URL: dictionaryURL,
	}
	client := httpclient.NewHTTPClient(params)
	response, err := client.Get()
	if err != nil {
		log.Panic(err)
	}

	dict := new(dictionary)
	err = xml.Unmarshal(response, dict)
	if err != nil {
		log.Panic(err)
	}
	return dict.Tags
}

func writeDictionaryTags(tags []tag) {
	if FileExists(dictionaryTagsFile) {
		err := os.Remove(dictionaryTagsFile)
		if err != nil {
			log.Panic(err)
		}
	}
	f, err := os.Create(dictionaryTagsFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	f.WriteString("package tags\n\n")

	for _, tag := range tags {
		if strings.Contains(tag.Group, "x") || strings.Contains(tag.Element, "x") {
			continue
		}
		f.WriteString(fmt.Sprintf("// %s - (%s,%s) %s\n", tag.Keyword, tag.Group, tag.Element, tag.Name))
		f.WriteString(fmt.Sprintf("var %s = &Tag{\n", tag.Keyword))
		f.WriteString(fmt.Sprintf("  Group: 0x%s,\n", tag.Group))
		f.WriteString(fmt.Sprintf("  Element: 0x%s,\n", tag.Element))
		f.WriteString(fmt.Sprintf("  VR: \"%s\",\n", tag.VR))
		f.WriteString(fmt.Sprintf("  VM: \"%s\",\n", tag.VM))
		f.WriteString(fmt.Sprintf("  Name: \"%s\",\n", tag.Keyword))
		f.WriteString(fmt.Sprintf("  Description: \"%s\",\n", tag.Name))
		f.WriteString("}\n")
	}

	f.Sync()
}

func writeTagsFile(tags []tag) {
	if FileExists(tagsFileName) {
		err := os.Remove(tagsFileName)
		if err != nil {
			log.Panic(err)
		}
	}
	f, err := os.Create(tagsFileName)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	f.WriteString("package tags\n\n")
	f.WriteString("var tags = []*Tag{\n")

	for _, tag := range tags {
		if strings.Contains(tag.Group, "x") || strings.Contains(tag.Element, "x") {
			continue
		}
		f.WriteString("  {\n")
		f.WriteString(fmt.Sprintf("    Name: \"%s\",\n", tag.Keyword))
		f.WriteString(fmt.Sprintf("    Description: \"%s\",\n", tag.Name))
		f.WriteString(fmt.Sprintf("    Group: 0x%s,\n", tag.Group))
		f.WriteString(fmt.Sprintf("    Element: 0x%s,\n", tag.Element))
		f.WriteString(fmt.Sprintf("    VR: \"%s\",\n", tag.VR))
		f.WriteString(fmt.Sprintf("    VM: \"%s\",\n", tag.VM))
		f.WriteString("  },\n")
	}

	f.WriteString("}\n")
	f.Sync()
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}