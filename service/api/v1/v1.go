package v1

import (
	"go-sso/db/inter"
	"go-sso/service/api/viewset"
)

func NewAuthViewset() *AuthViewset {
	vs := &viewset.ViewSet{}
	authVS := &AuthViewset{
		itemInter: inter.GetQuery(),
		ViewSet:   *vs,
	}
	return authVS
}
