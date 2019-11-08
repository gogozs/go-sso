package v1

import (
	"go-qiuplus/db/query"
	"go-qiuplus/service/api/viewset"
)

func GetAuthVS() *AuthViewset {
	vs := &viewset.ViewSet{}
	authVS := &AuthViewset{
		itemInter: query.UserQ,
		ViewSet:   *vs,
	}
	return authVS
}

func GetWxVS() *WxViewset {
	vs := &viewset.ViewSet{}
	wxVS := &WxViewset{
		ViewSet: *vs,
	}
	return wxVS
}
