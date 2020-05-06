package media

// DCMStudy study information structure
type DCMStudy struct {
	PatientID          string
	PatientName        string
	PatientBD          string
	PatientSex         string
	ReferringPhysician string
	ReportDate         string
	ReportTime         string
	AccessionNumber    string
	Modality           string
	InstitutionName    string
	Description        string
	StudyInstanceUID   string
	ReportText         string
	ObserverName       string
}
