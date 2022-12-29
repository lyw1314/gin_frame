package v1

import (
	"gin_frame/model"
	_ "gin_frame/model"
	"gin_frame/pkg/e"
	"gin_frame/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接收接口传参
type reqList struct {
	Title    string `json:"title" form:"title"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	From     int    `json:"from" form:"from" binding:"omitempty,numeric,eq=1|eq=2"`
}

type reqAdd struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Desc    string `json:"desc" form:"desc" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
}

// GetList 获取列表
func GetList(c *gin.Context) {
	appG := util.Gin{C: c}
	// 参数合法性校验
	var request reqList
	if err := c.ShouldBind(&request); err != nil {
		// 写日志
		util.Log.Error(c, err.Error())

		appG.Response(http.StatusOK, e.INVALID_PARAMS, "")
		return
	}

	var blogModel model.BlogArticle
	var conf map[string]interface{}
	if request.Title != "" {
		conf = map[string]interface{}{
			"title": request.Title,
		}
	}

	list, _ := blogModel.GetList(conf)
	appG.Response(http.StatusOK, e.SUCCESS, list)
}

// Add 添加数据
func Add(c *gin.Context) {
	appG := util.Gin{C: c}
	// 参数合法性校验
	var request reqAdd
	if err := c.ShouldBind(&request); err != nil {
		// 写日志
		util.Log.Error(c, err.Error())

		// 返回结果
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "")
		return
	}

	// 声明要操作的model
	var blogModel model.BlogArticle

	// 填充结构体
	blogModel.Title = request.Title
	blogModel.Content = request.Content
	blogModel.Desc = request.Desc

	// 调用添加方法
	err := blogModel.Add()

	if err != nil {
		// 写日志
		util.Log.Error(c, err.Error())
		appG.Response(http.StatusOK, e.HANDLE_FAIL, "")
	}
	appG.Response(http.StatusOK, e.SUCCESS, "")
}

func UserRedis(c *gin.Context) {
	model.RedisHandle["master"].Get(c, "ddd")
}
