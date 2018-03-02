package security

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
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
func (this *GoogleAuthenticator) VerifyCode(secret, code string) (bool, error) {
	// now time
	calculatedCode, err := this.GetCode(secret)
	if err != nil {
		return false, err
	}
	if calculatedCode == code {
		return true, nil
	}
	return false, nil
}

// GetCode Calculate the code, with given secret and point in time
func (this *GoogleAuthenticator) GetCode(secret string) (string, error) {
	// decode the key from the first argument
	inputNoSpaces := strings.Replace(secret, " ", "", -1)
	inputNoSpacesUpper := strings.ToUpper(inputNoSpaces)
	secretKey, err := base32.StdEncoding.DecodeString(inputNoSpacesUpper)
	if err != nil {
		return "", err
	}

	epochSeconds := time.Now().Unix()
	value := oneTimePassword(secretKey, toBytes(epochSeconds/30))
	return fmt.Sprintf("%d", value), nil
}

/**
* Get QR-Code URL for image, from google charts.
*
* @param string $name
* @param string $secret
* @param string $issuer
*
* @return string
 */
func (this *GoogleAuthenticator) GetQRCodeUrl(name, secret, issuer string) string {
	strUrl := fmt.Sprintf("otpauth://totp/%s--%s?secret=%s&issuer=%s", name, time.Now().Format("2006-01-02,15:04"), secret, issuer)
	return fmt.Sprintf(`http://s.jiathis.com/qrcode.php?url=%s`, url.QueryEscape(strUrl))
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

func arrayFlip(oldArr []string) map[string]int {
	newArr := make(map[string]int, len(oldArr))
	for key, value := range oldArr {
		newArr[value] = key
	}
	return newArr
}
