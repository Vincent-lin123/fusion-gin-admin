package middleware

import (
	"bytes"
	"compress/gzip"
	"fusion-gin-admin/config"
	"fusion-gin-admin/ginx"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

func CopyBodyMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	var maxMemory int64 = 64 << 20 //64MB
	if v := config.C.HTTP.MaxContentLength; v > 0 {
		maxMemory = v
	}

	return func(context *gin.Context) {
		if SkipHandler(context, skippers...) || context.Request.Body == nil {
			context.Next()
			return
		}

		var requestBody []byte
		isGzip := false
		safe := &io.LimitedReader{R: context.Request.Body, N: maxMemory}

		if context.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(safe)
			if err == nil {
				isGzip = true
				requestBody, _ = ioutil.ReadAll(reader)
			}
		}

		if !isGzip {
			requestBody, _ = ioutil.ReadAll(safe)
		}

		context.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		context.Request.Body = http.MaxBytesReader(context.Writer, ioutil.NopCloser(bf), maxMemory)
		context.Set(ginx.ReqBodyKey, requestBody)

		context.Next()
	}
}
