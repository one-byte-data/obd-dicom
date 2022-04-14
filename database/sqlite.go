package database

import (
	"database/sql"

	"git.onebytedata.com/odb/go-dicom/media"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
}

func NewSQLiteDatabase(dbFileName string) Database {
	_, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		panic(err)
	}

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
