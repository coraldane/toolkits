package encode

import (
	"strings"
	"github.com/pborman/uuid"
)

func GetUUID() string {
	text := uuid.New()
	return strings.Replace(text, "-", "", -1)
}