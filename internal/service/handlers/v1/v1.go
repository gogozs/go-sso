package v1

import (
	"go-sso/internal/service/handlers/viewset"
)

func NewAuthViewset() *AuthViewset {
	vs := &viewset.ViewSet{}
	authVS := &AuthViewset{
		ViewSet: *vs,
	}
	return authVS
}
