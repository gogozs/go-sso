package query

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

var (
	validTestCases = []struct{
		account string
		accountType string
		expected bool
	}{
		{"12345", "username", false},
		{"a12345", "username", true},
		{"a123@45", "username", false},
		{"a123@45", "telephone", false},
		{"1234567890", "telephone", false},
		{"12345678901", "telephone", false},
		{"18817012345", "telephone", true},
		{"12345678901", "email", false},
		{"a@22", "email", false},
		{"abc@aliyun.com", "email", true},
	}
)


func TestIsValid(t *testing.T) {
	uq := &UserQuery{}
	for _, tc := range validTestCases {
		assert.Equal(t, uq.IsValid(tc.account, tc.accountType), tc.expected)
	}
}