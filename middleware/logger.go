package middleware

import (
	"fusion-gin-admin/config"
	"fusion-gin-admin/ginx"
	"fusion-gin-admin/lib/logger"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"time"
)

func LoggerMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		method := c.Request.Method

		entry := logger.WithContext(logger.NewTagContext(c.Request.Context(), "__request__"))
		start := time.Now()
		fields := make(map[string]interface{})
		fields["ip"] = c.ClientIP()
		fields["method"] = method
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["header"] = c.Request.Header
		fields["user_agent"] = c.GetHeader("User-Agent")
		fields["content_length"] = c.Request.ContentLength

		if method == http.MethodPost || method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(c.GetHeader("Content-Type"))
			if mediaType != "multipart/form-data" {
				if v, ok := c.Get(ginx.ReqBodyKey); ok {
					if b, ok := v.([]byte); ok && len(b) <= config.C.HTTP.MaxLoggerLength {
						fields["body"] = string(b)
					}
				}
			}
		}
		c.Next()

		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		fields["resp_status"] = c.Writer.Status()
		fields["resp_length"] = c.Writer.Size()

		if v, ok := c.Get(ginx.LoggerReqBodyKey); ok {
			if b, ok := v.([]byte); ok && len(b) <= config.C.HTTP.MaxLoggerLength {
				fields["body"] = string(b)
			}
		}

		if v, ok := c.Get(ginx.RespBodyKey); ok {
			if b, ok := v.([]byte); ok && len(b) <= config.C.HTTP.MaxLoggerLength {
				fields["resp_body"] = string(b)
			}
		}

		fields[logger.UserIDKey] = ginx.GetUserID(c)
		entry.WithFields(fields).Infof("[http] %s-%s-%s-%d(%dms)",
			p, c.Request.Method, c.ClientIP(), c.Writer.Status(), timeConsuming)
	}
}
