package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// AESEncrypt AES加密
func AESEncrypt(encodeStr string, key []byte) (decodeStr string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("加解密Key不正确/匹配")
		}
	}()
	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encodeBytes = PKCS5Padding(encodeBytes, blockSize)

	// 对IV有随机性要求，但没有保密性要求，所以常见的做法是将IV包含在加密文本当中
	ciphertext := make([]byte, blockSize+len(encodeBytes))
	//随机一个block大小作为IV
	//采用不同的IV时相同的秘钥将会产生不同的密文，可以理解为一次加密的session
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	//crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(ciphertext[blockSize:], encodeBytes)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//填充
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

// AESDecrypt AES解密
func AESDecrypt(decodeStr string, key []byte) (origin []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("加解密Key不正确/匹配")
		}
	}()
	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if len(decodeBytes) == 0 {
		return nil, errors.New(fmt.Sprintf("base64解码失败；%s", decodeStr))
	}

	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(decodeBytes) < blockSize {
		return nil, errors.New(fmt.Sprintf("base64解码失败；%s", decodeStr))
	}
	iv := decodeBytes[:blockSize]
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes)-blockSize)

	blockMode.CryptBlocks(origData, decodeBytes[blockSize:])
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
