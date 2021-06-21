package api

import (
	"fusion-gin-admin/config"
	"fusion-gin-admin/ginx"
	"fusion-gin-admin/lib/errors"
	"fusion-gin-admin/lib/logger"
	"fusion-gin-admin/schema"
	"fusion-gin-admin/service"
	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

type Login struct {
	LoginService *service.Login
}

func (a *Login) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginService.GetCaptcha(ctx, config.C.Captcha.Length)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *Login) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		ginx.ResError(c, errors.New400Response("请提供验证码ID"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginx.ResError(c, errors.New400Response("未找到验证码ID"))
			return
		}
	}

	cfg := config.C.Captcha
	err := a.LoginService.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginx.ResError(c, err)
	}
}

func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginx.ResError(c, errors.New400Response("无效的验证码"))
		return
	}

	user, err := a.LoginService.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	userID := user.ID
	// 将用户ID放入上下文
	ginx.SetUserID(c, userID)

	tokenInfo, err := a.LoginService.GenerateToken(ctx, userID)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("登入系统")
	ginx.ResSuccess(c, tokenInfo)
}

func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := ginx.GetUserID(c)
	if userID != "" {
		ctx = logger.NewTagContext(ctx, "__logout__")
		err := a.LoginService.DestroyToken(ctx, ginx.GetToken(c))
		if err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
		logger.WithContext(ctx).Infof("登出系统")
	}
	ginx.ResOK(c)
}

func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginService.GenerateToken(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, tokenInfo)
}

func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginService.GetLoginInfo(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, info)
}

func (a *Login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginService.QueryUserMenuTree(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, menus)
}

func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.LoginService.UpdatePassword(ctx, ginx.GetUserID(c), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
