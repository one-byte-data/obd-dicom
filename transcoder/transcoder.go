package transcoder

import (
	"C"

	"git.onebytedata.com/odb/go-dicom/jpeglib"
	"git.onebytedata.com/odb/go-dicom/openjpeg"
)

func TranscodeJ2kToJpeg8(j2kData []byte, width uint16, height uint16, samples uint16, mode int, outData *[]byte, outSize *int) error {
	rawData := make([]byte, width * height * 3)
	if err := openjpeg.J2Kdecode(j2kData, uint32(len(j2kData)), rawData); err != nil {
		return err
	}

	if err := jpeglib.EIJG8encode(rawData, width, height, samples, outData, outSize, mode); err != nil {
		return err
	}

	return nil
}