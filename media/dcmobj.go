package media

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/tags"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/uid"
)

// DcmObj - DICOM Object structure
type DcmObj interface {
	Add(tag DcmTag)
	AddConceptNameSeq(group uint16, element uint16, CodeValue string, CodeMeaning string)
	AddSRText(text string)
	DumpTags()
	IsExplicitVR() bool
	SetExplicitVR(explicit bool)
	IsBigEndian() bool
	SetBigEndian(bigEndian bool)
	GetTag(i int) DcmTag
	SetTag(i int, tag DcmTag)
  	InsertTag(i int, tag DcmTag)
	DelTag(i int)
	GetTags() []DcmTag
	GetUShort(tag tags.Tag) uint16
	GetUInt(tag tags.Tag) uint32
	GetString(tag tags.Tag) string
	GetUShortGE(group uint16, element uint16) uint16
	GetUIntGE(group uint16, element uint16) uint32
	GetStringGE(group uint16, element uint16) string
	WriteDate(tag tags.Tag, date time.Time)
	WriteTime(tag tags.Tag, date time.Time)
	WriteUint16(tag tags.Tag, val uint16)
	WriteUint32(tag tags.Tag, val uint32)
	WriteString(tag tags.Tag, content string)
	WriteUint16GE(group uint16, element uint16, vr string, val uint16)
	WriteUint32GE(group uint16, element uint16, vr string, val uint32)
	WriteStringGE(group uint16, element uint16, vr string, content string)
	GetTransferSyntax() string
	SetTransferSyntax(ts string)
	TagCount() int
	CreateSR(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string)
	CreatePDF(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string, fileName string)
	WriteToBytes() []byte
	WriteToFile(fileName string) error
	dumpSeq(indent int)
}

type dcmObj struct {
	Tags           []DcmTag
	TransferSyntax string
	ExplicitVR     bool
	BigEndian      bool
	SQtag          DcmTag
}

// NewEmptyDCMObj - Create as an interface to a new empty dcmObj
func NewEmptyDCMObj() DcmObj {
	return &dcmObj{
		Tags:           make([]DcmTag, 0),
		TransferSyntax: "",
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          DcmTag{},
	}
}

// NewDCMObjFromFile - Read from a DICOM file into a DICOM Object
func NewDCMObjFromFile(fileName string) (DcmObj, error) {
	BigEndian := false

	if !fileExists(fileName) {
		return nil, errors.New("ERROR, DcmObj::Read, file does not exist")
	}

	bufdata, err := NewBufDataFromFile(fileName)
	if err != nil {
		return nil, err
	}

	obj := &dcmObj{
		Tags:           make([]DcmTag, 0),
		TransferSyntax: "",
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          DcmTag{},
	}

	obj.TransferSyntax = bufdata.ReadMeta()
	if len(obj.TransferSyntax) > 0 {
		if obj.TransferSyntax != uid.ImplicitVRLittleEndian {
			obj.ExplicitVR = true
		}
		if obj.TransferSyntax == uid.ExplicitVRBigEndian {
			BigEndian = true
		}
		bufdata.SetBigEndian(BigEndian)

		bufdata.ReadObj(obj)
	}

	return obj, nil
}

// NewDCMObjFromBytes - Read from a DICOM bytes into a DICOM Object
func NewDCMObjFromBytes(data []byte) DcmObj {
	BigEndian := false
	bufdata := NewBufDataFromBytes(data)

	obj := &dcmObj{
		Tags:           make([]DcmTag, 0),
		TransferSyntax: "",
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          DcmTag{},
	}

	obj.TransferSyntax = bufdata.ReadMeta()
	if len(obj.TransferSyntax) > 0 {
		if obj.TransferSyntax == uid.ImplicitVRLittleEndian {
			obj.ExplicitVR = false
		} else {
			obj.ExplicitVR = true
		}
		if obj.TransferSyntax == uid.ExplicitVRBigEndian {
			BigEndian = true
		}
		bufdata.SetBigEndian(BigEndian)
		bufdata.ReadObj(obj)
	}

	return obj
}

