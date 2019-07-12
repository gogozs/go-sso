package sdk

const (
	SUCCESS = 0
	SYSTEM_BUSY = -1
	INVALID_APPID = 40013
	SECRET_ERROR = 40001
	GRANTTYPE_ERROR = 40002
	IP_NOT_ALLOW = 40164
	UN_KNOWN = -500
)

var message = make(map[int]string)

func init() {
	// login
	message[SUCCESS] = "请求成功"
	message[SYSTEM_BUSY] = "系统繁忙"
	message[INVALID_APPID] = "invalid appid"
	message[SECRET_ERROR] = "AppSecret错误或者AppSecret不属于这个公众号，请开发者确认AppSecret的正确性"
	message[GRANTTYPE_ERROR] = "请确保grant_type字段值为client_credential"
	message[IP_NOT_ALLOW] = "调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置。（小程序及小游戏调用不要求IP地址在白名单内。）"
	message[UN_KNOWN] = "SDK未知错误"
}

func GetMessage(errorCode int) string {
	if msg, ok := message[errorCode]; ok {
		return msg
	}
	return message[UN_KNOWN]
}

