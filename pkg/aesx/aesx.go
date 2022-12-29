// Package aesx
// AES CBC模式加密
package aesx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type AesEncrypt struct {
	key   []byte //密钥
	iv    []byte // IV偏移量，AES数据块长度为128位，所以IV长度需要为16个字符，超过长度则截取，不足则在末尾填充'0'补足
	block cipher.Block
}

func NewAesEncrypt(key, iv string) (*AesEncrypt, error) {
	ret := &AesEncrypt{
		key:   []byte(key),
		iv:    []byte(iv),
		block: nil,
	}
	var err error
	ret.block, err = aes.NewCipher(ret.key)
	return ret, err
}

// Encrypt 加密
func (a *AesEncrypt) Encrypt(in string) ([]byte, error) {
	origData := []byte(in)
	origData = PKCS5Padding(origData, a.block.BlockSize())
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	iv := make([]byte, a.block.BlockSize())
	copy(iv, a.iv)
	bm := cipher.NewCBCEncrypter(a.block, iv)
	bm.CryptBlocks(crypted, origData)
	//var b = base64.StdEncoding.EncodeToString(crypted)
	return crypted, nil
}

// Decrypt 解密
func (a *AesEncrypt) Decrypt(crypted []byte) (string, error) {

	origData := make([]byte, len(crypted))
	iv := make([]byte, a.block.BlockSize())
	copy(iv, a.iv)
	bm := cipher.NewCBCDecrypter(a.block, iv)
	bm.CryptBlocks(origData, crypted)
	origData, err := PKCS5UnPadding(origData)
	var out = string(origData)
	return out, err
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length < 1 {
		return nil, fmt.Errorf("invalid unpadding length")
	}
	// 去掉最后一个字节 unPadding 次
	unPadding := int(origData[length-1])

	if length < unPadding {
		return nil, fmt.Errorf("invalid unpadding length")
	}
	return origData[:(length - unPadding)], nil
}

// EncryptBase64 加密返回base64编码
func (a *AesEncrypt) EncryptBase64(original string) (string, error) {
	bytes, err := a.Encrypt(original)
	if err != nil {
		return "", err
	}
	toString := base64.URLEncoding.EncodeToString(bytes)
	return toString, nil
}

// DecryptBase64 解密base64编码结果
func (a *AesEncrypt) DecryptBase64(cipherText string) (string, error) {
	decodeString, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	original, err := a.Decrypt(decodeString)
	if err != nil {
		return "", err
	}
	return original, nil
}
