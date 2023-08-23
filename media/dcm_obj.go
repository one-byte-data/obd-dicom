package media

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/one-byte-data/obd-dicom/dictionary/sopclass"
	"github.com/one-byte-data/obd-dicom/dictionary/tags"
	"github.com/one-byte-data/obd-dicom/dictionary/transfersyntax"
	"github.com/one-byte-data/obd-dicom/jpeglib"
	"github.com/one-byte-data/obd-dicom/openjpeg"
	"github.com/one-byte-data/obd-dicom/transcoder"
)

// DcmObj - DICOM Object structure
type DcmObj interface {
	Add(tag *DcmTag)
	AddConceptNameSeq(group uint16, element uint16, CodeValue string, CodeMeaning string)
	AddSRText(text string)
	DumpTags()
	IsExplicitVR() bool
	SetExplicitVR(explicit bool)
	IsBigEndian() bool
	SetBigEndian(bigEndian bool)
	GetDate(tag *tags.Tag) time.Time
	GetPixelData(frame int) ([]byte, error)
	GetTagAt(i int) *DcmTag
	GetTag(tag *tags.Tag) *DcmTag
	GetTagGE(group uint16, element uint16) *DcmTag
	SetTag(i int, tag *DcmTag)
	InsertTag(i int, tag *DcmTag)
	DelTag(i int)
	GetTags() []*DcmTag
	GetUShort(tag *tags.Tag) uint16
	GetUInt(tag *tags.Tag) uint32
	GetString(tag *tags.Tag) string
	GetUShortGE(group uint16, element uint16) uint16
	GetUIntGE(group uint16, element uint16) uint32
	GetStringGE(group uint16, element uint16) string
	WriteDate(tag *tags.Tag, date time.Time)
	WriteDateRange(tag *tags.Tag, startDate time.Time, endDate time.Time)
	WriteTime(tag *tags.Tag, date time.Time)
	WriteUint16(tag *tags.Tag, val uint16)
	WriteUint32(tag *tags.Tag, val uint32)
	WriteString(tag *tags.Tag, content string)
	WriteUint16GE(group uint16, element uint16, vr string, val uint16)
	WriteUint32GE(group uint16, element uint16, vr string, val uint32)
	WriteStringGE(group uint16, element uint16, vr string, content string)
	GetTransferSyntax() *transfersyntax.TransferSyntax
	SetTransferSyntax(ts *transfersyntax.TransferSyntax)
	ChangeTransferSynx(ts *transfersyntax.TransferSyntax) error
	TagCount() int
	CreateSR(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string)
	CreatePDF(study DCMStudy, SeriesInstanceUID string, SOPInstanceUID string, fileName string)
	WriteToBytes() []byte
	WriteToFile(fileName string) error
	dumpSeq(indent int)
	compress(i *int, img []byte, RGB bool, cols uint16, rows uint16, bitss uint16, bitsa uint16, pixelrep uint16, planar uint16, frames uint32, outTS string) error
	uncompress(i int, img []byte, size uint32, frames uint32, bitsa uint16, PhotoInt string) error
}

type dcmObj struct {
	Tags           []*DcmTag
	TransferSyntax *transfersyntax.TransferSyntax
	ExplicitVR     bool
	BigEndian      bool
	SQtag          *DcmTag
}

// NewEmptyDCMObj - Create as an interface to a new empty dcmObj
func NewEmptyDCMObj() DcmObj {
	return &dcmObj{
		Tags:           make([]*DcmTag, 0),
		TransferSyntax: nil,
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          &DcmTag{},
	}
}

// NewDCMObjFromFile - Read from a DICOM file into a DICOM Object
func NewDCMObjFromFile(fileName string) (DcmObj, error) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("ERROR, DcmObj::Read, file does not exist")
		}
		return nil, fmt.Errorf("ERROR: %s", err.Error())
	}

	bufdata, err := NewBufDataFromFile(fileName)
	if err != nil {
		return nil, err
	}

	return parseBufData(bufdata)
}

// NewDCMObjFromBytes - Read from a DICOM bytes into a DICOM Object
func NewDCMObjFromBytes(data []byte) (DcmObj, error) {
	return parseBufData(NewBufDataFromBytes(data))
}

