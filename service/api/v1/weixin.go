package v1

/*
class Handle(object):
    def GET(self):
        try:
            data = web.input()
            if len(data) == 0:
                return "hello, this is handle view"
            signature = data.signature
            timestamp = data.timestamp
            nonce = data.nonce
            echostr = data.echostr
            token = "xxxx" #请按照公众平台官网\基本配置中信息填写

            list = [token, timestamp, nonce]
            list.sort()
            sha1 = hashlib.sha1()
            map(sha1.update, list)
            hashcode = sha1.hexdigest()
            print "handle/GET func: hashcode, signature: ", hashcode, signature
            if hashcode == signature:
                return echostr
            else:
                return ""
        except Exception, Argument:
            return Argument
 */

import (
	"github.com/gin-gonic/gin"
	"go-weixin/pkg/log"
	"go-weixin/sdk"
)

// token 认证
func ViewWx(c *gin.Context) {
	appG := Gin{C: c}
	s := sdk.TokenSignature{}
	err := c.BindQuery(&s)
	if err != nil {
		log.Error(err)
		appG.SuccessResponse(false)
	} else if s.Confirm() {
		appG.SuccessResponse(true)
	} else {
		appG.SuccessResponse(false)
	}
}
