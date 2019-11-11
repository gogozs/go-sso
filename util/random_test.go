package util

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRandomCode(t *testing.T) {
	for i:=0;i<10;i++ {
		c := RandomCode()
		assert.Equal(t, 6, len(c))
		ok, err := regexp.MatchString(`^\d{6}`, c)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, ok)
	}
}
