package e

var MsgFlags = map[int]string{
	SUCCESS:                  "ok",
	ERROR:                    "fail",
	INVALID_PARAMS:           "请求参数错误",
	EXIST_TAG:                "已存在该标签名称",
	NOT_EXIST_TAG:            "该标签不存在",
	NOT_EXIST_ARTICLE:        "该文章不存在",
	AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	AUTH_CHECK_TOKEN_TIMEOUT: "Token已失效",
	AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:               "Token错误",
	ACL_NO_AUTH:              "无访问权限",
	HANDLE_FAIL:              "处理失败",
}

func GetMsg(code int) string {
	if msg, ok := MsgFlags[code]; ok {
		return msg
	}

	return MsgFlags[ERROR]
}