func (obj *dcmObj) IsExplicitVR() bool {
	return obj.ExplicitVR
}

func (obj *dcmObj) SetExplicitVR(explicit bool) {
	obj.ExplicitVR = explicit
}

func (obj *dcmObj) IsBigEndian() bool {
	return obj.BigEndian
}

func (obj *dcmObj) SetBigEndian(bigEndian bool) {
	obj.BigEndian = bigEndian
}

// TagCount - return the Tags number
func (obj *dcmObj) TagCount() int {
	return len(obj.Tags)
}

// GetTag - return the Tag at position i
func (obj *dcmObj) GetTag(i int) DcmTag {
	return obj.Tags[i]
}

func (obj *dcmObj) SetTag(i int, tag DcmTag) {
	if i<=obj.TagCount() {
		obj.Tags[i] = tag
	}
}

func (obj *dcmObj) InsertTag(index int, tag DcmTag) {
	obj.Tags = append(obj.Tags[:index+1], obj.Tags[index:]...)
	obj.Tags[index] = tag
}

func (obj *dcmObj) GetTags() []DcmTag {
	return obj.Tags
}

func (obj *dcmObj) DelTag(i int) {
	obj.Tags = append(obj.Tags[:i], obj.Tags[i+1:]...)
}

func (obj *dcmObj) DumpTags() {
	for _, tag := range obj.Tags {
		if tag.VR == "SQ" {
			fmt.Printf("\t(%04X,%04X) %s - %s\n", tag.Group, tag.Element, tag.VR, tag.Description)
			seq := tag.ReadSeq(obj.IsExplicitVR())
			seq.dumpSeq(1)
			continue
		}
		if tag.Length > 128 {
			fmt.Printf("\t(%04X,%04X) %s - %s : (Not displayed)\n", tag.Group, tag.Element, tag.VR, tag.Description)
			continue
		}
		fmt.Printf("\t(%04X,%04X) %s - %s : %s\n", tag.Group, tag.Element, tag.VR, tag.Description, tag.Data)
	}
	fmt.Println()
}

func (obj *dcmObj) dumpSeq(indent int) {
	tabs := "\t"
	for i := 0; i < indent; i++ {
		tabs += "\t"
	}

	for _, tag := range obj.Tags {
		if tag.VR == "SQ" {
			fmt.Printf("%s(%04X,%04X) %s - %s\n", tabs, tag.Group, tag.Element, tag.VR, tag.Description)
			seq := tag.ReadSeq(obj.IsExplicitVR())
			seq.dumpSeq(indent + 1)
			continue
		}
		if tag.Length > 128 {
			fmt.Printf("%s(%04X,%04X) %s - %s : (Not displayed)\n", tabs, tag.Group, tag.Element, tag.VR, tag.Description)
			continue
		}
		fmt.Printf("%s(%04X,%04X) %s - %s : %s\n", tabs, tag.Group, tag.Element, tag.VR, tag.Description, tag.Data)
	}
}

func (obj *dcmObj) GetUShort(tag tags.Tag) uint16 {
	return obj.GetUShortGE(tag.Group, tag.Element)
}

