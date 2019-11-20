package email_tool

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRegisterTpl(t *testing.T) {
	body := RegisterTmpl("25432")
	err := SendEmail([]string{"zhushen@datagrand.com"}, "test", body)
	assert.Equal(t, nil, err)
}