func parseBufData(bufdata BufData) (DcmObj, error) {
	BigEndian := false

	transferSyntax, err := bufdata.ReadMeta()
	if err != nil {
		return nil, err
	}

	obj := &dcmObj{
		Tags:           make([]*DcmTag, 0),
		TransferSyntax: transferSyntax,
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          &DcmTag{},
	}

	if obj.TransferSyntax == nil {
		return nil, fmt.Errorf("unable to read transfer syntax from data")
	}

	if obj.TransferSyntax == transfersyntax.ImplicitVRLittleEndian {
		obj.ExplicitVR = false
	} else {
		obj.ExplicitVR = true
	}
	if obj.TransferSyntax == transfersyntax.ExplicitVRBigEndian {
		BigEndian = true
	}
	bufdata.SetBigEndian(BigEndian)

	if err := bufdata.ReadObj(obj); err != nil {
		return nil, err
	}

	return obj, nil
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

// GetTagAt - return the Tag at position i
func (obj *dcmObj) GetTagAt(i int) *DcmTag {
	return obj.Tags[i]
}

func (obj *dcmObj) GetTag(tag *tags.Tag) *DcmTag {
	for _, t := range obj.Tags {
		if t.Group == tag.Group && t.Element == tag.Element {
			return t
		}
	}
	return nil
}

func (obj *dcmObj) GetTagGE(group uint16, element uint16) *DcmTag {
	for _, t := range obj.Tags {
		if t.Group == group && t.Element == element {
			return t
		}
	}
	return nil
}

func (obj *dcmObj) SetTag(i int, tag *DcmTag) {
	FillTag(tag)
	if i <= obj.TagCount() {
		obj.Tags[i] = tag
	}
}

func (obj *dcmObj) InsertTag(index int, tag *DcmTag) {
	FillTag(tag)
	obj.Tags = append(obj.Tags[:index+1], obj.Tags[index:]...)
	obj.Tags[index] = tag
}

func (obj *dcmObj) GetTags() []*DcmTag {
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
		switch tag.VR {
		case "US":
			fmt.Printf("\t(%04X,%04X) %s - %s : %d\n", tag.Group, tag.Element, tag.VR, tag.Description, binary.LittleEndian.Uint16(tag.Data))
		default:
			fmt.Printf("\t(%04X,%04X) %s - %s : %s\n", tag.Group, tag.Element, tag.VR, tag.Description, tag.Data)
		}
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
		switch tag.VR {
		case "US":
			fmt.Printf("%s(%04X,%04X) %s - %s : %d\n", tabs, tag.Group, tag.Element, tag.VR, tag.Description, binary.LittleEndian.Uint16(tag.Data))
		default:
			fmt.Printf("%s(%04X,%04X) %s - %s : %s\n", tabs, tag.Group, tag.Element, tag.VR, tag.Description, tag.Data)
		}
	}
}

func (obj *dcmObj) GetDate(tag *tags.Tag) time.Time {
	date := obj.GetString(tag)
	data, _ := time.Parse("20060102", date)
	return data
}

func (obj *dcmObj) GetUShort(tag *tags.Tag) uint16 {
	return obj.GetUShortGE(tag.Group, tag.Element)
}

