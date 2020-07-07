package main

import (
	"fmt"
	"os"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
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

func GetPixelData(obj media.DcmObj, index *int) []uint8 {
	var i int
	var tag media.DcmTag
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
	var obj media.DcmObj

	if obj.Read("test.dcm") {
		// Need an uncompressed image
		if obj.TransferSyntax == "1.2.840.10008.1.2.1" {
			var index int
			width := int(obj.GetUShort(0x28, 0x10))
			height := int(obj.GetUShort(0x28, 0x11))
			pixelData := GetPixelData(obj, &index)
			if len(pixelData) > 0 {
				var out_data []byte
				// Encode image
				if Encode16(pixelData, width, height, 1, &out_data) == true {
					tag := obj.GetTag(index)
					tag.Length = 0xFFFFFFFF
					tag.VR = "OB"
					tag.Data = nil
					obj.Tags[index] = tag
					obj.TransferSyntax = "1.2.840.10008.1.2.4.70"
					index++
					tag = media.DcmTag{0xFFFE, 0xE000, 0, "DL", nil, false}
					obj.Tags = insert(obj.Tags, index, tag)
					index++
					tag = media.DcmTag{0xFFFE, 0xE000, uint32(len(out_data)), "DL", out_data, false}
					obj.Tags = insert(obj.Tags, index, tag)
					index++
					tag = media.DcmTag{0xFFFE, 0xE0DD, 0, "DL", nil, false}
					obj.Tags = insert(obj.Tags, index, tag)
					obj.Write("out.dcm")
				}
			}
		}
	}
}

func test8() {
	var jpeg_data []byte
	var out_data []byte

	if LoadFromFile("test.jpg", &jpeg_data) {
		out_size := 1576 * 1134 * 3 // Image Size, have to know in advance.
		out_data = make([]byte, out_size)

		if Decode8(jpeg_data, len(jpeg_data), out_data, out_size) {
			fmt.Println("Decode Success!")
			if SaveToFile("out.raw", out_data) {
				fmt.Println("Saved out.raw")
			} else {
				fmt.Println("ERROR, Saving out.raw")
			}
		} else {
			fmt.Println("ERROR, Decode Failed!")
		}
	}
	jpeg_data = nil
	out_data = nil
	if LoadFromFile("test.raw", &out_data) {
		if Encode8(out_data, 1576, 1134, 3, &jpeg_data) {
			fmt.Println("Encode Success!")
			if SaveToFile("out.jpg", jpeg_data) {
				fmt.Println("Saved out.jpg")
			} else {
				fmt.Println("ERROR, Saving out.jpg")
			}
		} else {
			fmt.Println("ERROR, Encode Failed!")
		}
	}
}
