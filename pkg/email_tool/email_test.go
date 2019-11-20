package email_tool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendEmail(t *testing.T) {
	err := SendEmail([]string{"810909753@qq.com"}, "test", "test")
	assert.Equal(t, nil, err)
}
