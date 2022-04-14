package codingscheme

type CodingScheme struct {
	UID         string
	Name        string
	Description string
	Type        string
}

func GetCodingSchemeFromName(name string) *CodingScheme {
	for _, cs := range codingSchemes {
		if cs.Name == name {
			return cs
		}
	}
	return nil
}

func GetCodingSchemeFromUID(uid string) *CodingScheme {
	for _, cs := range codingSchemes {
		if cs.UID == uid {
			return cs
		}
	}
	return nil
}
