package encode

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/coraldane/logger"
	"time"
)

var (
	DATE_FORMAT_DEFAULT    = "2006-01-02 15:04:05"
	DATE_FORMAT_SIMPLE     = "20060102150405"
	DATE_FORMAT_DATE_SHORT = "20060102"
)

func FormatNow(layout string) string {
	return time.Now().Format(layout)
}

func FormatUTCTime(t time.Time) string {
	loc, _ := time.LoadLocation("UTC")
	return t.In(loc).Format(DATE_FORMAT_DEFAULT)
}

func ToJsonString(v interface{}) string {
	bs, err := json.Marshal(v)
	if nil != err {
		logger.Error("marshal json error", err)
		return ""
	}
	return string(bs)
}

func Md5Hex(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}


