package main

import (
	"database/sql"
	"log"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	_ "github.com/mattn/go-sqlite3"

)

type DCMStudy struct {
	PatientID       string
	PatientName     string
	PatientBD       string
	PatientSex      string
	PatientComments string

	StudyInstanceUID   string
	StudyDate          string
	StudyTime          string
	Modality           string
	InstitutionName    string
	Description        string
	AccessionNumber    string
	ReferringPhysician string
}

func (study *DCMStudy) Query(obj media.DcmObj) string {
	var tag media.DcmTag
	var query, partial string
	previous := false

	for i := 0; i < len(obj.GetTags()); i++ {
		tag = obj.GetTag(i)
		if tag.Length > 0 {
			switch tag.Group {
			case 0x08:
				switch tag.Element {
				case 0x20:
					study.StudyDate = tag.GetString()
					partial = "StudyDate='" + study.StudyDate + "'"
					break
				case 0x30:
					study.StudyTime = tag.GetString()
					partial = "StudyTime='" + study.StudyTime + "'"
					break
				case 0x50:
					study.AccessionNumber = tag.GetString()
					partial = "AccessionNumber='" + study.AccessionNumber + "'"
					break
				case 0x60:
					study.Modality = tag.GetString()
					partial = "Modality='" + study.Modality + "'"
					break
				case 0x80:
					study.InstitutionName = tag.GetString()
					partial = "InstitutionName LIKE '%" + study.InstitutionName + "'%"
					break
				case 0x90:
					study.ReferringPhysician = tag.GetString()
					partial = "ReferringPhysician='" + study.ReferringPhysician + "'"
					break
				case 0x1030:
					study.Description = tag.GetString()
					partial = "Description LIKE '%" + study.Description + "%'"
					break
				}
				break
			case 0x10:
				switch tag.Element {
				case 0x0010:
					study.PatientName = tag.GetString()
					partial = "PatientName LIKE '" + study.PatientName + "%'"
					break
				case 0x0020:
					study.PatientID = tag.GetString()
					partial = "PatientID='" + study.PatientID + "'"
					break
				case 0x0030: //Patient Birth Date
					study.PatientBD = tag.GetString()
					partial = "PatientBD='" + study.PatientBD + "'"
					break
				case 0x0040:
					study.PatientSex = tag.GetString()
					partial = "PatientSex='" + study.PatientSex + "'"
					break
				}
				break
			case 0x20:
				switch tag.Element {
				case 0x000D:
					study.StudyInstanceUID = tag.GetString()
					partial = "StudyInstanceUID='" + study.StudyInstanceUID + "'"
					break
				}
				break
			}
			if len(partial) > 0 {
				if previous == true {
					query = query + " AND " + partial
				} else {
					query = " WHERE " + partial
					previous = true
				}
			}
		}
	}
	return query
}

func (study *DCMStudy) QueryResult(obj media.DcmObj) media.DcmObj {
	var added bool
	var tag media.DcmTag
	query := media.NewEmptyDCMObj()
	query.SetTransferSyntax("1.2.840.10008.1.2")

	for i := 0; i < len(obj.GetTags()); i++ {
		tag = obj.GetTag(i)
		added=true
		switch tag.Group {
		case 0x08:
			switch tag.Element {
			case 0x20:
				query.WriteStringGE(0x08, 0x20, "DA", study.StudyDate)
				break
			case 0x30:
				query.WriteStringGE(0x08, 0x30, "TM", study.StudyTime)
				break
			case 0x50:
				query.WriteStringGE(0x08, 0x50, "SH", study.AccessionNumber)
				break
			case 0x52:
				query.WriteStringGE(0x08, 0x52, "CS", "STUDY")
				break
			case 0x61:
				query.WriteStringGE(0x08, 0x61, "CS", study.Modality)
				break
			case 0x1030:
				query.WriteStringGE(0x08, 0x1030, "LO", study.Description)
				break
			default:
				added=false
			}
			break
		case 0x10:
			switch tag.Element {
			case 0x10:
				query.WriteStringGE(0x10, 0x10, "PN", study.PatientName)
				break
			case 0x20:
				query.WriteStringGE(0x10, 0x20, "LO", study.PatientID)
				break
			case 0x30:
				query.WriteStringGE(0x10, 0x30, "DA", study.PatientBD)
				break
			case 0x40:
				query.WriteStringGE(0x10, 0x40, "CS", study.PatientSex)
				break
			default:
				added=false
			}
			break
		default:
				added=false
		}
		if (tag.Group==0x20)&&(tag.Element==0x0D) {
			query.WriteStringGE(0x20, 0x0d, "UI", study.StudyInstanceUID)
			added=true
		}
		if added==false {
			query.Add(tag)
		}
	}
	return query
}

func (study *DCMStudy) Select(query media.DcmObj) (error, []media.DcmObj) {
	QueryString := study.Query(query)
	results := make([]media.DcmObj, 0)
	db, err := sql.Open("sqlite3", "./pacs.db")
	if err != nil {
		log.Println(err.Error())
		return err, nil
	}

	fields := "StudyDate, StudyTime, StudyDescription, AccessionNumber, ReferPhysician, StudyModality, PatientID, PatientName, PatientSex, PatientBD, StudyInstanceUID"
	QueryString = "SELECT " + fields + " FROM Study " + QueryString
	rows, err := db.Query(QueryString)
	if err != nil {
		log.Println(err.Error())
		return err, nil
	}

	for rows.Next() {
		rows.Scan(&study.StudyDate, &study.StudyTime, &study.Description, &study.AccessionNumber, &study.ReferringPhysician, &study.Modality, &study.PatientID, &study.PatientName, &study.PatientSex, &study.PatientBD, &study.StudyInstanceUID)
		obj := study.QueryResult(query)
		results = append(results, obj)
	}
	rows.Close()
	return nil, results
}
