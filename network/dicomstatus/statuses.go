package dicomstatus

// Success - 0x0000
const Success uint16 = 0x0000

// Cancel - 0xfe00
const Cancel uint16 = 0xfe00

// Pending - 0xff00
const Pending uint16 = 0xff00

// PendingWithWarnings - 0xff01
const PendingWithWarnings uint16 = 0xff01

// Warning - 0x0001
const Warning uint16 = 0x0001

// Failure - 0xa700
const FailureOutOfResources uint16 = 0xa700

// Failure - 0xa900
const FailureDoesNotMatchSOPClass uint16 = 0xa900

// Failure - 0x0122
const FailureSOPClassNotSupported uint16 = 0x0122

// Failure - 0xc000
const FailureUnableToProcess uint16 = 0xc000
