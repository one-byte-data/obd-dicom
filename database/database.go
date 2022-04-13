package database

import "git.onebytedata.com/odb/go-dicom/media"

type Database interface {
	AddPatient(dcmObj media.DcmObj) error
	AddStudy(dcmObj media.DcmObj) error
	AddSeries(dcmObj media.DcmObj) error
	AddInstance(dcmObj media.DcmObj) error
}
