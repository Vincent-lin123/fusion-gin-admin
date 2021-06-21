package middleware

import (
	"fusion-gin-admin/config"
	"fusion-gin-admin/contextx"
	"fusion-gin-admin/ginx"
	"fusion-gin-admin/lib/auth"
	"fusion-gin-admin/lib/errors"
	"fusion-gin-admin/lib/logger"
	"github.com/gin-gonic/gin"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	ginx.SetUserID(c, userID)
	ctx := contextx.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, config.C.Root.UserName)
			c.Next()
		}
	}

	return func(context *gin.Context) {
		if SkipHandler(context, skippers...) {
			context.Next()
		}

		userID, err := a.ParseUserID(context.Request.Context(), ginx.GetToken(context))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config.C.IsDebugMode() {
					wrapUserAuthContext(context, config.C.Root.UserName)
					context.Next()
					return
				}
				ginx.ResError(context, errors.ErrInvalidToken)
				return
			}
			ginx.ResError(context, errors.WithStack(err))
			return
		}

		wrapUserAuthContext(context, userID)
		context.Next()
	}
}
