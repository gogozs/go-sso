package apierror

var MsgFlags = map[int]string{
	SUCCESS:        "OK",
	ERROR:          "Fail",
	NOT_FOUND:      "请求资源不存在",
	INVALID_PARAMS: "请求参数错误",
	ERROR_AUTH:     "权限不足，拒绝访问",
	UNAUTHORIZED:   "身份认证信息未提供",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
