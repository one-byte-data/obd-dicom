package media

import (
	"os"
	"time"
)

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// DcmObj - DICOM Object structure
type DcmObj struct {
	tags           []DcmTag
	TransferSyntax string
	ExplicitVR     bool
	BigEndian      bool
	SQtag          DcmTag
}

// TagCount - return the Tags number
func (obj *DcmObj) TagCount() int {
	return len(obj.tags)
}

// GetTag - return the Tag at position i
func (obj *DcmObj) GetTag(i int) DcmTag {
	return obj.tags[i]
}

// GetString - return the String for this group & element
func (obj *DcmObj) GetString(group uint16, element uint16) string {
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
func (obj *DcmObj) Add(tag DcmTag) {
	obj.tags = append(obj.tags, tag)
}

// Read - Read from a DICOM file into a DICOM Object
func (obj *DcmObj) Read(FileName string) bool {
	flag := false
	BigEndian := false
	var bufdata BufData

	if fileExists(FileName) {
		if bufdata.Ms.LoadFromFile(FileName) {
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
				bufdata.BigEndian = BigEndian
				flag = bufdata.ReadObj(obj)
			}
		}
	}
	return flag
}

// Wrote - Write a DICOM Object to a DICOM File
func (obj *DcmObj) Write(FileName string) bool {
	var bufdata BufData
	bufdata.BigEndian = false
	if obj.TransferSyntax == "1.2.840.10008.1.2.2" {
		bufdata.BigEndian = true
	}
	SOPClassUID := obj.GetString(0x08, 0x16)
	SOPInstanceUID := obj.GetString(0x08, 0x18)
	bufdata.WriteMeta(SOPClassUID, SOPInstanceUID, obj.TransferSyntax)
	bufdata.WriteObj(obj)
	bufdata.Ms.Position = 0
	return (bufdata.Ms.SaveToFile(FileName))
}

// WriteString - Writes a String to a DICOM tag
func (obj *DcmObj) WriteString(group uint16, element uint16, vr string, content string) {
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
	obj.tags = append(obj.tags, tag)
}

// AddConceptNameSeq - Concept Name Sequence for DICOM SR
func (obj *DcmObj) AddConceptNameSeq(group uint16, element uint16, CodeValue string, CodeMeaning string) {
	var item DcmObj
	var seq DcmObj
	var tag DcmTag

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
func (obj *DcmObj) AddSRText(text string) {
	var item DcmObj
	var seq DcmObj
	var tag DcmTag

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
func (obj *DcmObj) CreateSR(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string) {
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
func (obj *DcmObj) CreatePDF(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string, FileName string) {
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

	var mstream MemoryStream
	mstream.LoadFromFile(FileName)
	mstream.Position = 0
	size := uint32(mstream.Size)
	if size%2 == 1 {
		size++
		mstream.data=append(mstream.data, 0)
	}
	obj.WriteString(0x42, 0x10, "ST", FileName)
	obj.Add(DcmTag{0x42, 0x11, size, "OB", mstream.data, obj.BigEndian})
	obj.WriteString(0x42, 0x12, "LO", "application/pdf")
}
