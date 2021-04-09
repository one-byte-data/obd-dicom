package main

import (
	"git.onebytedata.com/odb/go-dicom/tags"
	"log"
	"strconv"
	"strings"
)

func CheckDate(date string) bool {
	if len(date) != 8 {
		return false
	}

	temp := date[0:3]
	year, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		return false
	}
	if (year < 1900) || (year > 2100) {
		return false
	}
	temp = date[4:5]
	month, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		return false
	}
	if (month < 1) || (month > 12) {
		return false
	}
	temp = date[6:7]
	day, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		return false
	}
	if (day < 1) || (day > 31) {
		return false
	}
	return true
}

func CheckSex(Sex string) bool {
	if (Sex == "M") || (Sex == "F") || (Sex == "O") {
		return true
	}
	return false
}

func CheckNumeric(number string) bool {
	_, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return false
	}
	return true
}

func CheckModality(Modality string) bool {
	if (Modality == "CT") || (Modality == "MR") || (Modality == "US") || (Modality == "OT") || (Modality == "DR") ||
		(Modality == "DX") || (Modality == "RF") || (Modality == "NM") || (Modality == "XA") || (Modality == "CR") ||
		(Modality == "ES") || (Modality == "XC") || (Modality == "GM") || (Modality == "IO") || (Modality == "IVUS") ||
		(Modality == "MG") || (Modality == "PX") || (Modality == "PT") || (Modality == "RG") || (Modality == "SR") ||
		(Modality == "TCD") {
		return true
	}
	return false
}

func CheckOperator(Operator string) bool {
	if (Operator == "=") || (Operator == "CONTAINS") || (Operator == "<>") || (Operator == ">=") || (Operator == "<=") {
		return true
	}
	return false
}

func ValidateValue(Value string, group uint16, element uint16) bool {
	flag := true
	switch group {
	case 0x08:
		switch element {
		case 0x20:
			flag = CheckDate(Value)
			break
		case 0x30:
			flag = CheckNumeric(Value)
			break
		case 0x60:
			flag = CheckModality(Value)
			break
		}
		break
	case 0x10:
		switch element {
		case 0x30:
			flag = CheckDate(Value)
			break
		case 0x40:
			flag = CheckSex(Value)
			break
		}
		break
	case 0x20:
		switch element {
		case 0x11:
			flag = CheckNumeric(Value)
			break
		}
		break
	}
	return flag
}

func ValidateRule(Conditions string, Replacements string) bool {
	var flag bool
	var group, element uint16

	// First I verify that all conditions are met
	flag = true
	Cond := strings.Split(Conditions, "&")
	for i := 0; i < len(Cond) && flag; i++ {
		components := strings.Split(Cond[i], "|")
		if len(components) == 3 {
			group, element = tags.GetGroupElement(components[0])
			if ValidateValue(components[2], group, element) {
				flag = CheckOperator(components[1])
			} else {
				log.Println("ERROR, Condition[" + string(i) + "] fails ValidateValue")
				flag = false
			}
		} else {
			log.Println("ERROR, Condition[" + string(i) + "] not enough components")
			flag = false
		}
	}
	if flag {
		Rep := strings.Split(Replacements, "&")
		for i := 0; i < len(Rep); i++ {
			components := strings.Split(Rep[i], "|")
			if len(components) == 2 {
				group, element = tags.GetGroupElement(components[0])
				flag = ValidateValue(components[1], group, element)
				if flag == false {
					log.Println("ERROR, Replacement[" + string(i) + "] fails ValidateValue")
					break
				}
			}
		}
	}
	return flag
}
