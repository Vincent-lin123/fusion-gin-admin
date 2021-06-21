package middleware

import (
	"fmt"
	"fusion-gin-admin/ginx"
	"fusion-gin-admin/lib/errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func NoMethodHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		ginx.ResError(context, errors.ErrMethodNotAllow)
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ginx.ResError(c, errors.ErrNotFound)
	}
}

type SkipperFunc func(*gin.Context) bool

func AllowPathPrefixSkipper(prefixs ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixs {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(context *gin.Context) bool {
		path := JoinRouter(context.Request.Method, context.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

func EmptyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}
