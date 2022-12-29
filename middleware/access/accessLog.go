package access

import (
	"gin_frame/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 访问日志记录
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 日志格式
		util.Log.AccessInfo(c, "access_log", latencyTime)
	}
}
