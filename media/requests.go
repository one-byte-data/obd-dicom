package media

import "git.onebytedata.com/OneByteDataPlatform/go-dicom/tags"

// DefaultCFindRequest - Creates a default C-Find request
func DefaultCFindRequest() DcmObj {
	query := NewEmptyDCMObj()
	query.WriteString(tags.StudyDate, "")
	query.WriteString(tags.StudyTime, "")
	query.WriteString(tags.AccessionNumber, "")
	query.WriteString(tags.QueryRetrieveLevel, "STUDY")
	query.WriteString(tags.ModalitiesInStudy, "")
	query.WriteString(tags.StudyDescription, "")
	query.WriteString(tags.PatientName, "")
	query.WriteString(tags.PatientID, "")
	query.WriteString(tags.PatientBirthDate, "")
	query.WriteString(tags.PatientSex, "")
	query.WriteString(tags.StudyInstanceUID, "")
	return query
}

// DefaultCMoveRequest - Creates a default C-Move request
func DefaultCMoveRequest(studyUID string) DcmObj {
	query := NewEmptyDCMObj()
	query.WriteString(tags.StudyDate, "")
	query.WriteString(tags.StudyTime, "")
	query.WriteString(tags.AccessionNumber, "")
	query.WriteString(tags.QueryRetrieveLevel, "STUDY")
	query.WriteString(tags.ModalitiesInStudy, "")
	query.WriteString(tags.StudyDescription, "")
	query.WriteString(tags.PatientName, "")
	query.WriteString(tags.PatientID, "")
	query.WriteString(tags.PatientBirthDate, "")
	query.WriteString(tags.PatientSex, "")
	query.WriteString(tags.StudyInstanceUID, studyUID)
	return query
}
