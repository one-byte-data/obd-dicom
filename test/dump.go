package main

import (
	"fmt"
	"strconv"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

func DumpDcm(obj media.DcmObj) {
	sq := 0

	for i := 0; i < obj.TagCount(); i++ {
		tag := obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if sq == 0 {
			if tag.Length > 0 {
				var val string
				if (tag.VR == "SL") || (tag.VR == "SS") || (tag.VR == "US") {
					val = strconv.Itoa(int(tag.GetUShort()))
				} else if tag.VR == "UL" {
					val = strconv.Itoa(int(tag.GetUInt()))
				} else if tag.VR == "SQ" {
					seq := tag.ReadSeq(obj.IsExplicitVR())
					DumpDcm(seq)
				} else {
					if tag.Length < 256 {
						val = tag.GetString()
					} else {
						val = ""
					}
				}
				fmt.Printf("(%04x,%04x),%d, %s, %s, %s\n", tag.Group, tag.Element, tag.Length, tag.VR, val, media.TagDescription(tag.Group, tag.Element))
			} else {
				fmt.Printf("(%04x,%04x),%d, %s, %s, %s\n", tag.Group, tag.Element, tag.Length, tag.VR, "", media.TagDescription(tag.Group, tag.Element))
			}
		} else {
			fmt.Printf("--- (%04x,%04x),%d, %s, %s, %s\n", tag.Group, tag.Element, tag.Length, tag.VR, "", media.TagDescription(tag.Group, tag.Element))
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
}

/*
func main() {
	media.InitDict()
	obj, err := media.NewDCMObjFromFile("test.dcm")
	if err != nil {
		return
	}
	DumpDcm(obj)
}
*/