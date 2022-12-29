package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/spf13/viper"
)

/*
加密
u.CreatedAt = time.Now().Local().Format("2006-01-02 15:04:05")
u.Ip = c.ClientIP()

	if b, err = json.Marshal(u); err != nil {
	   ErrorPage(c, err)
	   return
	}

c.SetCookie("gin_cookie", encryption.EnString(string(b)), 0, "/", "", false, true)
c.Redirect(302, "/")
*/
func EnString(data string) (enText string) {
	var (
		origData []byte
	)
	origData, _ = AesEncrypt([]byte(data))
	enText, _ = EnPwdCode(origData)
	return enText
}

/*
解密

	if err = json.Unmarshal([]byte(encryption.DeString(s)), &u); err != nil {
	   return err
	}

c.SetCookie("gin_cookie", "", -1, "/", "", false, true)
c.Redirect(302, "/")
*/
func DeString(data string) (text string) {
	var (
		origData, deText []byte
		err              error
	)
	origData, _ = DePwdCode(data)

	if deText, err = AesDeCrypt(origData); err != nil {

		return err.Error()
	}
	return string(deText)
}

// AES PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充反向操作，删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	// 获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	return origData[:(length - int(origData[length-1]))], nil
}

// 实现加密
func AesEncrypt(origData []byte) (crypted []byte, err error) {
	var (
		blockSize int
		block     cipher.Block
	)
	// 创建加密算法实例
	if block, err = aes.NewCipher([]byte(viper.GetString("aes.PW_KEY"))); err != nil {
		return nil, err
	}
	// 获取块大小
	blockSize = block.BlockSize()
	// 对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	// 采用AES加密方法中的CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(viper.GetString("aes.PW_KEY"))[:blockSize])
	crypted = make([]byte, len(origData))
	// 执行加密
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 加密base64
func EnPwdCode(pwd []byte) (s string, err error) {
	var result []byte
	if result, err = AesEncrypt(pwd); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// 解密base64
func DePwdCode(pwd string) (b []byte, err error) {
	if b, err = base64.StdEncoding.DecodeString(pwd); err != nil {
		return nil, err
	}
	// 执行AES解密
	return AesDeCrypt(b)
}

// AES解密
func AesDeCrypt(cypted []byte) (b []byte, err error) {
	var (
		block     cipher.Block
		blockSize int
	)

	if block, err = aes.NewCipher([]byte(viper.GetString("aes.PW_KEY"))); err != nil {
		return nil, err
	}

	// 获取块大小
	blockSize = block.BlockSize()
	// 创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, []byte(viper.GetString("aes.PW_KEY"))[:blockSize])
	b = make([]byte, len(cypted))
	blockMode.CryptBlocks(b, cypted)

	return PKCS7UnPadding(b)
}
