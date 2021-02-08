package util

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRandomCode(t *testing.T) {
	reg, _ := regexp.Compile(`^\d{6}`)

	for i := 0; i < 10; i++ {
		c := RandomCode(6)
		assert.Equal(t, 6, len(c))
		ok := reg.Match([]byte(c))
		assert.Equal(t, true, ok)
	}
}
