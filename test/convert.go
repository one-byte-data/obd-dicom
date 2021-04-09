package main

import (
	"encoding/binary"
	"git.onebytedata.com/odb/go-dicom/jpeglib"
	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/openjpeg"
	"log"
	"strconv"
	"strings"
)

func SupportedTS(TransferSyntax string) bool {
	return true
}

func GetUint32(in []byte, length int) uint32 {
	c := make([]byte, length)
	copy(c, in)
	return binary.LittleEndian.Uint32(c)
}

func ReadSegment(in []byte, out []byte, seg_offset uint32, seg_size uint32, i uint32, rawSize uint32) {
	var count int8
	out_offset := i * rawSize
	in_offset := seg_offset

	for (out_offset - i*rawSize) < rawSize {
		count = int8(in[in_offset])
		in_offset++
		if count >= 0 {
			copy(out[out_offset:out_offset+uint32(count+1)], in[in_offset:in_offset+uint32(count+1)])
			in_offset += uint32(count + 1)
			out_offset += uint32(count + 1)
		} else {
			if (count <= -1) && (count >= -127) {
				newByte := in[in_offset]
				in_offset++
				for j := uint32(0); j < uint32(-count+1); j++ {
					out[j+out_offset] = newByte
				}
				out_offset += uint32(-count + 1)
				if in_offset-seg_offset > seg_size {
					log.Println("ERROR, overflow decoding RLE")
					return
				}
			}
		}
	}
}

func RLEdecode(in []byte, out []byte, length uint32, size uint32, PhotoInt string) {
	var segment_count, offset, i uint32
	var segment_offset [15]uint32
	var segment_length [15]uint32

	offset = 0
	for i := 0; i < 15; i++ {
		segment_offset[i] = 0
		segment_length[i] = 0
	}

	segment_count = GetUint32(in, 4)
	for i := 0; i < 15; i++ {
		offset = offset + 4
		segment_offset[i] = GetUint32(in[offset:], 4)
	}

	segment_offset[segment_count] = length
	temp := make([]byte, size)
	for i = 0; i < segment_count; i++ {
		segment_length[i] = segment_offset[i+1] - segment_offset[i]
		ReadSegment(in, temp, segment_offset[i], segment_length[i], i, size/segment_count)
	}

	offset = size / segment_count
	if (strings.Contains(PhotoInt, "MONO")) && (segment_count == 2) {
		for i = 0; i < size/segment_count; i++ {
			out[2*i] = temp[i+offset]
			out[2*i+1] = temp[i]
		}
	} else if (strings.Contains(PhotoInt, "MONO")) && (segment_count == 1) {
		for i = 0; i < size; i++ {
			out[i] = temp[i]
		}
	} else if (PhotoInt == "YBR_FULL") && (segment_count == 3) {
		var Y, Cb, Cr float32
		for i = 0; i < size/segment_count; i++ {
			Y = float32(temp[i])
			Cb = float32(temp[i+offset])
			Cr = float32(temp[i+2*offset])
			out[3*i] = byte(Y + 1.402*(Cr-128.0))
			out[3*i+1] = byte(Y - 0.344136*(Cb-128.0) - 0.714136*(Cr-128.0))
			out[3*i+2] = byte(Y + 1.772*(Cb-128.0))
		}
	} else if (PhotoInt == "RGB") && (segment_count == 3) {
		for i = 0; i < size/segment_count; i++ {
			out[3*i] = temp[i]
			out[3*i+1] = temp[i+offset]
			out[3*i+2] = temp[i+2*offset]
		}
	} else {
		log.Println("ERROR, format not supported")
	}
}

