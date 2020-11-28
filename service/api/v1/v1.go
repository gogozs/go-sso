package v1

import (
	"go-sso/service/api/viewset"
)

func NewAuthViewset() *AuthViewset {
	vs := &viewset.ViewSet{}
	authVS := &AuthViewset{
		ViewSet: *vs,
	}
	return authVS
}
