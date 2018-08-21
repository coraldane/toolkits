package security

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

var (
	DEFAULT_DES_ENCRYPT_KEY = "!#$LiuTian@Copyright2017"
)

func DesEncryptWithDefaultKey(origData string) (string, error) {
	return TripleDesEncrypt([]byte(origData), []byte(DEFAULT_DES_ENCRYPT_KEY))
}

func DesDecryptWithDefaultKey(crypted string) (string, error) {
	return TripleDesDecrypt([]byte(crypted), []byte(DEFAULT_DES_ENCRYPT_KEY))
}

func TripleDesDecrypt(crypted, key []byte) (string, error) {
	decryptText, err := hex.DecodeString(string(crypted))
	block, err := des.NewTripleDESCipher(key)
	if nil != err {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(decryptText))
	blockMode.CryptBlocks(origData, decryptText)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

func TripleDesEncrypt(origData, key []byte) (string, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return hex.EncodeToString(crypted), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
