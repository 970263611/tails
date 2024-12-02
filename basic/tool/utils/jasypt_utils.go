package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// Decrypt 解密 基于PBEWithMD5AndDES算法 Jasypt Java 包: org.jasypt.encryption.pbe.StandardPBEStringEncryptor 的 golang 实现
func Decrypt(msg, password string) (string, error) {
	msgBytes, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}
	salt := msgBytes[:8]
	encText := msgBytes[8:]

	dk, iv := getDerivedKey(password, salt, 1000)

	text, err := DesDecrypt(encText, dk, iv)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func DesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	padding, err := PKCS5UnPadding(origData)
	if err != nil {
		return nil, err
	}
	origData = padding
	return origData, nil
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if length-unPadding < 0 || length-unPadding > len(origData) {
		return nil, errors.New("解密异常,请确保是否使用了正确的jasypt.salt")
	}
	return origData[:(length - unPadding)], nil
}

// Encrypt 加密 基于PBEWithMD5AndDES算法 Jasypt Java 包: org.jasypt.encryption.pbe.StandardPBEStringEncryptor 的 golang 实现
func Encrypt(msg, password string) (string, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	padNum := 8 - (len(msg) % 8)
	for i := 0; i <= padNum; i++ {
		msg += string(rune(padNum))
	}
	dk, iv := getDerivedKey(password, salt, 1000)
	encText, err := DesEncrypt([]byte(msg), dk, iv)
	if err != nil {
		return "", err
	}
	r := append(salt, encText...)
	encodeString := base64.StdEncoding.EncodeToString(r)
	return encodeString, nil
}

func DesEncrypt(origData, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func pKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func getDerivedKey(password string, salt []byte, count int) ([]byte, []byte) {
	key := md5.Sum([]byte(password + string(salt)))
	for i := 0; i < count-1; i++ {
		key = md5.Sum(key[:])
	}
	return key[:8], key[8:]
}
