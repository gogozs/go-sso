package viewset

import (
	"github.com/gin-gonic/gin"
	"go-sso/pkg/api_error"
	"net/http"
)

type ViewSet struct {
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (this *ViewSet) ErrorHandler(f func(c *gin.Context) error, c *gin.Context) {
	err := f(c)
	switch err.(type) {
	case nil:
	case *api_error.NotFoundError:
		this.NotFoundResponse(c)
	default:
		if e, ok := err.(api_error.ApiError); ok {
			this.FailResponse(c, e)
		} else {
			this.FailResponse(c, api_error.NewError(err))
		}
	}
}

func GetSuccessResponse(data interface{}) Response {
	return Response{
		Code: api_error.SUCCESS,
		Msg:  api_error.SUCCESS_MSG,
		Data: data,
	}
}

func GetFailResponse(err api_error.ApiError, data interface{}) Response {
	return Response{
		Code: err.Code(),
		Msg:  err.GetMsg(),
		Data: data,
	}
}

func (this *ViewSet) GetId(c *gin.Context) string {
	if i := c.Param("id"); i != "" {
		return i
	}
	return ""
}

// 封装通用response
// Response 返回的数据

func (this *ViewSet) SuccessResponse(c *gin.Context, data interface{}) error {
	c.JSON(http.StatusOK, GetSuccessResponse(data))
	return nil
}

func (this *ViewSet) SuccessBlackResponse(c *gin.Context) error {
	c.JSON(http.StatusOK, GetSuccessResponse(nil))
	return nil
}

func (this *ViewSet) SuccessListResponse(c *gin.Context, data interface{}, PageNum, PageSize, Total int) error {
	c.JSON(http.StatusOK,
		GetSuccessResponse(map[string]interface{}{
			"page_num":  PageNum,
			"page_size": PageSize,
			"total":     Total,
			"data":      data,
		}),
	)
	return nil
}

func (this *ViewSet) FailResponse(c *gin.Context, err api_error.ApiError, data ...interface{}) {
	c.JSON(http.StatusBadRequest, GetFailResponse(err, data))
	return
}

func (this *ViewSet) NotFoundResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, GetFailResponse(api_error.ErrNotFound, nil))
	return
}
