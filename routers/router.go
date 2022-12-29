package routers

import (
	"gin_frame/middleware/access"
	v1 "gin_frame/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	//健康检查
	//r.Group("/api/v1/health").GET("/status", api.Health)

	// --------------中间件 start--------------
	r.Use(access.Logger())   //记录访问日志
	r.Use(access.Recovery()) //错误日志收集
	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	//r.Use(requestID.RequestIdMiddleware())
	//r.Use(access.CsrfMiddleware())       //csrf验证
	//r.Use(access.LoginCheckMiddleware()) //登录状态检查
	//r.Use(access.AclMiddleware())        //访问权限检查
	//r.Use(err.ErrorsTransMiddleware())   //错误翻译
	// --------------中间件 end--------------

	// -------------设置路由-------------

	// 获取token
	r.POST("/api/v1/auth/getToken", v1.GetToken)
	r.GET("/api/v1/auth/parseToken", v1.ParseToken)

	demo := r.Group("/api/v1/demo") //demo
	//demo.Use(access.JWT())          // 添加JWT验证
	{
		demo.GET("/getList", v1.GetList)
		demo.POST("/add", v1.Add)
	}

	return r
}
