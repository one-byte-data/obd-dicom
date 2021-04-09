package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/uuids"
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
	obj, err := media.NewDCMObjFromFile(DCMFile)
	if err != nil {
		return false
	}

	var study media.DCMStudy

	pdfobj := media.NewEmptyDCMObj()

	study.GetInfo(obj)
	pdfobj.SetExplicitVR(true)
	pdfobj.SetBigEndian(false)
	pdfobj.SetTransferSyntax("1.2.840.10008.1.2.1")
	RootUID := uuids.CreateStudyUID(study.PatientName, study.PatientID, study.AccessionNumber, study.StudyDate)
	SeriesUID := uuids.CreateSeriesUID(RootUID, study.Modality, "300")
	InstanceUID := uuids.CreateInstanceUID(RootUID, "1")
	pdfobj.CreatePDF(study, SeriesUID, InstanceUID, PDFFile)
	pdfobj.WriteToFile("samplepdf.dcm")

	return true
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
