package uid

var VerificationSOPClass = &SOPClass{
	UID:         "1.2.840.10008.1.1",
	Name:        "VerificationSOPClass",
	Description: "Verification SOP Class",
}

var PatientRootQueryRetrieveInformationModelFIND = &SOPClass{
	UID:         "1.2.840.10008.5.1.4.1.2.1.1",
	Name:        "PatientRootQueryRetrieveInformationModelFIND",
	Description: "Patient Root Query Retrieve Information Model FIND",
}

var PatientRootQueryRetrieveInformationModelMOVE = &SOPClass{
	UID:         "1.2.840.10008.5.1.4.1.2.1.2",
	Name:        "PatientRootQueryRetrieveInformationModelMOVE",
	Description: "Patient Root Query Retrieve Information Model MOVE",
}

var PatientRootQueryRetrieveInformationModelGET = &SOPClass{
	UID:         "1.2.840.10008.5.1.4.1.2.1.3",
	Name:        "PatientRootQueryRetrieveInformationModelGET",
	Description: "Patient Root Query Retrieve Information Model GET",
}

var StudyRootQueryRetrieveInformationModelFIND = &SOPClass{
	UID:         "1.2.840.10008.5.1.4.1.2.2.1",
	Name:        "StudyRootQueryRetrieveInformationModelFIND",
	Description: "Study Root Query Retrieve Information Model FIND",
}

var StudyRootQueryRetrieveInformationModelMOVE = &SOPClass{
	UID:         "1.2.840.10008.5.1.4.1.2.2.2",
	Name:        "StudyRootQueryRetrieveInformationModelMOVE",
	Description: "Study Root Query Retrieve Information Model MOVE",
}

var StudyRootQueryRetrieveInformationModelGET = &SOPClass{
	UID:         "1.2.840.10008.5.1.4.1.2.2.3",
	Name:        "StudyRootQueryRetrieveInformationModelGET",
	Description: "Study Root Query Retrieve Information Model GET",
}

var SOPClasses = []*SOPClass{
	VerificationSOPClass,
	PatientRootQueryRetrieveInformationModelFIND,
	PatientRootQueryRetrieveInformationModelMOVE,
	PatientRootQueryRetrieveInformationModelGET,
	StudyRootQueryRetrieveInformationModelFIND,
	StudyRootQueryRetrieveInformationModelMOVE,
	StudyRootQueryRetrieveInformationModelGET,
}
