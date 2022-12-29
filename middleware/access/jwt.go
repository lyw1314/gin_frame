package access

import (
	"gin_frame/pkg/e"
	"gin_frame/pkg/util"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := util.Gin{C: c}
		var token string
		token = c.Query("token")
		if token == "" {
			token = c.GetHeader("Authorization")
		}

		if token == "" {
			util.Log.Warning(c, "token鉴权失败：token为空")
			appG.Response(http.StatusOK, e.AUTH_CHECK_TOKEN_FAIL, "")
			c.Abort()
			return
		}

		mJwt := util.NewJwt("gin-frame")
		parseToken, err := mJwt.ParseToken(token)
		if err != nil {
			var code int
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.AUTH_CHECK_TOKEN_FAIL
			}

			util.Log.Warning(c, "token鉴权失败："+err.Error())
			appG.Response(http.StatusOK, code, "")
			c.Abort()
			return
		}

		// 设置上下文，传递用户信息
		c.Set(util.LoginUserNameKey, parseToken.Username)

		c.Next()
	}
}
