package database

import "git.onebytedata.com/odb/go-dicom/media"

type SQLite struct {
}

func NewSQLiteDatabase() Database {
	return &SQLite{}
}

func (s *SQLite) AddPatient(dcmObj media.DcmObj) error {
	return nil
}

func (s *SQLite) AddStudy(dcmObj media.DcmObj) error {
	return nil
}

func (s *SQLite) AddSeries(dcmObj media.DcmObj) error {
	return nil
}

func (s *SQLite) AddInstance(dcmObj media.DcmObj) error {
	return nil
}