// GetUShortGE - return the Uint16 for this group & element
func (obj *dcmObj) GetUShortGE(group uint16, element uint16) uint16 {
	var i int
	var tag *DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTagAt(i)
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

func (obj *dcmObj) GetUInt(tag *tags.Tag) uint32 {
	return obj.GetUIntGE(tag.Group, tag.Element)
}

// GetUIntGE - return the Uint32 for this group & element
func (obj *dcmObj) GetUIntGE(group uint16, element uint16) uint32 {
	var i int
	var tag *DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTagAt(i)
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

func (obj *dcmObj) GetString(tag *tags.Tag) string {
	return obj.GetStringGE(tag.Group, tag.Element)
}

// GetStringGE - return the String for this group & element
func (obj *dcmObj) GetStringGE(group uint16, element uint16) string {
	var i int
	var tag *DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTagAt(i)
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
func (obj *dcmObj) Add(tag *DcmTag) {
	obj.Tags = append(obj.Tags, tag)
}

func (obj *dcmObj) WriteToBytes() []byte {
	bufdata := NewEmptyBufData()

	if obj.TransferSyntax.UID == transfersyntax.ExplicitVRBigEndian.UID {
		bufdata.SetBigEndian(true)
	}
	SOPClassUID := obj.GetStringGE(0x08, 0x16)
	SOPInstanceUID := obj.GetStringGE(0x08, 0x18)
	bufdata.WriteMeta(SOPClassUID, SOPInstanceUID, obj.TransferSyntax.UID)
	bufdata.WriteObj(obj)
	bufdata.SetPosition(0)
	return bufdata.GetAllBytes()
}

// Wrote - Write a DICOM Object to a DICOM File
func (obj *dcmObj) WriteToFile(fileName string) error {
	data := obj.WriteToBytes()
	return os.WriteFile(fileName, data, 0644)
}

func (obj *dcmObj) WriteDate(tag *tags.Tag, date time.Time) {
	obj.WriteString(tag, date.Format("20060102"))
}

func (obj *dcmObj) WriteDateRange(tag *tags.Tag, startDate time.Time, endDate time.Time) {
	obj.WriteString(tag, fmt.Sprintf("%s-%s", startDate.Format("20060102"), endDate.Format("20060102")))
}

func (obj *dcmObj) WriteTime(tag *tags.Tag, date time.Time) {
	obj.WriteString(tag, date.Format("150405"))
}

func (obj *dcmObj) WriteUint16(tag *tags.Tag, val uint16) {
	obj.WriteUint16GE(tag.Group, tag.Element, tag.VR, val)
}

func (obj *dcmObj) WriteUint32(tag *tags.Tag, val uint32) {
	obj.WriteUint32GE(tag.Group, tag.Element, tag.VR, val)
}

func (obj *dcmObj) WriteString(tag *tags.Tag, content string) {
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

	tag := &DcmTag{
		Group:     group,
		Element:   element,
		Length:    2,
		VR:        vr,
		Data:      c,
		BigEndian: obj.BigEndian,
	}
	FillTag(tag)
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

	tag := &DcmTag{
		Group:     group,
		Element:   element,
		Length:    4,
		VR:        vr,
		Data:      c,
		BigEndian: obj.BigEndian,
	}
	FillTag(tag)
	obj.Tags = append(obj.Tags, tag)
}

// WriteStringGE - Writes a String to a DICOM tag
func (obj *dcmObj) WriteStringGE(group uint16, element uint16, vr string, content string) {
	data := []byte(content)
	length := len(data)
	if length%2 == 1 {
		length++
		if vr == "UI" {
			data = append(data, 0x00)
		} else {
			data = append(data, 0x20)
		}
	}
	tag := &DcmTag{
		Group:     group,
		Element:   element,
		Length:    uint32(length),
		VR:        vr,
		Data:      data,
		BigEndian: false,
	}
	FillTag(tag)
	obj.Tags = append(obj.Tags, tag)
}

func (obj *dcmObj) GetTransferSyntax() *transfersyntax.TransferSyntax {
	return obj.TransferSyntax
}

func (obj *dcmObj) SetTransferSyntax(ts *transfersyntax.TransferSyntax) {
	obj.TransferSyntax = ts
}

func (obj *dcmObj) GetPixelData(frame int) ([]byte, error) {
	var i int
	var rows, cols, bitsa, planar uint16
	var PhotoInt string
	sq := 0
	frames := uint32(0)
	RGB := false
	icon := false

	if !transfersyntax.SupportedTransferSyntax(obj.TransferSyntax.UID) {
		return nil, fmt.Errorf("unsupported transfer synxtax %s", obj.TransferSyntax.Name)
	}

	for i = 0; i < len(obj.Tags); i++ {
		tag := obj.GetTagAt(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if sq == 0 {
			if (tag.Group == 0x0028) && (!icon) {
				switch tag.Element {
				case 0x04:
					PhotoInt = tag.GetString()
					if !strings.Contains(PhotoInt, "MONO") {
						RGB = true
					}
				case 0x06:
					planar = tag.GetUShort()
				case 0x08:
					uframes, err := strconv.Atoi(tag.GetString())
					if err != nil {
						frames = 0
					} else {
						frames = uint32(uframes)
					}
				case 0x10:
					rows = tag.GetUShort()
				case 0x11:
					cols = tag.GetUShort()
				case 0x0100:
					bitsa = tag.GetUShort()
				}
			}
			if (tag.Group == 0x0088) && (tag.Element == 0x0200) && (tag.Length == 0xFFFFFFFF) {
				icon = true
			}
			if (tag.Group == 0x6003) && (tag.Element == 0x1010) && (tag.Length == 0xFFFFFFFF) {
				icon = true
			}
			if (tag.Group == 0x7FE0) && (tag.Element == 0x0010) && (!icon) {
				size := uint32(cols) * uint32(rows) * uint32(bitsa) / 8
				if RGB {
					size = 3 * size
				}
				if frames > 0 {
					size = uint32(frames) * size
				} else {
					frames = 1
				}
				if size == 0 {
					return nil, errors.New("ERROR, DcmObj::ConvertTransferSyntax, size=0")
				}

				if frame > int(frames) {
					return nil, errors.New("ERROR, invalid frame")
				}

				if tag.Length == 0xFFFFFFFF {
					return obj.GetTagAt(i + 2 + frame).Data, nil
				} else {
					if RGB && (planar == 1) {
						var img_offset, img_size uint32
						img_size = size / frames
						img := make([]byte, img_size)
						for f := uint32(0); f < frames; f++ {
							img_offset = img_size * f
							for j := uint32(0); j < img_size/3; j++ {
								img[3*j] = tag.Data[j+img_offset]
								img[3*j+1] = tag.Data[j+img_size/3+img_offset]
								img[3*j+2] = tag.Data[j+2*img_size/3+img_offset]
							}
							if f == uint32(frame) {
								return img, nil
							}
						}
						planar = 0
					} else {
						return tag.Data, nil
					}
				}
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	return nil, fmt.Errorf("there was an error getting pixel data")
}

func (obj *dcmObj) ChangeTransferSynx(outTS *transfersyntax.TransferSyntax) error {
	flag := false

	var i int
	var rows, cols, bitss, bitsa, planar, pixelrep uint16
	var PhotoInt string
	sq := 0
	frames := uint32(0)
	RGB := false
	icon := false

	if obj.TransferSyntax.UID == outTS.UID {
		return nil
	}

	if !transfersyntax.SupportedTransferSyntax(outTS.UID) {
		return fmt.Errorf("unsupported transfer synxtax %s", outTS.Name)
	}

	for i = 0; i < len(obj.Tags); i++ {
		tag := obj.GetTagAt(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if sq == 0 {
			if (tag.Group == 0x0028) && (!icon) {
				switch tag.Element {
				case 0x04:
					PhotoInt = tag.GetString()
					if !strings.Contains(PhotoInt, "MONO") {
						RGB = true
					}
				case 0x06:
					planar = tag.GetUShort()
				case 0x08:
					uframes, err := strconv.Atoi(tag.GetString())
					if err != nil {
						frames = 0
					} else {
						frames = uint32(uframes)
					}
				case 0x10:
					rows = tag.GetUShort()
				case 0x11:
					cols = tag.GetUShort()
				case 0x0100:
					bitsa = tag.GetUShort()
				case 0x0101:
					bitss = tag.GetUShort()
				case 0x0103:
					pixelrep = tag.GetUShort()
				}
			}
			if (tag.Group == 0x0088) && (tag.Element == 0x0200) && (tag.Length == 0xFFFFFFFF) {
				icon = true
			}
			if (tag.Group == 0x6003) && (tag.Element == 0x1010) && (tag.Length == 0xFFFFFFFF) {
				icon = true
			}
			if (tag.Group == 0x7FE0) && (tag.Element == 0x0010) && (!icon) {
				size := uint32(cols) * uint32(rows) * uint32(bitsa) / 8
				if RGB {
					size = 3 * size
				}
				if frames > 0 {
					size = uint32(frames) * size
				} else {
					frames = 1
				}
				if size == 0 {
					return errors.New("ERROR, DcmObj::ConvertTransferSyntax, size=0")
				}
				img := make([]byte, size)
				if tag.Length == 0xFFFFFFFF {
					obj.uncompress(i, img, size, frames, bitsa, PhotoInt)
				} else { // Uncompressed
					if RGB && (planar == 1) { // change from planar=1 to planar=0
						var img_offset, img_size uint32
						img_size = size / frames
						for f := uint32(0); f < frames; f++ {
							img_offset = img_size * f
							for j := uint32(0); j < img_size/3; j++ {
								img[3*j+img_offset] = tag.Data[j+img_offset]
								img[3*j+1+img_offset] = tag.Data[j+img_size/3+img_offset]
								img[3*j+2+img_offset] = tag.Data[j+2*img_size/3+img_offset]
							}
						}
						planar = 0
					} else {
						copy(img, tag.Data)
					}
				}
				if err := obj.compress(&i, img, RGB, cols, rows, bitss, bitsa, pixelrep, planar, frames, outTS.UID); err != nil {
					return err
				} else {
					flag = true
				}
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	if flag {
		obj.TransferSyntax = outTS
		return nil
	}
	return fmt.Errorf("there was an error changing the transfer synxtax")
}

// AddConceptNameSeq - Concept Name Sequence for DICOM SR
func (obj *dcmObj) AddConceptNameSeq(group uint16, element uint16, CodeValue string, CodeMeaning string) {
	item := &dcmObj{
		Tags:           make([]*DcmTag, 0),
		TransferSyntax: nil,
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          new(DcmTag),
	}
	seq := &dcmObj{
		Tags:           make([]*DcmTag, 0),
		TransferSyntax: nil,
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          new(DcmTag),
	}
	tag := new(DcmTag)

	item.BigEndian = obj.BigEndian
	item.ExplicitVR = obj.ExplicitVR
	seq.BigEndian = obj.BigEndian
	seq.ExplicitVR = obj.ExplicitVR

	item.WriteString(tags.CodeValue, CodeValue)
	item.WriteString(tags.CodingSchemeDesignator, "odb")
	item.WriteString(tags.CodeMeaning, CodeMeaning)
	tag.WriteSeq(0xFFFE, 0xE000, item)
	seq.Add(tag)
	tag.WriteSeq(group, element, seq)
	obj.Add(tag)
}

// AddSRText - add Text to SR
func (obj *dcmObj) AddSRText(text string) {
	item := &dcmObj{
		Tags:           make([]*DcmTag, 0),
		TransferSyntax: nil,
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          new(DcmTag),
	}
	seq := &dcmObj{
		Tags:           make([]*DcmTag, 0),
		TransferSyntax: nil,
		ExplicitVR:     false,
		BigEndian:      false,
		SQtag:          new(DcmTag),
	}
	tag := new(DcmTag)

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
	obj.WriteString(tags.SOPClassUID, sopclass.BasicTextSRStorage.UID)
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
	obj.WriteString(tags.SOPClassUID, sopclass.EncapsulatedPDFStorage.UID)
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
		mstream.Append([]byte{0x00})
	}
	obj.WriteString(tags.DocumentTitle, fileName)
	obj.Add(&DcmTag{
		Group:     0x42,
		Element:   0x11,
		Length:    size,
		VR:        "OB",
		Data:      mstream.GetData(),
		BigEndian: obj.BigEndian,
	})
	obj.WriteString(tags.MIMETypeOfEncapsulatedDocument, "application/pdf")
}

func (obj *dcmObj) compress(i *int, img []byte, RGB bool, cols uint16, rows uint16, bitss uint16, bitsa uint16, pixelrep uint16, planar uint16, frames uint32, outTS string) error {
	var offset, size, jpeg_size, j uint32
	var JPEGData []byte
	var JPEGBytes, index int

	single := uint32(cols) * uint32(rows) * uint32(bitsa) / 8
	size = single * frames
	if RGB {
		size = 3 * size
	}

	index = *i
	tag := obj.GetTagAt(index)

	switch outTS {
	case transfersyntax.JPEGLosslessSV1.UID:
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := &DcmTag{
			Group:     0xFFFE,
			Element:   0xE000,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		for j = 0; j < frames; j++ {
			index++
			offset = j * uint32(cols) * uint32(rows) * uint32(bitsa) / 8
			if RGB {
				offset = 3 * offset
			}
			if bitsa == 8 {
				if RGB {
					if err := jpeglib.EIJG8encode(img[offset:], cols, rows, 3, &JPEGData, &JPEGBytes, 4); err != nil {
						return err
					}
				} else {
					if err := jpeglib.EIJG8encode(img[offset:], cols, rows, 1, &JPEGData, &JPEGBytes, 4); err != nil {
						return err
					}
				}
			} else {
				if err := jpeglib.EIJG16encode(img[offset/2:], cols, rows, 1, &JPEGData, &JPEGBytes, 0); err != nil {
					return err
				}
			}
			newtag = &DcmTag{
				Group:     0xFFFE,
				Element:   0xE000,
				Length:    uint32(JPEGBytes),
				VR:        "DL",
				Data:      JPEGData,
				BigEndian: obj.IsBigEndian(),
			}
			obj.InsertTag(index, newtag)
			JPEGData = nil
		}
		index++
		newtag = &DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	case transfersyntax.JPEGBaseline8Bit.UID:
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := &DcmTag{
			Group:     0xFFFE,
			Element:   0xE000,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		jpeg_size = 0
		for j = 0; j < frames; j++ {
			index++
			offset = j * uint32(cols) * uint32(rows) * uint32(bitsa) / 8
			if RGB {
				offset = 3 * offset
				if err := jpeglib.EIJG8encode(img[offset:], cols, rows, 3, &JPEGData, &JPEGBytes, 0); err != nil {
					return err
				}
			} else {
				if bitsa == 8 {
					if err := jpeglib.EIJG8encode(img[offset:], cols, rows, 1, &JPEGData, &JPEGBytes, 0); err != nil {
						return err
					}
				} else {
					if err := jpeglib.EIJG12encode(img[offset:], cols, rows, 1, &JPEGData, &JPEGBytes, 0); err != nil {
						return err
					}
				}
			}
			newtag = &DcmTag{
				Group:     0xFFFE,
				Element:   0xE000,
				Length:    uint32(JPEGBytes),
				VR:        "DL",
				Data:      JPEGData,
				BigEndian: obj.IsBigEndian(),
			}
			obj.InsertTag(index, newtag)
			JPEGData = nil
			jpeg_size = jpeg_size + uint32(JPEGBytes)
		}
		index++
		newtag = &DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	case transfersyntax.JPEGExtended12Bit.UID:
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := &DcmTag{
			Group:     0xFFFE,
			Element:   0xE000,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		jpeg_size = 0
		for j = 0; j < frames; j++ {
			index++
			offset = j * uint32(cols) * uint32(rows) * uint32(bitsa) / 8
			if err := jpeglib.EIJG12encode(img[offset/2:], cols, rows, 1, &JPEGData, &JPEGBytes, 0); err != nil {
				return err
			}
			newtag = &DcmTag{
				Group:     0xFFFE,
				Element:   0xE000,
				Length:    uint32(JPEGBytes),
				VR:        "DL",
				Data:      JPEGData,
				BigEndian: obj.IsBigEndian(),
			}
			obj.InsertTag(index, newtag)
			JPEGData = nil
			jpeg_size = jpeg_size + uint32(JPEGBytes)
		}
		index++
		newtag = &DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	case transfersyntax.JPEG2000Lossless.UID:
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := &DcmTag{
			Group:     0xFFFE,
			Element:   0xE000,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		for j = 0; j < frames; j++ {
			index++
			offset = j * uint32(cols) * uint32(rows) * uint32(bitsa) / 8
			if RGB {
				offset = 3 * offset
				if err := openjpeg.J2Kencode(img[offset:], cols, rows, 3, bitsa, &JPEGData, &JPEGBytes, 0); err != nil {
					return err
				}
			} else {
				if err := openjpeg.J2Kencode(img[offset:], cols, rows, 1, bitsa, &JPEGData, &JPEGBytes, 0); err != nil {
					return err
				}
			}
			newtag = &DcmTag{
				Group:     0xFFFE,
				Element:   0xE000,
				Length:    uint32(JPEGBytes),
				VR:        "DL",
				Data:      JPEGData,
				BigEndian: obj.IsBigEndian(),
			}
			obj.InsertTag(index, newtag)
			JPEGData = nil
		}
		index++
		newtag = &DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	case transfersyntax.JPEG2000.UID:
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := &DcmTag{
			Group:     0xFFFE,
			Element:   0xE000,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		jpeg_size = 0
		for j = 0; j < frames; j++ {
			index++
			offset = j * uint32(cols) * uint32(rows) * uint32(bitsa) / 8
			if RGB {
				offset = 3 * offset
				if err := openjpeg.J2Kencode(img[offset:], cols, rows, 3, bitsa, &JPEGData, &JPEGBytes, 10); err != nil {
					return err
				}
			} else {
				if err := openjpeg.J2Kencode(img[offset:], cols, rows, 1, bitsa, &JPEGData, &JPEGBytes, 10); err != nil {
					return err
				}
			}
			newtag = &DcmTag{
				Group:     0xFFFE,
				Element:   0xE000,
				Length:    uint32(JPEGBytes),
				VR:        "DL",
				Data:      JPEGData,
				BigEndian: obj.IsBigEndian(),
			}
			obj.InsertTag(index, newtag)
			JPEGData = nil
			jpeg_size = jpeg_size + uint32(JPEGBytes)
		}
		index++
		newtag = &DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	default:
		if bitss == 8 {
			tag.VR = "OB"
		} else {
			tag.VR = "OW"
		}
		tag.Length = size
		if tag.Data != nil {
			tag.Data = nil
		}
		tag.Data = make([]byte, tag.Length)
		copy(tag.Data, img)
		obj.SetTag(index, tag)
	}
	return nil
}

func (obj *dcmObj) uncompress(i int, img []byte, size uint32, frames uint32, bitsa uint16, PhotoInt string) error {
	var j, offset, single uint32
	single = size / frames

	obj.DelTag(i + 1) // Delete offset table.
	switch obj.TransferSyntax.UID {
	case transfersyntax.RLELossless.UID:
		for j = 0; j < frames; j++ {
			offset = j * single
			tag := obj.GetTagAt(i + 1)
			if err := transcoder.RLEdecode(tag.Data, img[offset:], tag.Length, single, PhotoInt); err != nil {
				return err
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	case transfersyntax.JPEGLosslessSV1.UID:
	case transfersyntax.JPEGLossless.UID:
		for j = 0; j < frames; j++ {
			offset = j * single
			tag := obj.GetTagAt(i + 1)
			if bitsa == 8 {
				if err := jpeglib.DIJG8decode(tag.Data, tag.Length, img[offset:], single); err != nil {
					return err
				}
			} else {
				if err := jpeglib.DIJG16decode(tag.Data, tag.Length, img[offset:], single); err != nil {
					return err
				}
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	case transfersyntax.JPEGBaseline8Bit.UID:
		for j = 0; j < frames; j++ {
			offset = j * single
			tag := obj.GetTagAt(i + 1)
			if bitsa == 8 {
				if err := jpeglib.DIJG8decode(tag.Data, tag.Length, img[offset:], single); err != nil {
					return err
				}
			} else {
				if err := jpeglib.DIJG12decode(tag.Data, tag.Length, img[offset:], single); err != nil {
					return err
				}
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	case transfersyntax.JPEGExtended12Bit.UID:
		for j = 0; j < frames; j++ {
			offset = j * single
			tag := obj.GetTagAt(i + 1)
			if err := jpeglib.DIJG12decode(tag.Data, tag.Length, img[offset:], single); err != nil {
				return err
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	case transfersyntax.JPEG2000Lossless.UID:
		for j = 0; j < frames; j++ {
			offset = j * single
			tag := obj.GetTagAt(i + 1)
			if err := openjpeg.J2Kdecode(tag.Data, tag.Length, img[offset:]); err != nil {
				return err
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	case transfersyntax.JPEG2000.UID:
		for j = 0; j < frames; j++ {
			offset = j * single
			tag := obj.GetTagAt(i + 1)
			if err := openjpeg.J2Kdecode(tag.Data, tag.Length, img[offset:]); err != nil {
				return err
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	}
	return nil
}
