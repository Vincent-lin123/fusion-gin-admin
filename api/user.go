package api

import (
	"fusion-gin-admin/ginx"
	"fusion-gin-admin/lib/errors"
	"fusion-gin-admin/schema"
	"fusion-gin-admin/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"strings"
)

var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

type User struct {
	UserService *service.User
}

func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}
	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	params.Pagination = true
	result, err := a.UserService.QueryShow(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserService.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item.CleanSecure())
}

func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	} else if item.Password == "" {
		ginx.ResError(c, errors.New400Response("密码不能为空"))
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.UserService.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.UserService.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserService.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *User) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserService.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *User) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserService.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
