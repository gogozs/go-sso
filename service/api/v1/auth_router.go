package v1

import (
	"go-weixin/service/api"
)


func init() {
	router := api.GetRouter()
	apiAuth := router.Group("api/v1/auth/")

	{
		apiAuth.POST("login/", ViewLogin)
		apiAuth.POST("register/", ViewRegister)
	}
}
