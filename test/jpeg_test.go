package main

import (
	"fmt"
	"log"
	"os"

	"git.onebytedata.com/odb/go-dicom/jpeglib"
	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/tags"
	"git.onebytedata.com/odb/go-dicom/uid"
)

// LoadFromFile - Load from File into MemoryStream
func LoadFromFile(FileName string, buffer *[]byte) bool {
	flag := false

	file, err := os.Open(FileName)
	if err != nil {
		fmt.Println("ERROR, opening file: " + FileName)
		return flag
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("ERROR, getting file Stats: " + FileName)
		return flag
	}
	size := int(stat.Size())
	bs := make([]byte, size)
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println("ERROR, reading file: " + FileName)
		return flag
	}
	*buffer = append(*buffer, bs...)
	return true
}

// SaveToFile - Save MemoryStream to File
func SaveToFile(FileName string, buffer []byte) bool {
	flag := false

	file, err := os.Create(FileName)
	if err != nil {
		fmt.Println("ERROR, opening file: " + FileName)
		return flag
	}
	defer file.Close()
	_, err = file.Write(buffer)
	if err != nil {
		fmt.Println("ERROR, writing to file: " + FileName)
		return flag
	}
	return true
}

// GetPixelData - Gets pixel data from DcmObj
func GetPixelData(obj media.DcmObj, index *int) []uint8 {
	var i int
	tag := new(media.DcmTag)
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if (sq == 0) && (tag.Length > 0) && (tag.Length != 0xFFFFFFFF) {
			if (tag.Group == 0x7fe0) && (tag.Element == 0x10) {
				*index = i
				return tag.Data
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	*index = -1
	return nil
}

func insert(a []media.DcmTag, index int, value media.DcmTag) []media.DcmTag {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func test16() {
	obj, err := media.NewDCMObjFromFile("test.dcm")
	if err != nil {
		log.Panic(err)
	}

	// Need an uncompressed image
	if obj.GetTransferSyntax().UID == uid.ExplicitVRLittleEndian.UID {
		var index int
		width := obj.GetUShort(tags.Rows)
		height := obj.GetUShort(tags.Columns)
		pixelData := GetPixelData(obj, &index)
		if len(pixelData) > 0 {
			var outData []byte
			var outSize int
			// Encode image
			err := jpeglib.EIJG16encode(pixelData, width, height, 1, &outData, &outSize, 0)
			if err != nil {
				log.Panic(err)
			}

			tag := obj.GetTag(index)
			tag.Length = 0xFFFFFFFF
			tag.VR = "OB"
			tag.Data = nil
			obj.InsertTag(index, tag)
			obj.SetTransferSyntax(uid.JPEGExtendedHierarchical1719)
			index++
			tag = &media.DcmTag{
				Group:     0xFFFE,
				Element:   0xE000,
				Length:    0,
				VR:        "DL",
				Data:      nil,
				BigEndian: false,
			}
			obj.InsertTag(index, tag)
			index++
			tag = &media.DcmTag{
				Group:     0xFFFE,
				Element:   0xE000,
				Length:    uint32(outSize),
				VR:        "DL",
				Data:      outData,
				BigEndian: false,
			}
			obj.InsertTag(index, tag)
			index++
			tag = &media.DcmTag{
				Group:     0xFFFE,
				Element:   0xE0DD,
				Length:    0,
				VR:        "DL",
				Data:      nil,
				BigEndian: false,
			}
			obj.InsertTag(index, tag)

			err = obj.WriteToFile("out.dcm")
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

func test8() {
	var jpegData []byte
	var outData []byte
	var jpegSize int

	if LoadFromFile("test.jpg", &jpegData) {
		outSize := 1576 * 1134 * 3 // Image Size, have to know in advance.
		outData = make([]byte, outSize)
		err := jpeglib.DIJG8decode(jpegData, uint32(len(jpegData)), outData, uint32(outSize))
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("Decode Success!")
		if SaveToFile("out.raw", outData) {
			fmt.Println("Saved out.raw")
		} else {
			fmt.Println("ERROR, Saving out.raw")
		}
	} else {
		fmt.Println("ERROR, Decode Failed!")
	}

	jpegData = nil
	outData = nil
	if LoadFromFile("test.raw", &outData) {
		err := jpeglib.EIJG8encode(outData, 1576, 1134, 3, &jpegData, &jpegSize, 4)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("Encode Success!")
		if SaveToFile("out.jpg", jpegData) {
			fmt.Println("Saved out.jpg")
		} else {
			fmt.Println("ERROR, Saving out.jpg")
		}
	} else {
		fmt.Println("ERROR, Encode Failed!")
	}
}

/*
func main(){
	test8()
	test16()
}
*/
