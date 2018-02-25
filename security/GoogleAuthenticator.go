package security

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrSecretLengthLss     = errors.New("secret length lss 6 error")
	ErrSecretLength        = errors.New("secret length error")
	ErrPaddingCharCount    = errors.New("padding char count error")
	ErrPaddingCharLocation = errors.New("padding char Location error")
	ErrParam               = errors.New("param error")

	Table = []string{
		"A", "B", "C", "D", "E", "F", "G", "H", // 7
		"I", "J", "K", "L", "M", "N", "O", "P", // 15
		"Q", "R", "S", "T", "U", "V", "W", "X", // 23
		"Y", "Z", "2", "3", "4", "5", "6", "7", // 31
		"=", // padding char
	}

	allowedValues = map[int]string{
		6: "======",
		4: "====",
		3: "===",
		1: "=",
		0: "",
	}
)

type GoogleAuthenticator struct {
	codeLen float64
	table   map[string]int
}

func NewGoogleAuthenticator() *GoogleAuthenticator {
	return &GoogleAuthenticator{
		codeLen: 6,
		table:   arrayFlip(Table),
	}
}

// SetCodeLength Set the code length, should be >=6
func (this *GoogleAuthenticator) SetCodeLength(length float64) error {
	if length < 6 {
		return ErrSecretLengthLss
	}
	this.codeLen = length
	return nil
}

// CreateSecret create new secret
// 16 characters, randomly chosen from the allowed base32 characters.
func (this *GoogleAuthenticator) CreateSecret(lens ...int) (string, error) {
	var (
		length int
		secret []string
	)
	// init length
	switch len(lens) {
	case 0:
		length = 16
	case 1:
		length = lens[0]
	default:
		return "", ErrParam
	}
	for i := 0; i < length; i++ {
		secret = append(secret, Table[rand.Intn(len(Table))])
	}
	return strings.Join(secret, ""), nil
}

// VerifyCode Check if the code is correct. This will accept codes starting from $discrepancy*30sec ago to $discrepancy*30sec from now
func (this *GoogleAuthenticator) VerifyCode(secret, code string, discrepancy int64) (bool, error) {
	// now time
	curTimeSlice := time.Now().Unix() / 30
	for i := -discrepancy; i <= discrepancy; i++ {
		calculatedCode, err := this.GetCode(secret, curTimeSlice+i)
		if err != nil {
			return false, err
		}
		if calculatedCode == code {
			return true, nil
		}
	}
	return false, nil
}

// GetCode Calculate the code, with given secret and point in time
func (this *GoogleAuthenticator) GetCode(secret string, timeSlices ...int64) (string, error) {
	var timeSlice int64
	switch len(timeSlices) {
	case 0:
		timeSlice = time.Now().Unix() / 30
	case 1:
		timeSlice = timeSlices[0]
	default:
		return "", ErrParam
	}
	secret = strings.ToUpper(secret)
	secretKey, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	tim, err := hex.DecodeString(fmt.Sprintf("%016x", timeSlice))
	if err != nil {
		return "", err
	}
	hm := HmacSha1(secretKey, tim)
	offset := hm[len(hm)-1] & 0x0F
	hashpart := hm[offset : offset+4]
	value, err := strconv.ParseInt(hex.EncodeToString(hashpart), 16, 0)
	if err != nil {
		return "", err
	}
	value = value & 0x7FFFFFFF
	modulo := int64(math.Pow(10, this.codeLen))
	format := fmt.Sprintf("%%0%dd", int(this.codeLen))
	return fmt.Sprintf(format, value%modulo), nil
}

/**
* Get QR-Code URL for image, from google charts.
*
* @param string $name
* @param string $secret
* @param array  $params
*
* @return string
 */
func (this *GoogleAuthenticator) GetQRCodeUrl(name, secret string, params ...int) string {
	var width, height int
	width, height = 150, 150
	switch len(params) {
	case 0:
	case 1:
		width = params[0]
	case 2:
		width = params[0]
		height = params[1]
	}

	strUrl := fmt.Sprintf("otpauth://totp/%s?secret=%s", name, secret)
	return fmt.Sprintf(`https://chart.googleapis.com/chart?chs=%dx%d&chld=H|0&cht=qr&chl=%s`, width, height, url.QueryEscape(strUrl))
}

func HmacSha1(key, data []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func arrayFlip(oldArr []string) map[string]int {
	newArr := make(map[string]int, len(oldArr))
	for key, value := range oldArr {
		newArr[value] = key
	}
	return newArr
}
