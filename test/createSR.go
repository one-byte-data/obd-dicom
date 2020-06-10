package main

import (
	"fmt"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"strconv"
)

func createSR() {
	sq := 0
	media.InitDict()
	var obj media.DcmObj
	obj.Read("test.dcm")

	for i := 0; i < obj.TagCount(); i++ {
		tag := obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if (sq == 0) && (tag.Length > 0) && (tag.Length != 0xFFFFFFFF) {
			if tag.Length < 256 {
				var val string
				if (tag.VR == "SL") || (tag.VR == "SS") || (tag.VR == "US") {
					val = strconv.Itoa(int(tag.GetUShort()))
				} else if tag.VR == "UL" {
					val = strconv.Itoa(int(tag.GetUInt()))
				} else {
					val = tag.GetString()
				}
				fmt.Printf("(%04x,%04x),%d, %s, %s, %s\n", tag.Group, tag.Element, tag.Length, tag.VR, val, media.TagDescription(tag.Group, tag.Element))
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	//	obj.Write("output.dcm")
	var srobj media.DcmObj
	var study media.DCMStudy

	srobj.ExplicitVR = true
	srobj.BigEndian = false
	srobj.TransferSyntax = "1.2.840.10008.1.2.1"

	study.AccessionNumber = "123456"
	study.Description = "Complete Thorax"
	study.InstitutionName = "Central Hospital"
	study.Modality = "CR"
	study.ObserverName = "Senior Radiologist"
	study.PatientBD = "20000101"
	study.PatientID = "99999"
	study.PatientName = "Jose Perez"
	study.PatientSex = "M"
	study.ReferringPhysician = "Asking Forstudies"
	study.ReportText = "This is a normal study, nothing to report."
	study.StudyInstanceUID = "9999.9999.1"
	srobj.CreateSR(study, "8888.8888.1", "7777.7777.1")
	srobj.Write("samplesr.dcm")
}
