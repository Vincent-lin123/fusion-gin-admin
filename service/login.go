package service

import (
	"context"
	"fusion-gin-admin/lib/auth"
	"fusion-gin-admin/lib/errors"
	"fusion-gin-admin/lib/util/hash"
	"fusion-gin-admin/model/gormx/model"
	"fusion-gin-admin/schema"
	"github.com/LyricTian/captcha"
	"github.com/google/wire"
	"net/http"
	"sort"
)

var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

type Login struct {
	Auth            auth.Auther
	UserModel       *model.User
	UserRoleModel   *model.UserRole
	RoleModel       *model.Role
	RoleMenuModel   *model.RoleMenu
	MenuModel       *model.Menu
	MenuActionModel *model.MenuAction
}

func (a *Login) GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error) {
	captchaID := captcha.NewLen(length)
	item := &schema.LoginCaptcha{
		CaptchaID: captchaID,
	}
	return item, nil
}

func (a *Login) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	err := captcha.WriteImage(w, captchaID, width, height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.WithStack(err)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

func (a *Login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	root := schema.GetRootUser()
	if userName == root.UserName && root.Password == password {
		return root, nil
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.ErrInvalidUserName
	}

	item := result.Data[0]
	if item.Password != hash.SHA1String(password) {
		return nil, errors.ErrInvalidPassword
	} else if item.Status != 1 {
		return nil, errors.ErrUserDisable
	}

	return item, nil
}

func (a *Login) GenerateToken(ctx context.Context, userID string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

func (a *Login) DestroyToken(ctx context.Context, tokenString string) error {
	err := a.Auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Login) checkAndGetUser(ctx context.Context, userID string) (*schema.User, error) {
	user, err := a.UserModel.Get(ctx, userID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.ErrInvalidUser
	} else if user.Status != 1 {
		return nil, errors.ErrUserDisable
	}
	return user, nil
}

func (a *Login) GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error) {
	if isRoot := schema.CheckIsRootUser(userID); isRoot {
		root := schema.GetRootUser()
		loginInfo := &schema.UserLoginInfo{
			UserName: root.UserName,
			RealName: root.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	info := &schema.UserLoginInfo{
		UserID:   user.ID,
		UserName: user.UserName,
		RealName: user.RealName,
	}

	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if roleIDs := userRoleResult.Data.ToRoleIDs(); len(roleIDs) > 0 {
		roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}
		info.Roles = roleResult.Data
	}

	return info, nil
}

func (a *Login) QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error) {
	isRoot := schema.CheckIsRootUser(userID)
	// 如果是root用户，则查询所有显示的菜单树
	if isRoot {
		result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
			Status: 1,
		}, schema.MenuQueryOptions{
			OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
		})
		if err != nil {
			return nil, err
		}

		menuActionResult, err := a.MenuActionModel.Query(ctx, schema.MenuActionQueryParam{})
		if err != nil {
			return nil, err
		}
		return result.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap()).ToTree(), nil
	}

	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	} else if len(userRoleResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	roleMenuResult, err := a.RoleMenuModel.Query(ctx, schema.RoleMenuQueryParam{
		RoleIDs: userRoleResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(roleMenuResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	menuResult, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		IDS:    roleMenuResult.Data.ToMenuIDs(),
		Status: 1,
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	mData := menuResult.Data.ToMap()
	// 获取授权菜单的父级菜单，判断哪些父级菜单不在之前的授权菜单中，存放于qIDs切片
	var qIDs []string
	for _, pid := range menuResult.Data.SplitParentIDS() {
		if _, ok := mData[pid]; !ok {
			qIDs = append(qIDs, pid)
		}
	}
	// 获取这些差异的父级菜单的信息，补充到menuResult.Data中
	if len(qIDs) > 0 {
		pmenuResult, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
			IDS: qIDs,
		})
		if err != nil {
			return nil, err
		}
		menuResult.Data = append(menuResult.Data, pmenuResult.Data...)
	}

	sort.Sort(menuResult.Data)
	menuActionResult, err := a.MenuActionModel.Query(ctx, schema.MenuActionQueryParam{
		IDs: roleMenuResult.Data.ToActionIDs(),
	})
	if err != nil {
		return nil, err
	}
	return menuResult.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap()).ToTree(), nil
}

func (a *Login) UpdatePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error {
	if schema.CheckIsRootUser(userID) {
		return errors.New400Response("root用户不允许更新密码")
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return err
	} else if hash.SHA1String(params.OldPassword) != user.Password {
		return errors.New400Response("旧密码不正确")
	}

	params.NewPassword = hash.SHA1String(params.NewPassword)
	return a.UserModel.UpdatePassword(ctx, userID, params.NewPassword)
}
