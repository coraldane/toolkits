package container

import "reflect"

func HasElem(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func RemoveElem(arr []string, index int) []string {
	retValue := make([]string, 0)
	for i, val := range arr {
		if i == index {
			continue
		}
		retValue = append(retValue, val)
	}
	return retValue
}

func MergeArray(arr []string, elem string) []string {
	newCards := make([]string, 0)
	newCards = append(newCards, arr...)
	if "" != elem {
		newCards = append(newCards, elem)
	}
	return newCards
}

func IndexOf(s interface{}, elem interface{}) int {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return i
			}
		}
	}
	return -1
}

func LastIndexOf(s interface{}, elem interface{}) int {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := arrV.Len() - 1; i >= 0; i-- {
			if arrV.Index(i).Interface() == elem {
				return i
			}
		}
	}
	return -1
}
