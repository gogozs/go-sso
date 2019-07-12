package v1

import (
	"github.com/gin-gonic/gin"
	"go-weixin/pkg/api_error"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 封装通用response
// Response 返回的数据
func (g *Gin) Response(httpCode, msgCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: msgCode,
		Data: data,
		Msg:  api_error.GetMsg(msgCode),
	})
	return
}

func (g *Gin) SuccessResponse(data interface{}) {
	g.C.JSON(200, Response{
		Code: api_error.SUCCESS,
		Data: data,
		Msg:  api_error.GetMsg(api_error.SUCCESS),
	})
	return
}

func (g *Gin) FailResponse(errCode int) {
	g.C.JSON(400, Response{
		Code: -1,
		Msg:  api_error.GetMsg(errCode),
	})
	return
}

func (g *Gin) NotFoundResponse(errCode int) {
	g.C.JSON(404, Response{
		Code: api_error.NOT_FOUND,
		Msg:  api_error.GetMsg(errCode),
	})
	return
}

// 处理post data
//func (g *Gin) GetBodyData(obj interface{}) (interface{}, error) {
//	var err error
//	contentType := g.C.Request.Header.Get("Content-Type")
//	switch contentType {
//	case "application/json":
//		body, _ := ioutil.ReadAll(g.C.Request.Body)
//		err = json.Unmarshal(body, &obj)
//	case "application/x-www-form-urlencoded":
//		err = g.C.BindWith(&obj, binding.Form)
//	default:
//		err = g.C.Bind(&obj) //使用自动推断
//	}
//	return obj, err
//}

func (g *Gin) List() {

}

func (g *Gin) Retreive() {

}

func (g *Gin) Create() {

}

func (g *Gin) Patch() {

}

func (g *Gin) Update() {

}

func (g *Gin) Delete() {

}
