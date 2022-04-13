package uuids

import (
	"hash/fnv"
	"strconv"

	"git.onebytedata.com/odb/go-dicom/imp"
)

func hash32(text string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

func CreateStudyUID(patName string, patID string, accNum string, stDate string) string {
	StudyUID := imp.GetImpClassUID()
	value := int(hash32(patName + patID + accNum + stDate))
	StudyUID = StudyUID + "." + strconv.Itoa(value) // 25 bytes + 11 bytes
	return StudyUID
}

func CreateSeriesUID(RootUID string, Modality string, SeriesNumber string) string {
	value := int(hash32(Modality + SeriesNumber))
	return (RootUID + "." + strconv.Itoa(value)) // 36 bytes + 11 bytes
}

func CreateInstanceUID(RootUID string, InstNumber string) string {
	return (RootUID + "." + InstNumber) // 47 bytes + 2 bytes
}
