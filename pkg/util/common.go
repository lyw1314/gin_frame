package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var ZeroFloat64 = float64(0)
var ZeroInt64 = int64(0)

// RandString 取得随机字符包含数字、大小写等，可以自己随意扩展。
func RandString(l int) string {
	var inibyte []byte
	var result bytes.Buffer
	for i := 48; i < 123; i++ {
		switch {
		case i < 58:
			inibyte = append(inibyte, byte(i))
		case i >= 97 && i < 123:
			inibyte = append(inibyte, byte(i))
		}
	}
	var temp byte
	for i := 0; i < l; {
		if inibyte[randInt(0, len(inibyte))] != temp {
			temp = inibyte[randInt(0, len(inibyte))]
			result.WriteByte(temp)
			i++
		}
	}

	return result.String()
}

func randInt(min int, max int) byte {
	rand.Seed(time.Now().UnixNano())
	return byte(min + rand.Intn(max-min))
}

func GetServerIp() (string, error) {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrList {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("获取服务器IP失败")
}

func FileDownload(c *gin.Context, fileAlias, fileName string) {
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			Log.Error(c, err.Error())
		}
	}()
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileAlias))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fileName)
}

// GetCallerInfoForLog 获取调用者的函数名、文件名、行号
func GetCallerInfoForLog(skip int) map[string]interface{} {

	pc, file, line, ok := runtime.Caller(skip) // 回溯层数
	if !ok {
		return nil
	}
	ret := make(map[string]interface{})
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名
	ret["func"] = funcName
	ret["file"] = file
	ret["line"] = line
	return ret
}

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
