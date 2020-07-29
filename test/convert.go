package main

import (
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"strconv"
)

func SupportedTS(TransferSyntax string) bool {
	return true
}

func ConvertTS(obj media.DcmObj, outTS string) bool {
	flag:=false
	ExplicitVROUT:=true
	var i int
	var tag media.DcmTag
	var rows, cols, bitss, bitsa, planar, pixelrep uint16
	var PhotometricInterpretation string
	sq := 0
	frames:=0
	RGB :=false
	icon:=false

	if len(outTS)==0 {
		return true
	}
	if obj.GetTransferSynxtax() == outTS {
		return true
	}
	// We don't process MPEG2 or MPEG4
	if (obj.GetTransferSynxtax() == "1.2.840.10008.1.2.4.100") || (obj.GetTransferSynxtax() == "1.2.840.10008.1.2.4.102") {
		return true
	}
	if !SupportedTS(obj.GetTransferSynxtax()) {
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
		if (sq == 0){
			if (tag.Group==0x0028)&&(!icon){
				switch(tag.Element){
					 case 0x04:
						PhotometricInterpretation=tag.GetString()
						break
					 case 0x06:
						planar = tag.GetUShort()
						break
					 case 0x08:
						frames, err := strconv.Atoi(tag.GetString())
						if err != nil {
							frames = 0
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
				if (tag.Group==0x0088)&&(tag.Element==0x0200)&&(tag.Length==0xFFFFFFFF) {
					icon=true
				}
				if (tag.Group==0x6003)&&(tag.Element==0x1010)&&(tag.Length==0xFFFFFFFF) {
					icon=true
				}
				if (tag.Group==0x7FE0)&&(tag.Element==0x0010)&&(!icon) {
					size:= uint32(cols)*uint32(rows)*uint32(bitsa)/8
					if RGB {
						size = 3*size
					}
					if frames>0 {
						size = uint32(frames)*size
					} else {
						frames=1
					}
					if(size==0){
						log.Println("ERROR, DcmObj::ConvertTransferSyntax, size=0")
						return false
					}
					
					if tag.Length==0xFFFFFFFF { // Compressed

					} else { // Uncompressed

					}
				}
			} 
		} 
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	return flag
}

func main() {
	media.InitDict()
	obj, err := media.NewDCMObjFromFile("test.dcm")
	if err != nil {
		log.Panic(err)
	}
	ConvertTS(obj, "")
}
