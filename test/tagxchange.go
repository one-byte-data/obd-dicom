package main

import (
	"log"
	"strconv"
	"strings"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/tags"
)

func GetTag(obj media.DcmObj, group uint16, element uint16) media.DcmTag {
	var i int
	var tag media.DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if (sq == 0) && (tag.Length > 0) && (tag.Length != 0xFFFFFFFF) {
			if (tag.Group == group) && (tag.Element == element) {
				break
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
	return tag
}

func Insert(obj media.DcmObj, intag media.DcmTag) {
	var i int
	var tag media.DcmTag
	sq := 0
	for i = 0; i < obj.TagCount(); i++ {
		tag = obj.GetTag(i)
		if ((tag.VR == "SQ") && (tag.Length == 0xFFFFFFFF)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE000) && (tag.Length == 0xFFFFFFFF)) {
			sq++
		}
		if (sq == 0) && (tag.Length > 0) && (tag.Length != 0xFFFFFFFF) {
			if (tag.Group == intag.Group) && (tag.Element == intag.Element) {
				obj.SetTag(i, intag)
				return		
			}
			if (tag.Group==intag.Group)&&(tag.Element>intag.Element) {
				obj.InsertTag(i, intag)
				return
			}
		   if tag.Group>intag.Group {
				obj.InsertTag(i, intag)
				return
			}
		}
		if ((tag.Group == 0xFFFE) && (tag.Element == 0xE00D)) || ((tag.Group == 0xFFFE) && (tag.Element == 0xE0DD)) {
			sq--
		}
	}
}

func ApplyCondition(DICOMValue string, Condition string, Value string) bool {
	var DFUp, ValUp string

	DFUp = strings.ToUpper(DICOMValue)
	if strings.Index(Value, "\"") != -1 {
		ValUp = strings.Split(strings.ToUpper(Value), "\"")[1]
	} else {
		ValUp = strings.ToUpper(Value)
	}

	// Condition Equal to.
	if Condition == "=" {
		if strings.TrimSpace(ValUp) == "*" {
			return true
		}
		if (len(DFUp) == 0) && (strings.TrimSpace(ValUp) == "NULL") {
			return true
		}
		return (DFUp == ValUp)
	}

	// Condition Contains.
	if Condition == "CONTAINS" {
		return strings.Index(DFUp, ValUp) != -1
	}

	// Condition Greater than
	if Condition == ">=" {
		field, err := strconv.ParseInt(DFUp, 10, 32)
		if err==nil {
			val, err := strconv.ParseInt(ValUp, 10, 32)
			if err==nil {
				return (field >= val)
			} else {
				return false
			}
		} else {
			return false
		}
	}

	// Condition Less than
	if Condition == "<=" {
		field, err := strconv.ParseInt(DFUp, 10, 32)
		if err == nil {
			val, err := strconv.ParseInt(ValUp, 10, 32)
			if err== nil {
				return (field <= val)
			} else {
				return false
			}
		} else {
			return false
		}
	}

	// Condition Different than
	if Condition == "<>" {
		return (DFUp != ValUp)
	}

	return (false)
}

func CopyDCM(inobj media.DcmObj) media.DcmObj{
	outobj:= media.NewEmptyDCMObj()
	outobj.SetExplicitVR(inobj.IsExplicitVR())
	outobj.SetBigEndian(inobj.IsBigEndian())
	outobj.SetTransferSyntax(inobj.GetTransferSyntax())
	var tag media.DcmTag
	for i:=0; i< inobj.TagCount(); i++ {
		tag = inobj.GetTag(i)
		outobj.Add(tag)  
	}
	return outobj
}

func MultipleReplace(inobj media.DcmObj, Conditions string, Replacements string) media.DcmObj{
	// Mutiple Rules have this syntax:
	// if Cond1&Cond2&Cond3 then apply Rep1&Rep2
	// There can be n Conditions and m Replacements.
	var tag media.DcmTag
	var group, element, out_group, out_element uint16
	var DICOMTag, DICOMValue string

	outobj := CopyDCM(inobj)

	// First I verify that all conditions are met
	flag := true
	Cond := strings.Split(Conditions, "&")
	for i:=0; i<len(Cond) && flag; i++ {
		components:=strings.Split(Cond[i], "|")
		if len(components)==3 {
			DICOMTag = components[0]
			group, element = tags.GetGroupElement(DICOMTag)
			DICOMValue = inobj.GetStringGE(group, element)
			tag = GetTag(inobj, group, element)
			if tag.Group!=0 {
				if !ApplyCondition(DICOMValue, components[1], components[2]) {
					flag = false
				}
			} else {
				flag = false
			}
		} else {
			flag=false
		}
	}

	// If flag is still true do the replacements
	if flag {
		Rep := strings.Split(Replacements, "&")
		for i:=0; i<len(Rep); i++ {
			components:=strings.Split(Rep[i], "|")
			DICOMTag = components[0]
			WithValue := components[1]
			// Si existe lo modifico, si no existe
			// Tengo que crearlo y agregarlo en un lugar adecuado.
			out_group, out_element = tags.GetGroupElement(DICOMTag)
			// Si group y element == 0 ERROR!!
			if (out_group == 0) || (out_element == 0) {
				log.Println("ERROR, Tag Name not found: " + DICOMTag)
				break
			}
			tag = GetTag(inobj, out_group, out_element)
			// Si es un tag...
			if strings.Index(WithValue, "\"") == -1 {
				group, element = tags.GetGroupElement(DICOMTag)
				WithValue = inobj.GetStringGE(group, element)
			} else {
				WithValue = strings.Split(WithValue, "\"")[1]
			}
			length := len(WithValue)
			if length%2 != 0 {
				WithValue = WithValue + " "
				length++
			}
			if tag.Length != 0 {
				tag.Data = nil
			}
			tag.Length = uint32(length)
			tag.Data = make([]byte, tag.Length)
			copy(tag.Data, WithValue)
			if tag.Group == 0 {
				//insert...
				tag.Group = out_group
				tag.Element = out_element
			}
			Insert(outobj, tag)
		}
	}
	return outobj
}

func main() {
	media.InitDict()
	obj, err := media.NewDCMObjFromFile("images/rle_gray.dcm")
	if err != nil {
		log.Panic(err)
	}
	Conditions:="PatientName|CONTAINS|xyz"
	Replacements:="AccessionNumber|\"CONTA\""
	if ValidateRule(Conditions, Replacements)==true {
		out:=MultipleReplace(obj, Conditions, Replacements)
		out.WriteToFile("out.dcm")
	} else {
		log.Println("ERROR, Failed Rule Validation")
	}
}
