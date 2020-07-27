package media

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// DcmObj - DICOM Object structure
type DcmObj interface {
	Add(tag DcmTag)
	AddConceptNameSeq(group uint16, element uint16, CodeValue string, CodeMeaning string)
	AddSRText(text string)
	DumpTags()
	Clear()
	IsExplicitVR() bool
	SetExplicitVR(explicit bool)
	IsBigEndian() bool
	SetBigEndian(bigEndian bool)
	GetTag(i int) DcmTag
	SetTag(i int, tag DcmTag)
	GetTags() []DcmTag
	GetUShort(group uint16, element uint16) uint16
	GetUInt(group uint16, element uint16) uint32
	GetString(group uint16, element uint16) string
	WriteUint16(group uint16, element uint16, vr string, val uint16)
	WriteUint32(group uint16, element uint16, vr string, val uint32)
	WriteString(group uint16, element uint16, vr string, content string)
	GetTransferSynxtax() string
	SetTransferSyntax(ts string)
	TagCount() int
	CreateSR(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string)
	CreatePDF(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string, fileName string)
	WriteToBytes() []byte
	WriteToFile(fileName string) error
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
		if obj.TransferSyntax != "1.2.840.10008.1.2" {
			obj.ExplicitVR = true
		}
		if obj.TransferSyntax == "1.2.840.10008.1.2.2" {
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
		if obj.TransferSyntax == "1.2.840.10008.1.2" {
			obj.ExplicitVR = false
		} else {
			obj.ExplicitVR = true
		}
		if obj.TransferSyntax == "1.2.840.10008.1.2.2" {
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

// Clear - clear tags array
func (obj *dcmObj) Clear() {
	obj.Tags = nil
}

// GetTag - return the Tag at position i
func (obj *dcmObj) GetTag(i int) DcmTag {
	return obj.Tags[i]
}

func (obj *dcmObj) SetTag(i int, tag DcmTag) {
	obj.Tags[i] = tag
}

func (obj *dcmObj) GetTags() []DcmTag {
	return obj.Tags
}

func (obj *dcmObj) DumpTags() {
	for _, tag := range obj.Tags {
		if tag.VR == "SQ" {
			fmt.Printf("\t(%04X,%04X) %s - %s\n", tag.Group, tag.Element, tag.VR, TagDescription(tag.Group, tag.Element))
			continue
		}
		if tag.Length > 128 {
			fmt.Printf("\t(%04X,%04X) %s - %s : (Not displayed)\n", tag.Group, tag.Element, tag.VR, TagDescription(tag.Group, tag.Element))
			continue
		}
		fmt.Printf("\t(%04X,%04X) %s - %s : %s\n", tag.Group, tag.Element, tag.VR, TagDescription(tag.Group, tag.Element), tag.Data)
	}
}

// GetUShort - return the Uint16 for this group & element
func (obj *dcmObj) GetUShort(group uint16, element uint16) uint16 {
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

// GetUInt - return the Uint32 for this group & element
func (obj *dcmObj) GetUInt(group uint16, element uint16) uint32 {
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

// GetString - return the String for this group & element
func (obj *dcmObj) GetString(group uint16, element uint16) string {
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

	if obj.TransferSyntax == "1.2.840.10008.1.2.2" {
		bufdata.SetBigEndian(true)
	}
	SOPClassUID := obj.GetString(0x08, 0x16)
	SOPInstanceUID := obj.GetString(0x08, 0x18)
	bufdata.WriteMeta(SOPClassUID, SOPInstanceUID, obj.TransferSyntax)
	bufdata.WriteObj(obj)
	bufdata.SetPosition(0)
	return bufdata.GetAllBytes()
}

// Wrote - Write a DICOM Object to a DICOM File
func (obj *dcmObj) WriteToFile(fileName string) error {
	bufdata := NewEmptyBufData()

	if obj.TransferSyntax == "1.2.840.10008.1.2.2" {
		bufdata.SetBigEndian(true)
	}
	SOPClassUID := obj.GetString(0x08, 0x16)
	SOPInstanceUID := obj.GetString(0x08, 0x18)
	bufdata.WriteMeta(SOPClassUID, SOPInstanceUID, obj.TransferSyntax)
	bufdata.WriteObj(obj)
	bufdata.SetPosition(0)
	return bufdata.SaveToFile(fileName)
}

// WriteUint16 - Writes a Uint16 to a DICOM tag
func (obj *dcmObj) WriteUint16(group uint16, element uint16, vr string, val uint16) {
	c := make([]byte, 2)
	if obj.BigEndian {
		binary.BigEndian.PutUint16(c, val)
	} else {
		binary.LittleEndian.PutUint16(c, val)
	}

	tag := DcmTag{group, element, 2, vr, c, obj.BigEndian}
	obj.Tags = append(obj.Tags, tag)
}

// WriteUint32 - Writes a Uint32 to a DICOM tag
func (obj *dcmObj) WriteUint32(group uint16, element uint16, vr string, val uint32) {
	c := make([]byte, 4)
	if obj.BigEndian {
		binary.BigEndian.PutUint32(c, val)
	} else {
		binary.LittleEndian.PutUint32(c, val)
	}

	tag := DcmTag{group, element, 4, vr, c, obj.BigEndian}
	obj.Tags = append(obj.Tags, tag)
}

// WriteString - Writes a String to a DICOM tag
func (obj *dcmObj) WriteString(group uint16, element uint16, vr string, content string) {
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
	tag := DcmTag{group, element, length, vr, []byte(content), false}
	obj.Tags = append(obj.Tags, tag)
}

func (obj *dcmObj) GetTransferSynxtax() string {
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

	// Code Value
	item.WriteString(0x08, 0x100, "SH", CodeValue)
	// Coding Scheme Designator
	item.WriteString(0x08, 0x102, "SH", "OneByteData")
	// Code Meaning
	item.WriteString(0x08, 0x104, "LO", CodeMeaning)
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

	// Relationship Type
	item.WriteString(0x40, 0xA010, "CS", "CONTAINS")
	// Value Type
	item.WriteString(0x40, 0xA040, "CS", "TEXT")
	item.AddConceptNameSeq(0x40, 0xA043, "2222", "Report Text")
	item.WriteString(0x40, 0xA160, "UT", text)
	tag.WriteSeq(0xFFFE, 0xE000, item)
	seq.Add(tag)
	tag.WriteSeq(0x40, 0xA730, seq)
	obj.Add(tag)
}

// CreateSR - Create a DICOM SR object
func (obj *dcmObj) CreateSR(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string) {
	// Instance Creation Date
	obj.WriteString(0x08, 0x12, "DA", time.Now().Format("20060102"))
	// Instance Creation Time
	obj.WriteString(0x08, 0x13, "TM", time.Now().Format("150405"))
	// Basic Text SOP Class
	obj.WriteString(0x08, 0x16, "UI", "1.2.840.10008.5.1.4.1.1.88.11")
	obj.WriteString(0x0008, 0x0018, "UI", SOPInstanceUID)
	obj.WriteString(0x0008, 0x0050, "SH", study.AccessionNumber) // Accession Number
	obj.WriteString(0x0008, 0x0060, "CS", "SR")
	obj.WriteString(0x0008, 0x0080, "LO", study.InstitutionName)
	obj.WriteString(0x0008, 0x0090, "PN", study.ReferringPhysician)
	obj.WriteString(0x0008, 0x1030, "LO", study.Description)
	obj.WriteString(0x0008, 0x103E, "LO", "REPORT")
	obj.WriteString(0x0010, 0x0010, "PN", study.PatientName)
	obj.WriteString(0x0010, 0x0020, "LO", study.PatientID)  // Patient ID
	obj.WriteString(0x0010, 0x0030, "DA", study.PatientBD)  // Patient's Birth Date
	obj.WriteString(0x0010, 0x0040, "CS", study.PatientSex) // Patient's Sex
	obj.WriteString(0x0020, 0x000d, "UI", study.StudyInstanceUID)
	obj.WriteString(0x0020, 0x000e, "UI", SeriesInstanceUID)
	obj.WriteString(0x0020, 0x0011, "IS", "200")
	obj.WriteString(0x0020, 0x0013, "IS", "1")
	obj.WriteString(0x0040, 0xA040, "CS", "CONTAINER") // Value Type
	obj.AddConceptNameSeq(0x0040, 0xA043, "1111", "Radiology Report")
	obj.WriteString(0x0040, 0xA050, "CS", "SEPARATE") // Continuity of Context
	obj.WriteString(0x40, 0xa075, "PN", study.ObserverName)
	obj.WriteString(0x0040, 0xa491, "CS", "COMPLETE") // CompletionFlag
	obj.WriteString(0x0040, 0xa493, "CS", "VERIFIED") // VerificationFlag
	obj.AddSRText(study.ReportText)
}

// CreatePDF - Create a DICOM SR object
func (obj *dcmObj) CreatePDF(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string, fileName string) {
	// Instance Creation Date
	obj.WriteString(0x08, 0x12, "DA", time.Now().Format("20060102"))
	// Instance Creation Time
	obj.WriteString(0x08, 0x13, "TM", time.Now().Format("150405"))
	// DICOM PDF SOP Class
	obj.WriteString(0x08, 0x16, "UI", "1.2.840.10008.5.1.4.1.1.104.1")
	obj.WriteString(0x0008, 0x0018, "UI", SOPInstanceUID)
	obj.WriteString(0x0008, 0x0050, "SH", study.AccessionNumber) // Accession Number
	obj.WriteString(0x0008, 0x0060, "CS", "OT")
	obj.WriteString(0x0008, 0x0080, "LO", study.InstitutionName)
	obj.WriteString(0x0008, 0x0090, "PN", study.ReferringPhysician)
	obj.WriteString(0x0008, 0x1030, "LO", study.Description)
	obj.WriteString(0x0010, 0x0010, "PN", study.PatientName)
	obj.WriteString(0x0010, 0x0020, "LO", study.PatientID)  // Patient ID
	obj.WriteString(0x0010, 0x0030, "DA", study.PatientBD)  // Patient's Birth Date
	obj.WriteString(0x0010, 0x0040, "CS", study.PatientSex) // Patient's Sex
	obj.WriteString(0x0020, 0x000d, "UI", study.StudyInstanceUID)
	obj.WriteString(0x0020, 0x000e, "UI", SeriesInstanceUID)
	obj.WriteString(0x0020, 0x0011, "IS", "300")
	obj.WriteString(0x0020, 0x0013, "IS", "1")

	mstream, _ := NewMemoryStreamFromFile(fileName)

	mstream.SetPosition(0)
	size := uint32(mstream.GetSize())
	if size%2 == 1 {
		size++
		mstream.Append([]byte{0})
	}
	obj.WriteString(0x42, 0x10, "ST", fileName)
	obj.Add(DcmTag{0x42, 0x11, size, "OB", mstream.GetData(), obj.BigEndian})
	obj.WriteString(0x42, 0x12, "LO", "application/pdf")
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
