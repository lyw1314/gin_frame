package access

import (
	"gin_frame/pkg/e"
	"gin_frame/pkg/util"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// 访问权限控制
func CsrfMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := util.Gin{C: c}

		//referer验证，防csrf
		refer, err := url.Parse(c.Request.Referer())
		if err != nil {
			util.Log.Warning(c, "referer解析异常"+c.Request.Referer()+" "+err.Error())
			appG.Response(http.StatusOK, e.REFRER_ILLEGAL, nil)
			c.Abort()
			return
		}

		if refer.Host != c.Request.Host {
			util.Log.Warning(c, "referer非法 "+c.Request.Referer())
			appG.Response(http.StatusOK, e.REFRER_ILLEGAL, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
