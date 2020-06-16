package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/uuids"
)

// GetFolderFiles get files in a folder
func GetFolderFiles() []os.FileInfo {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(pathS)
	if err != nil {
		return nil
	}
	return files
}

// ConvertPDF convert a DICOM to PDF file
func ConvertPDF(DCMFile string, PDFFile string) bool {
	flag := false
	var obj media.DcmObj

	if obj.Read(DCMFile) {
		var study media.DCMStudy
		var pdfobj media.DcmObj

		study.GetInfo(obj)
		pdfobj.ExplicitVR = true
		pdfobj.BigEndian = false
		pdfobj.TransferSyntax = "1.2.840.10008.1.2.1"
		RootUID := uuids.CreateStudyUID(study.PatientName, study.PatientID, study.AccessionNumber, study.StudyDate)
		SeriesUID := uuids.CreateSeriesUID(RootUID, study.Modality, "300")
		InstanceUID := uuids.CreateInstanceUID(RootUID, "1")
		pdfobj.CreatePDF(study, SeriesUID, InstanceUID, PDFFile)
		pdfobj.Write("samplepdf.dcm")

		flag = true
	}
	return flag
}

func main() {
	PDFFile := ""
	DCMFile := ""
	files := GetFolderFiles()
	for i := 0; i < len(files); i++ {
		if filepath.Ext(files[i].Name()) == ".pdf" {
			PDFFile = files[i].Name()
		}
		if filepath.Ext(files[i].Name()) == ".dcm" {
			DCMFile = files[i].Name()
		}
	}
	if (len(PDFFile) > 0) && (len(DCMFile) > 0) {
		ConvertPDF(DCMFile, PDFFile)
	}
}
