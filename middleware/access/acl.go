package access

import (
	"gin_frame/pkg/e"
	"gin_frame/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 访问权限控制
func AclMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := util.Gin{C: c}

		if c.GetInt(util.LoginUserRoleKey) == 0 {
			if c.Request.URL.Path == "/api/v1/login/by_phone" {
				util.Log.Warning(c, err.Error())
				appG.Response(http.StatusOK, e.ACL_NO_AUTH, nil)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
