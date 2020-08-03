package auth

import (
	"fmt"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	a := &TokenAuthManager{}
	fmt.Println(a.RandomToken())
}