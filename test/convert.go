package main

import (
	"log"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"strconv"
)

func SupportedTS(TransferSyntax string) bool {
	return true
}

func Decomp(obj media.DcmObj, i int, img []byte, size uint32, frames uint32, bitsa uint16) bool {
	var tag media.DcmTag
	var j, offset, single uint32
	
	single = size/frames
	// DE-Compression
	obj.DelTag(i+1); // Delete offset table.
	if obj.GetTransferSyntax()=="1.2.840.10008.1.2.5" {
		for j=0; j<frames; j++ {
			offset = j*single
			tag = obj.GetTag(i+1)
			RLEdecode(tag.Data, img[offset:], tag.Length, single, bitsa)
			obj.DelTag(i+1)
			}
		obj.DelTag(i+1);
	} else if (obj.GetTransferSyntax()=="1.2.840.10008.1.2.4.70")||(obj.GetTransferSyntax()=="1.2.840.10008.1.2.4.57") {
		for j=0; j<frames; j++ {
			offset = j*single
			tag = obj.GetTag(i+1)
			if bitsa==8 {
				DIJG8decode(tag.Data, img[offset:], tag.Length)
			} else {
				DIJG16decode(tag.Data, img[offset:], tag.Length)
			}
		obj.DelTag(i+1)
		}
	obj.DelTag(i+1)
	} else if obj.GetTransferSyntax()=="1.2.840.10008.1.2.4.50" {
		for j=0; j<frames; j++ {
			offset = j*single
			tag = obj.GetTag(i+1)
			if bitsa==8 {
				DIJG8decode(tag.Data, img[offset:], tag.Length)
			} else {
				DIJG12decode(tag.Data, img[offset:], tag.Length)
			}
			obj.DelTag(i+1)
		}
		obj.DelTag(i+1);
	} else if obj.GetTransferSyntax()=="1.2.840.10008.1.2.4.51" {
		 for j=0; j<frames; j++ {
			offset = j*single
			tag = obj.GetTag(i+1)
			DIJG12decode(tag.Data, img[offset:], tag.Length)
			obj.DelTag(i+1)
		}
		obj.DelTag(i+1)
	} else if obj.GetTransferSyntax()=="1.2.840.10008.1.2.4.90" {
		for j=0; j< frames; j++ {
			offset = j*single
			tag = obj.GetTag(i+1)
			J2Kdecode(tag.Data, img[offset:], tag.Length)
			obj.DelTag(i+1)
		}
		obj.DelTag(i+1)
	} else if obj.GetTransferSyntax()=="1.2.840.10008.1.2.4.91" {
		for j=0; j<frames; j++ {
			offset = j*single
			tag = obj.GetTag(i+1)
			J2Kdecode(tag.Data, img[offset:], tag.Length)
			obj.DelTag(i+1)
		}
	obj.DelTag(i+1)
	}
}
	
