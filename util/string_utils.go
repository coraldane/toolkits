package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/coraldane/logger"
	"github.com/pborman/uuid"
	"reflect"
	"strings"
	"time"
)

var (
	DATE_FORMAT_DEFAULT    = "2006-01-02 15:04:05"
	DATE_FORMAT_SIMPLE     = "20060102150405"
	DATE_FORMAT_DATE_SHORT = "20060102"
)

func Md5Hex(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func FormatNow(layout string) string {
	return time.Now().Format(layout)
}

func FormatUTCTime(t time.Time) string {
	loc, _ := time.LoadLocation("UTC")
	return t.In(loc).Format(DATE_FORMAT_DEFAULT)
}

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

func GenerateUrlParams(paramMap map[string]string) string {
	var buffer bytes.Buffer
	var index int

	for key, value := range paramMap {
		if index > 0 {
			buffer.WriteString("&")
		}
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString(value)
		index = index + 1
	}
	return buffer.String()
}

func ToJsonString(v interface{}) string {
	bs, err := json.Marshal(v)
	if nil != err {
		logger.Error("marshal json error", err)
		return ""
	}
	return string(bs)
}

func GetUUID() string {
	text := uuid.New()
	return strings.Replace(text, "-", "", -1)
}
