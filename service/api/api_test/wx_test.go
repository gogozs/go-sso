package api_test

import (
	"github.com/magiconair/properties/assert"
	"go-weixin/pkg/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)


func TestWx(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/wx", nil)
	router.ServeHTTP(w, req)
	var b bool
	res, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal(res, &b)
	assert.Equal(t, false, b)
}