package v1

import (
	"go-sso/db/inter"
	"go-sso/service/api/viewset"
)

func GetAuthVS() *AuthViewset {
	vs := &viewset.ViewSet{}
	authVS := &AuthViewset{
		itemInter: inter.GetDao(),
		ViewSet:   *vs,
	}
	return authVS
}
