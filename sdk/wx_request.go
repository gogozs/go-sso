package sdk

import (
	"errors"
	"github.com/fatih/structs"
	"go-weixin/pkg/api_request"
	"go-weixin/pkg/json"
	"net/http"
)

const (
	apiUrl = "https://api.weixin.qq.com"
)

// 获取 access_token
func GetWxToken(access *WxAccess) error {
	// {"access_token":"ACCESS_TOKEN","expires_in":7200}
	res, statusCode, _ := api_request.Get("/cgi-bin/token", nil, structs.Map(access))
	if statusCode == 200 {
		var data map[string]interface{}
		json.Unmarshal(res, &data)
		if _, ok := data["access_token"];ok {
			access.access_token = data["access_token"].(string)
		} else {
			errcode := data["errcode"].(int)
			return errors.New(GetMessage(errcode))
		}
		return nil
	} else {
		return http.ErrHandlerTimeout
	}
}
