package v1

import (
	"gin_frame/pkg/e"
	"gin_frame/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type reqUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// GetToken 获取token
func GetToken(c *gin.Context) {
	appG := util.Gin{C: c}

	var request reqUser
	if err := c.ShouldBind(&request); err != nil {
		// 写日志
		util.Log.Error(c, err.Error())

		appG.Response(http.StatusOK, e.INVALID_PARAMS, "")
		return
	}

	// 验证用户的合法性 todo
	// 一般去db验证是否存在

	// 分发token
	mJwt := util.NewJwt("gin-frame")

	//expireTime := nowTime

	token, _ := mJwt.CreateToken(request.UserName, "333", time.Hour*3)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"token": token})

}

func ParseToken(c *gin.Context) {
	appG := util.Gin{C: c}
	token := c.Query("token")
	mJwt := util.NewJwt("gin-frame")

	parseToken, _ := mJwt.ParseToken(token)
	appG.Response(http.StatusOK, e.SUCCESS, parseToken)

}
