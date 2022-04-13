package imp

type Implementation interface {
	GetClassUID() string
	GetVersion() string
}

type implementation struct {
	classUID string
	version  string
}

var imp Implementation

func SetDefaultImplementation() {
	imp = &implementation{
		classUID: "1.2.826.0.1.3680043.10.90.999",
		version:  "One-Byte-Data",
	}
}

func SetImplementation(classUID string, version string) {
	imp = &implementation{
		classUID: classUID,
		version:  version,
	}
}

func GetImpClassUID() string {
	if imp == nil {
		SetDefaultImplementation()
	}
	return imp.GetClassUID()
}

func GetImpVersion() string {
	if imp == nil {
		SetDefaultImplementation()
	}
	return imp.GetVersion()
}

func (i *implementation) GetClassUID() string {
	if i.classUID == "" {
		i.setDefault()
	}
	return i.classUID
}

func (i *implementation) GetVersion() string {
	if i.classUID == "" {
		i.setDefault()
	}
	return i.version
}

func (i *implementation) setDefault() {
	SetDefaultImplementation()
}
