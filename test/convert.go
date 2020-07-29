package main

import (
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

func SupportedTS(TransferSyntax string) bool {
	return true
}

func ConvertTS(obj media.DcmObj, outTS string) bool {
	if len(outTS) == 0 {
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

	var i int
	var tag media.DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if sq == 0 {
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	return true
}

/*
func main() {
	media.InitDict()
	obj, err := media.NewDCMObjFromFile("test.dcm")
	if err != nil {
		log.Panic(err)
	}
	ConvertTS(obj, "")
}
*/
