package access

import (
	"gin_frame/pkg/setting"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type oauthValidateResult struct {
	Success bool
	Message string
}

var err error

func OauthAccess(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token参数来验证token是否合法
		token := c.PostForm("token")
		if token == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"errno": 1, "message": "身份验证失败"})
			return
		}
		// 去oauth服务器验证
		resp, err := http.Post(setting.OauthURI+"/site/validscope", "application/x-www-form-urlencoded",
			strings.NewReader("access_token="+token+"&scope="+scope+"get_uid=1"))
		if err != nil {
			log.Println(err)
			c.Abort()
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
			log.Println(err, body)
			c.Abort()
		}
		var res oauthValidateResult
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Println(err)
			c.Abort()
		}

		// 结果验证
		if res.Success == true {
			c.Next()
		} else {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"errno": 1, "message": "身份验证失败"})
			return // return也是可以省略的，执行了abort操作，会内置在中间件defer前，return，写出来也只是解答为什么Abort()之后，还能执行返回JSON数据
		}
	}
}
