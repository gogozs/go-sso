package api

import v1 "go-weixin/service/api/v1"

func WxRouterInit() {
	router := GetRouter()
	apiWx := router.Group("/")

	{
		apiWx.GET("wx", v1.ViewWx)
	}
}
