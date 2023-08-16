[![Build Status](https://drone.onebytedata.net/api/badges/one-byte-data/obd-dicom/status.svg?ref=refs/heads/main)](https://drone.onebytedata.net/one-byte-data/obd-dicom)

# obd-dicom

One Byte Data DICOM Golang Library

## Install

```bash
go get -u github.com/one-byte-data/obd-dicom
```

## Usage

### Load DICOM File

```golang
obj, err := media.NewDCMObjFromFile(fileName)
if err != nil {
  log.Panicln(err)
}
obj.DumpTags()
```

### Send C-Echo Request
```golang
scu := services.NewSCU(destination)
err := scu.EchoSCU(0)
if err != nil {
  log.Fatalln(err)
}
log.Println("CEcho was successful")
```

### Send C-Find Request
```golang
request := utils.DefaultCFindRequest()
scu := services.NewSCU(destination)
scu.SetOnCFindResult(func(result media.DcmObj) {
  log.Printf("Found study %s\n", result.GetString(tags.StudyInstanceUID))
  result.DumpTags()
})

count, status, err := scu.FindSCU(request, 0)
if err != nil {
  log.Fatalln(err)
}
```

### Send C-Store Request
```golang
scu := services.NewSCU(destination)
err := scu.StoreSCU(fileName, 0)
if err != nil {
  log.Fatalln(err)
}
```

### Send C-Move Request
```golang
request := utils.DefaultCMoveRequest(studyUID)

scu := services.NewSCU(destination)
_, err := scu.MoveSCU(destinationAE, request, 0)
if err != nil {
  log.Fatalln(err)
}
```

### Start SCP Server
```golang
scp := services.NewSCP(*port)

scp.OnAssociationRequest(func(request network.AAssociationRQ) bool {
  called := request.GetCalledAE()
  return *calledAE == called
})

scp.OnCFindRequest(func(request network.AAssociationRQ, queryLevel string, query media.DcmObj) ([]media.DcmObj, uint16) {
  query.DumpTags()
  results := make([]media.DcmObj, 0)
  for i := 0; i < 10; i++ {
    results = append(results, utils.GenerateCFindRequest())
  }
  return results, dicomstatus.Success
})

scp.OnCMoveRequest(func(request network.AAssociationRQ, moveLevel string, query media.DcmObj) uint16 {
  query.DumpTags()
  return dicomstatus.Success
})

scp.OnCStoreRequest(func(request network.AAssociationRQ, data media.DcmObj) uint16 {
  log.Printf("INFO, C-Store recieved %s", data.GetString(tags.SOPInstanceUID))
  directory := filepath.Join(*datastore, data.GetString(tags.PatientID), data.GetString(tags.StudyInstanceUID), data.GetString(tags.SeriesInstanceUID))
  os.MkdirAll(directory, 0755)

  path := filepath.Join(directory, data.GetString(tags.SOPInstanceUID)+".dcm")

  err := data.WriteToFile(path)
  if err != nil {
    log.Printf("ERROR: There was an error saving %s : %s", path, err.Error())
  }
  return dicomstatus.Success
})

err := scp.Start()
if err != nil {
  log.Fatal(err)
}
```