# OBD Dicom

## Usage

  -calledae string
    	AE of the destination (default "DICOM_SCP")

  -callingae string
    	AE of the client (default "DICOM_SCU")

  -cecho
    	Send C-Echo to the destination

  -cfind
    	Send C-Find request to the destination

  -cmove
    	Send C-Move request to the destination

  -cstore
    	Sends a C-Store request to the destination

  -datastore string
    	Directory to use as SCP storage

  -destinationae string
    	AE of the destination for a C-Move request

  -dump
    	Dump contents of DICOM file to stdout

  -file string
    	DICOM file to be sent

  -host string
    	Destination host name or IP (default "localhost")

  -port int
    	Port of the destination system (default 1040)

  -query string
    	Comma seperated query to be sent with request ex: 00080020=test

  -scp
    	Start a SCP
      
  -studyuid string
    	Study UID to be added to request