func Decomp(obj media.DcmObj, i int, img []byte, size uint32, frames uint32, bitsa uint16, PhotoInt string) {
	var tag media.DcmTag
	var j, offset, single uint32

	single = size / frames
	// DE-Compression
	obj.DelTag(i + 1) // Delete offset table.
	if obj.GetTransferSyntax() == "1.2.840.10008.1.2.5" {
		for j = 0; j < frames; j++ {
			offset = j * single
			tag = obj.GetTag(i + 1)
			RLEdecode(tag.Data, img[offset:], tag.Length, single, PhotoInt)
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	} else if (obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.70") || (obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.57") {
		for j = 0; j < frames; j++ {
			offset = j * single
			tag = obj.GetTag(i + 1)
			if bitsa == 8 {
				jpeglib.DIJG8decode(tag.Data, tag.Length, img[offset:], single)
			} else {
				jpeglib.DIJG16decode(tag.Data, tag.Length, img[offset:], single)
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	} else if obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.50" {
		for j = 0; j < frames; j++ {
			offset = j * single
			tag = obj.GetTag(i + 1)
			if bitsa == 8 {
				jpeglib.DIJG8decode(tag.Data, tag.Length, img[offset:], single)
			} else {
				jpeglib.DIJG12decode(tag.Data, tag.Length, img[offset:], single)
			}
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	} else if obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.51" {
		for j = 0; j < frames; j++ {
			offset = j * single
			tag = obj.GetTag(i + 1)
			jpeglib.DIJG12decode(tag.Data, tag.Length, img[offset:], single)
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	} else if obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.90" {
		for j = 0; j < frames; j++ {
			offset = j * single
			tag = obj.GetTag(i + 1)
			openjpeg.J2Kdecode(tag.Data, tag.Length, img[offset:])
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	} else if obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.91" {
		for j = 0; j < frames; j++ {
			offset = j * single
			tag = obj.GetTag(i + 1)
			openjpeg.J2Kdecode(tag.Data, tag.Length, img[offset:])
			obj.DelTag(i + 1)
		}
		obj.DelTag(i + 1)
	}
}

func Comp(obj media.DcmObj, i *int, img []byte, RGB bool, cols uint16, rows uint16, bitss uint16, bitsa uint16, pixelrep uint16, planar uint16, frames uint32, outTS string) bool {
	var tag media.DcmTag
	var offset, size, jpeg_size, j uint32
	var JPEGData []byte
	var JPEGBytes, index int

	single := uint32(cols) * uint32(rows) * uint32(bitsa) / 8
	size = single * frames
	if RGB {
		size = 3 * size
	}

	index = *i
	tag = obj.GetTag(index)
	if outTS == "1.2.840.10008.1.2.4.70" {
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := media.DcmTag{
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
					jpeglib.EIJG8encode(img[offset:], cols, rows, 3, &JPEGData, &JPEGBytes, 4)
				} else {
					jpeglib.EIJG8encode(img[offset:], cols, rows, 1, &JPEGData, &JPEGBytes, 4)
				}
			} else {
				jpeglib.EIJG16encode(img[offset/2:], cols, rows, 1, &JPEGData, &JPEGBytes, 0)
			}
			newtag = media.DcmTag{
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
		newtag = media.DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	} else if outTS == "1.2.840.10008.1.2.4.50" {
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := media.DcmTag{
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
				jpeglib.EIJG8encode(img[offset:], cols, rows, 3, &JPEGData, &JPEGBytes, 0)
			} else {
				if bitsa == 8 {
					jpeglib.EIJG8encode(img[offset:], cols, rows, 1, &JPEGData, &JPEGBytes, 0)
				} else { // ERROR...
					// Can't use this transfer Syntax with bitsa!=8
					return false
				}
			}
			newtag = media.DcmTag{
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
		newtag = media.DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	} else if outTS == "1.2.840.10008.1.2.4.51" {
		if (bitss == 8) && (bitsa != 16) {
			return false
		}
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := media.DcmTag{
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
			if bitss > 12 {
				return false
			}
			jpeglib.EIJG12encode(img[offset/2:], cols, rows, 1, &JPEGData, &JPEGBytes, 0)
			newtag = media.DcmTag{
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
		newtag = media.DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	} else if outTS == "1.2.840.10008.1.2.4.90" {
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := media.DcmTag{
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
				openjpeg.J2Kencode(img[offset:], cols, rows, 3, bitsa, &JPEGData, &JPEGBytes, 0)
			} else {
				openjpeg.J2Kencode(img[offset:], cols, rows, 1, bitsa, &JPEGData, &JPEGBytes, 0)
			}
			newtag = media.DcmTag{
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
		newtag = media.DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	} else if outTS == "1.2.840.10008.1.2.4.91" {
		tag.VR = "OB"
		tag.Length = 0xFFFFFFFF
		if tag.Data != nil {
			tag.Data = nil
		}
		obj.SetTag(index, tag)
		index++
		newtag := media.DcmTag{
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
				openjpeg.J2Kencode(img[offset:], cols, rows, 3, bitsa, &JPEGData, &JPEGBytes, 10)
			} else {
				openjpeg.J2Kencode(img[offset:], cols, rows, 1, bitsa, &JPEGData, &JPEGBytes, 10)
			}
			newtag = media.DcmTag{
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
		newtag = media.DcmTag{
			Group:     0xFFFE,
			Element:   0xE0DD,
			Length:    0,
			VR:        "DL",
			Data:      nil,
			BigEndian: obj.IsBigEndian(),
		}
		obj.InsertTag(index, newtag)
		*i = index
	} else {
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
	return true
}

func ConvertTS(obj media.DcmObj, outTS string) bool {
	flag := false
	//	ExplicitVROUT:=true
	var i int
	var tag media.DcmTag
	var rows, cols, bitss, bitsa, planar, pixelrep uint16
	var PhotoInt string
	sq := 0
	frames := uint32(0)
	RGB := false
	icon := false

	if len(outTS) == 0 {
		return true
	}
	if obj.GetTransferSyntax() == outTS {
		return true
	}
	// We don't process MPEG2 or MPEG4
	if (obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.100") || (obj.GetTransferSyntax() == "1.2.840.10008.1.2.4.102") {
		return true
	}
	if !SupportedTS(obj.GetTransferSyntax()) {
		return false
	}
	if !SupportedTS(outTS) {
		return false
	}
	if outTS == "1.2.840.10008.1.2.5" {
		return false
	}

	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if sq == 0 {
			if (tag.Group == 0x0028) && (!icon) {
				switch tag.Element {
				case 0x04:
					PhotoInt = tag.GetString()
					if strings.Contains(PhotoInt, "MONO") == false {
						RGB = true
					}
					break
				case 0x06:
					planar = tag.GetUShort()
					break
				case 0x08:
					uframes, err := strconv.Atoi(tag.GetString())
					if err != nil {
						frames = 0
					} else {
						frames = uint32(uframes)
					}
					break
				case 0x10:
					rows = tag.GetUShort()
					break
				case 0x11:
					cols = tag.GetUShort()
					break
				case 0x0100:
					bitsa = tag.GetUShort()
					break
				case 0x0101:
					bitss = tag.GetUShort()
					break
				case 0x0103:
					pixelrep = tag.GetUShort()
					break
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
					log.Println("ERROR, DcmObj::ConvertTransferSyntax, size=0")
					return false
				}
				img := make([]byte, size)
				if tag.Length == 0xFFFFFFFF {
					Decomp(obj, i, img, size, frames, bitsa, PhotoInt)
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
				flag = Comp(obj, &i, img, RGB, cols, rows, bitss, bitsa, pixelrep, planar, frames, outTS)
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	if flag {
		obj.SetTransferSyntax(outTS)
	}
	return flag
}

/*
func main() {
	media.InitDict()
	obj, err := media.NewDCMObjFromFile("images/rle_gray.dcm")
	if err != nil {
		log.Panic(err)
	}
	if ConvertTS(obj, "1.2.840.10008.1.2.1") {
		obj.WriteToFile("out.dcm")
	}
}
*/