// GetUShortGE - return the Uint16 for this group & element
func (obj *dcmObj) GetUShortGE(group uint16, element uint16) uint16 {
	var i int
	var tag DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if (sq == 0) && (tag.Length > 0) && (tag.Length != 0xFFFFFFFF) {
			if (tag.Group == group) && (tag.Element == element) {
				break
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	if i < obj.TagCount() {
		return tag.GetUShort()
	}
	return 0
}

func (obj *dcmObj) GetUInt(tag tags.Tag) uint32 {
	return obj.GetUIntGE(tag.Group, tag.Element)
}

// GetUIntGE - return the Uint32 for this group & element
func (obj *dcmObj) GetUIntGE(group uint16, element uint16) uint32 {
	var i int
	var tag DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if (sq == 0) && (tag.Length > 0) && (tag.Length != 0xFFFFFFFF) {
			if (tag.Group == group) && (tag.Element == element) {
				break
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	if i < obj.TagCount() {
		return tag.GetUInt()
	}
	return 0
}

func (obj *dcmObj) GetString(tag tags.Tag) string {
	return obj.GetStringGE(tag.Group, tag.Element)
}

// GetStringGE - return the String for this group & element
func (obj *dcmObj) GetStringGE(group uint16, element uint16) string {
	var i int
	var tag DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if (sq == 0) && (tag.Length > 0) && (tag.Length != 0xFFFFFFFF) {
			if (tag.Group == group) && (tag.Element == element) {
				break
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	if i < obj.TagCount() {
		return tag.GetString()
	}
	return ""
}

// Add - add a new DICOM Tag to a DICOM Object
func (obj *dcmObj) Add(tag DcmTag) {
	obj.Tags = append(obj.Tags, tag)
}

func (obj *dcmObj) WriteToBytes() []byte {
	bufdata := NewEmptyBufData()

	if obj.TransferSyntax == uid.ExplicitVRBigEndian {
		bufdata.SetBigEndian(true)
	}
	SOPClassUID := obj.GetStringGE(0x08, 0x16)
	SOPInstanceUID := obj.GetStringGE(0x08, 0x18)
	bufdata.WriteMeta(SOPClassUID, SOPInstanceUID, obj.TransferSyntax)
	bufdata.WriteObj(obj)
	bufdata.SetPosition(0)
	return bufdata.GetAllBytes()
}

// Wrote - Write a DICOM Object to a DICOM File
func (obj *dcmObj) WriteToFile(fileName string) error {
	bufdata := NewEmptyBufData()

	if obj.TransferSyntax == uid.ExplicitVRBigEndian {
		bufdata.SetBigEndian(true)
	}
	SOPClassUID := obj.GetStringGE(0x08, 0x16)
	SOPInstanceUID := obj.GetStringGE(0x08, 0x18)
	bufdata.WriteMeta(SOPClassUID, SOPInstanceUID, obj.TransferSyntax)
	bufdata.WriteObj(obj)
	bufdata.SetPosition(0)
	return bufdata.SaveToFile(fileName)
}

func (obj *dcmObj) WriteDate(tag tags.Tag, date time.Time) {
	obj.WriteString(tag, date.Format("20060102"))
}

func (obj *dcmObj) WriteTime(tag tags.Tag, date time.Time) {
	obj.WriteString(tag, date.Format("150405"))
}

func (obj *dcmObj) WriteUint16(tag tags.Tag, val uint16) {
	obj.WriteUint16GE(tag.Group, tag.Element, tag.VR, val)
}

func (obj *dcmObj) WriteUint32(tag tags.Tag, val uint32) {
	obj.WriteUint32GE(tag.Group, tag.Element, tag.VR, val)
}

func (obj *dcmObj) WriteString(tag tags.Tag, content string) {
	obj.WriteStringGE(tag.Group, tag.Element, tag.VR, content)
}

// WriteUint16GE - Writes a Uint16 to a DICOM tag
func (obj *dcmObj) WriteUint16GE(group uint16, element uint16, vr string, val uint16) {
	c := make([]byte, 2)
	if obj.BigEndian {
		binary.BigEndian.PutUint16(c, val)
	} else {
		binary.LittleEndian.PutUint16(c, val)
	}

	tag := DcmTag{
		Group:     group,
		Element:   element,
		Length:    2,
		VR:        vr,
		Data:      c,
		BigEndian: obj.BigEndian,
	}
	FillTag(&tag)
	obj.Tags = append(obj.Tags, tag)
}

// WriteUint32GE - Writes a Uint32 to a DICOM tag
func (obj *dcmObj) WriteUint32GE(group uint16, element uint16, vr string, val uint32) {
	c := make([]byte, 4)
	if obj.BigEndian {
		binary.BigEndian.PutUint32(c, val)
	} else {
		binary.LittleEndian.PutUint32(c, val)
	}

	tag := DcmTag{
		Group:     group,
		Element:   element,
		Length:    4,
		VR:        vr,
		Data:      c,
		BigEndian: obj.BigEndian,
	}
	FillTag(&tag)
	obj.Tags = append(obj.Tags, tag)
}

// WriteStringGE - Writes a String to a DICOM tag
func (obj *dcmObj) WriteStringGE(group uint16, element uint16, vr string, content string) {
	var length uint32

	length = uint32(len(content))
	if length%2 == 1 {
		length++
		if vr == "UI" {
			content = content + string(0)
		} else {
			content = content + " "
		}
	}
	tag := DcmTag{
		Group:     group,
		Element:   element,
		Length:    length,
		VR:        vr,
		Data:      []byte(content),
		BigEndian: false,
	}
	FillTag(&tag)
	obj.Tags = append(obj.Tags, tag)
}

func (obj *dcmObj) GetTransferSyntax() string {
	return obj.TransferSyntax
}

func (obj *dcmObj) SetTransferSyntax(ts string) {
	obj.TransferSyntax = ts
}

// AddConceptNameSeq - Concept Name Sequence for DICOM SR
func (obj *dcmObj) AddConceptNameSeq(group uint16, element uint16, CodeValue string, CodeMeaning string) {
	item := &dcmObj{
		Tags:           make([]DcmTag, 0),
		TransferSyntax: "",
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          DcmTag{},
	}
	seq := &dcmObj{
		Tags:           make([]DcmTag, 0),
		TransferSyntax: "",
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          DcmTag{},
	}
	tag := DcmTag{}

	item.BigEndian = obj.BigEndian
	item.ExplicitVR = obj.ExplicitVR
	seq.BigEndian = obj.BigEndian
	seq.ExplicitVR = obj.ExplicitVR

	item.WriteString(tags.CodeValue, CodeValue)
	item.WriteString(tags.CodingSchemeDesignator, "OneByteData")
	item.WriteString(tags.CodeMeaning, CodeMeaning)
	tag.WriteSeq(0xFFFE, 0xE000, item)
	seq.Add(tag)
	tag.WriteSeq(group, element, seq)
	obj.Add(tag)
}

// AddSRText - add Text to SR
func (obj *dcmObj) AddSRText(text string) {
	item := &dcmObj{
		Tags:           make([]DcmTag, 0),
		TransferSyntax: "",
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          DcmTag{},
	}
	seq := &dcmObj{
		Tags:           make([]DcmTag, 0),
		TransferSyntax: "",
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          DcmTag{},
	}
	tag := DcmTag{}

	item.BigEndian = obj.BigEndian
	item.ExplicitVR = obj.ExplicitVR
	seq.BigEndian = obj.BigEndian
	seq.ExplicitVR = obj.ExplicitVR

	item.WriteString(tags.RelationshipType, "CONTAINS")
	item.WriteString(tags.ValueType, "TEXT")
	item.AddConceptNameSeq(0x40, 0xA043, "2222", "Report Text")
	item.WriteString(tags.TextValue, text)
	tag.WriteSeq(0xFFFE, 0xE000, item)
	seq.Add(tag)
	tag.WriteSeq(0x40, 0xA730, seq)
	obj.Add(tag)
}

// CreateSR - Create a DICOM SR object
func (obj *dcmObj) CreateSR(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string) {
	obj.WriteString(tags.InstanceCreationDate, time.Now().Format("20060102"))
	obj.WriteString(tags.InstanceCreationTime, time.Now().Format("150405"))
	obj.WriteString(tags.SOPClassUID, "1.2.840.10008.5.1.4.1.1.88.11")
	obj.WriteString(tags.SOPInstanceUID, SOPInstanceUID)
	obj.WriteString(tags.AccessionNumber, study.AccessionNumber)
	obj.WriteString(tags.Modality, "SR")
	obj.WriteString(tags.InstitutionName, study.InstitutionName)
	obj.WriteString(tags.ReferringPhysicianName, study.ReferringPhysician)
	obj.WriteString(tags.StudyDescription, study.Description)
	obj.WriteString(tags.SeriesDescription, "REPORT")
	obj.WriteString(tags.PatientName, study.PatientName)
	obj.WriteString(tags.PatientID, study.PatientID)
	obj.WriteString(tags.PatientBirthDate, study.PatientBD)
	obj.WriteString(tags.PatientSex, study.PatientSex)
	obj.WriteString(tags.StudyInstanceUID, study.StudyInstanceUID)
	obj.WriteString(tags.SeriesInstanceUID, SeriesInstanceUID)
	obj.WriteString(tags.SeriesNumber, "200")
	obj.WriteString(tags.InstanceNumber, "1")
	obj.WriteString(tags.ValueType, "CONTAINER")
	obj.AddConceptNameSeq(0x0040, 0xA043, "1111", "Radiology Report")
	obj.WriteString(tags.ContinuityOfContent, "SEPARATE")
	obj.WriteString(tags.VerifyingObserverName, study.ObserverName)
	obj.WriteString(tags.CompletionFlag, "COMPLETE")
	obj.WriteString(tags.VerificationFlag, "VERIFIED")
	obj.AddSRText(study.ReportText)
}

// CreatePDF - Create a DICOM SR object
func (obj *dcmObj) CreatePDF(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string, fileName string) {
	obj.WriteString(tags.InstanceCreationDate, time.Now().Format("20060102"))
	obj.WriteString(tags.InstanceCreationTime, time.Now().Format("150405"))
	obj.WriteString(tags.SOPClassUID, "1.2.840.10008.5.1.4.1.1.104.1")
	obj.WriteString(tags.SOPInstanceUID, SOPInstanceUID)
	obj.WriteString(tags.AccessionNumber, study.AccessionNumber)
	obj.WriteString(tags.Modality, "OT")
	obj.WriteString(tags.InstitutionName, study.InstitutionName)
	obj.WriteString(tags.ReferringPhysicianName, study.ReferringPhysician)
	obj.WriteString(tags.StudyDescription, study.Description)
	obj.WriteString(tags.PatientName, study.PatientName)
	obj.WriteString(tags.PatientID, study.PatientID)
	obj.WriteString(tags.PatientBirthDate, study.PatientBD)
	obj.WriteString(tags.PatientSex, study.PatientSex)
	obj.WriteString(tags.StudyInstanceUID, study.StudyInstanceUID)
	obj.WriteString(tags.SeriesInstanceUID, SeriesInstanceUID)
	obj.WriteString(tags.SeriesNumber, "300")
	obj.WriteString(tags.InstanceNumber, "1")

	mstream, _ := NewMemoryStreamFromFile(fileName)

	mstream.SetPosition(0)
	size := uint32(mstream.GetSize())
	if size%2 == 1 {
		size++
		mstream.Append([]byte{0})
	}
	obj.WriteString(tags.DocumentTitle, fileName)
	obj.Add(DcmTag{
		Group:     0x42,
		Element:   0x11,
		Length:    size,
		VR:        "OB",
		Data:      mstream.GetData(),
		BigEndian: obj.BigEndian,
	})
	obj.WriteString(tags.MIMETypeOfEncapsulatedDocument, "application/pdf")
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			log.Println("ERROR, dcmobj::fileExists, " + err.Error())
			return false
		}
	}
	return true
}