func Comp(obj media.DcmObj, i *int, img []byte, RGB bool, cols uint16, rows uint16, bitss uint16, bitsa uint16, pixelrep uint16, planar uint16, frames uint32, outTS string) bool {
	var tag media.DcmTag
	var offset, size, jpeg_size, j uint32
	var JPEGData []byte
	var JPEGBytes, index int
		
	single := uint32(cols)*uint32(rows)*uint32(bitsa)/8
	size = single*frames
	if RGB {
		size = 3*size
	}
	
	index = *i
	tag = obj.GetTag(index)
	if outTS=="1.2.840.10008.1.2.4.70" {
		tag.VR="OB"
		tag.Length=0xFFFFFFFF
		if tag.Data!= nil {
			tag.Data=nil
		}
		index++
		newtag := media.DcmTag {0xFFFE, 0xE000, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		for j=0; j<frames; j++ {
			index++
			offset = j*uint32(cols)*uint32(rows)*uint32(bitsa)/8
			if RGB {
				offset=3*offset
			}
			// Bug de MRI Neusoft. 08/05/2015
			if bitsa==8 {
				if RGB {
					encode8(cols, rows, 3, img[offset:], JPEGData, &JPEGBytes, 4)
				} else {
					encode8(cols, rows, 1, img[offset:], JPEGData, &JPEGBytes, 4)
				}
			} else {
				encode16(cols, rows, 1, img[offset/2:], JPEGData, &JPEGBytes)
			}
			newtag = media.DcmTag {0xFFFE, 0xE000, uint32(JPEGBytes), "DL", JPEGData, obj.IsBigEndian()}
			obj.SetTag(index, newtag)
			JPEGData=nil
		}
		index++
		newtag = media.DcmTag {0xFFFE, 0xE0DD, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		*i=index
	} else if outTS=="1.2.840.10008.1.2.4.50" {
		tag.VR="OB"
		tag.Length=0xFFFFFFFF
		if tag.Data!= nil {
			tag.Data=nil
		}
		index++
		newtag := media.DcmTag {0xFFFE, 0xE000, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		jpeg_size=0
		for j=0; j<frames; j++ {
			index++;
			offset = j*uint32(cols)*uint32(rows)*uint32(bitsa)/8
			if(RGB) {
				offset=3*offset
				encode8(cols, rows, 3, img[offset:], JPEGData, &JPEGBytes, 0)
			} else {
				if(bitsa==8) {
					encode8(cols, rows, 1, img[offset:], JPEGData, &JPEGBytes, 0)
				} else { // ERROR...
					// Can't use this transfer Syntax with bitsa!=8
					return false
				}
			}
			newtag = media.DcmTag {0xFFFE, 0xE000, uint32(JPEGBytes), "DL", JPEGData, obj.IsBigEndian()}
			obj.SetTag(index, newtag)
			JPEGData=nil
			jpeg_size = jpeg_size + uint32(JPEGBytes)
		}
		index++
		newtag = media.DcmTag {0xFFFE, 0xE0DD, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		*i=index
	} else if outTS=="1.2.840.10008.1.2.4.51" {
		if (bitss==8)&&(bitsa!=16) {
			return false
			}
		tag.VR="OB"
		tag.Length=0xFFFFFFFF
		if tag.Data!= nil {
			tag.Data=nil
			}
		index++
		newtag := media.DcmTag {0xFFFE, 0xE000, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		jpeg_size=0;
		for j=0; j<frames; j++ {
			index++
			offset = j*uint32(cols)*uint32(rows)*uint32(bitsa)/8
			if(bitss>12) {
				return false
			}
			encode12(cols, rows, 1, img[offset/2:], JPEGData, &JPEGBytes)
			newtag = media.DcmTag {0xFFFE, 0xE000, uint32(JPEGBytes), "DL", JPEGData, obj.IsBigEndian()}
			obj.SetTag(index, newtag)
			JPEGData=nil
			jpeg_size = jpeg_size + uint32(JPEGBytes)
		}
		index++
		newtag = media.DcmTag {0xFFFE, 0xE0DD, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		*i=index
	} else if outTS=="1.2.840.10008.1.2.4.90" {
		tag.VR="OB"
		tag.Length=0xFFFFFFFF
		if tag.Data!= nil {
			tag.Data=nil
			}
		index++
		newtag := media.DcmTag {0xFFFE, 0xE000, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		for j=0; j<frames; j++ {
			index++
			offset = j*uint32(cols)*uint32(rows)*uint32(bitsa)/8
			if(RGB) {
				offset=3*offset
				J2Kencode(img[offset:], JPEGData, &JPEGBytes, cols, rows, 3, bitsa, 0)
			} else {
				J2Kencode(img[offset:], JPEGData, &JPEGBytes, cols, rows, 1, bitsa, 0)
			}
			newtag = media.DcmTag {0xFFFE, 0xE000, uint32(JPEGBytes), "DL", JPEGData, obj.IsBigEndian()}
			obj.SetTag(index, newtag)
			JPEGData=nil
			}
		index++
		newtag = media.DcmTag {0xFFFE, 0xE0DD, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		*i=index
	} else if outTS=="1.2.840.10008.1.2.4.91" {
		tag.VR="OB"
		tag.Length=0xFFFFFFFF
		if tag.Data!= nil {
			tag.Data=nil
			}
		index++
		newtag := media.DcmTag {0xFFFE, 0xE000, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		jpeg_size=0;
		for j=0; j<frames; j++ {
			index++
			offset = j*uint32(cols)*uint32(rows)*uint32(bitsa)/8
			if RGB {
				offset=3*offset;
				J2Kencode(img[offset:], JPEGData, &JPEGBytes, cols, rows, 3, bitsa, 10)
			} else {
				J2Kencode(img[offset:], JPEGData, &JPEGBytes, cols, rows, 1, bitsa, 10);
			}
			newtag = media.DcmTag {0xFFFE, 0xE000, uint32(JPEGBytes), "DL", JPEGData, obj.IsBigEndian()}
			obj.SetTag(index, newtag)
			JPEGData=nil
			jpeg_size = jpeg_size + uint32(JPEGBytes)
		}
		index++
		newtag = media.DcmTag {0xFFFE, 0xE0DD, 0, "DL", nil, obj.IsBigEndian()}
		obj.SetTag(index, newtag)
		*i=index
	} else {
		if bitss==8 {
			tag.VR="OB"
		} else {
			tag.VR="OW"
		}
		tag.Length=size
		if tag.Data!=nil {
			tag.Data=nil
		}
		tag.Data = make([]byte, tag.Length)
		copy(tag.Data, img)
	}
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
	frames:=uint32(0)
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
					img := make([]byte, size)
					if tag.Length==0xFFFFFFFF { // 
						if Decomp(i, img, size, frames, bitsa) {
							flag = Comp(i, img, RGB, cols, rows, bitss, bitsa, pixelrep, planar, frames, outTS)
                        }
 					} else { // Uncompressed
						if RGB&&(planar==1) { // change from planar=1 to planar=0
							var img_offset, img_size uint32
							img_size = size/frames
							for f:=uint32(0); f<frames; f++ {
								img_offset = img_size*f
								 for j:=uint32(0); j<img_size/3; j++ {
									  img[3*j+img_offset] = tag.Data[j+img_offset]
									  img[3*j+1+img_offset] = tag.Data[j+img_size/3+img_offset]
				 					  img[3*j+2+img_offset] = tag.Data[j+2*img_size/3+img_offset]
									  }
								 }
							 planar=0
						} else {
							copy(img, tag.Data)
						}
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
