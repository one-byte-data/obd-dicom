package uid

type SOPClass struct {
	UID         string
	Name        string
	Description string
}

func GetSOPClassFromName(name string) *SOPClass {
	for _, sop := range SOPClasses {
		if sop.Name == name {
			return sop
		}
	}
	for _, sop := range TransferSyntaxes {
		if sop.Name == name {
			return sop
		}
	}
	return nil
}

func GetSOPClassFromUID(uid string) *SOPClass {
	for _, sop := range SOPClasses {
		if sop.UID == uid {
			return sop
		}
	}
	for _, sop := range TransferSyntaxes {
		if sop.UID == uid {
			return sop
		}
	}
	return nil
}

// VerificationSOPClass - (1.2.840.10008.1.1)
const VerificationSOPClass = "1.2.840.10008.1.1"

// PatientRootQueryRetrieveInformationModelFIND - (1.2.840.10008.5.1.4.1.2.1.1)
const PatientRootQueryRetrieveInformationModelFIND = "1.2.840.10008.5.1.4.1.2.1.1"

// PatientRootQueryRetrieveInformationModelMOVE - (1.2.840.10008.5.1.4.1.2.1.2)
const PatientRootQueryRetrieveInformationModelMOVE = "1.2.840.10008.5.1.4.1.2.1.2"

// PatientRootQueryRetrieveInformationModelGET - (1.2.840.10008.5.1.4.1.2.1.3)
const PatientRootQueryRetrieveInformationModelGET = "1.2.840.10008.5.1.4.1.2.1.3"

// StudyRootQueryRetrieveInformationModelFIND = "1.2.840.10008.5.1.4.1.2.2.1"
const StudyRootQueryRetrieveInformationModelFIND = "1.2.840.10008.5.1.4.1.2.2.1"

// StudyRootQueryRetrieveInformationModelMOVE = "1.2.840.10008.5.1.4.1.2.2.2"
const StudyRootQueryRetrieveInformationModelMOVE = "1.2.840.10008.5.1.4.1.2.2.2"

// StudyRootQueryRetrieveInformationModelGET = "1.2.840.10008.5.1.4.1.2.2.3"
const StudyRootQueryRetrieveInformationModelGET = "1.2.840.10008.5.1.4.1.2.2.3"

var SOPClasses = []*SOPClass{
	{
		UID:         "1.2.840.10008.1.1",
		Name:        "VerificationSOPClass",
		Description: "Verification SOP Class",
	},
	{
		UID:         "1.2.840.10008.5.1.4.1.2.1.1",
		Name:        "PatientRootQueryRetrieveInformationModelFIND",
		Description: "Patient Root Query Retrieve Information Model FIND",
	},
	{
		UID:         "1.2.840.10008.5.1.4.1.2.1.2",
		Name:        "PatientRootQueryRetrieveInformationModelMOVE",
		Description: "Patient Root Query Retrieve Information Model MOVE",
	},
	{
		UID:         "1.2.840.10008.5.1.4.1.2.1.3",
		Name:        "PatientRootQueryRetrieveInformationModelGET",
		Description: "Patient Root Query Retrieve Information Model GET",
	},
	{
		UID:         "1.2.840.10008.5.1.4.1.2.2.1",
		Name:        "StudyRootQueryRetrieveInformationModelFIND",
		Description: "Study Root Query Retrieve Information Model FIND",
	},
	{
		UID:         "1.2.840.10008.5.1.4.1.2.2.2",
		Name:        "StudyRootQueryRetrieveInformationModelMOVE",
		Description: "Study Root Query Retrieve Information Model MOVE",
	},
	{
		UID:         "1.2.840.10008.5.1.4.1.2.2.3",
		Name:        "StudyRootQueryRetrieveInformationModelGET",
		Description: "Study Root Query Retrieve Information Model GET",
	},
}
