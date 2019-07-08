package api

import v1 "go-weixin/service/api/v1"

func AuthRouterInit() {
	router := GetRouter()
	apiAuth := router.Group("api/v1/auth/")

	{
		apiAuth.POST("login/", v1.ViewLogin)
		apiAuth.POST("register/", v1.ViewRegister)
	}
}
