package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-sso/pkg/email_tool"
	"go-sso/pkg/log"
	"net/http"
	"net/http/httputil"
	"time"
)

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}

func ErrEmailWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				errMsg := fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s", timeFormat(time.Now()), string(httpRequest), err)
				if err := email_tool.SendEmail(nil, "request error", errMsg); err != nil {
					log.Error("send email error %s", err.Error())
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
