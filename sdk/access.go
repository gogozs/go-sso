package sdk

import (
	"go-weixin/config"
	"go-weixin/util"
	"sort"
	"strconv"
	"strings"
)

type WxAccess struct {
	appid string
	secret string
	grant_type string
	access_token string
}

type TokenSignature struct {
	Signature string
	Timestamp int
	Nonce int
	Echostr string
}

// 签名验证
func (ts *TokenSignature) Confirm() bool {
	wxConfig := config.GetConfig().Weixin
	token := wxConfig.Token
	list := []string{strconv.Itoa(ts.Nonce), strconv.Itoa(ts.Timestamp), token}
	sort.Strings(list)
	confirmStr := strings.Join(list, "")
	shaHash := util.Sha1(confirmStr)
	if shaHash == ts.Signature {
		return true
	}
	return false
}


var wxAccess *WxAccess

func init() {
	access := config.GetConfig().Weixin
	wxAccess = &WxAccess{
		appid:      access.AppID,
		secret:     access.AppSecret,
		grant_type: access.GrantType,
	}
}

func GetAccessToken() string {
	return wxAccess.access_token
}
