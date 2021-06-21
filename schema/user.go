package schema

import (
	"fusion-gin-admin/config"
	"fusion-gin-admin/lib/util/hash"
	"fusion-gin-admin/lib/util/json"
	"fusion-gin-admin/lib/util/structure"
	"time"
)

type User struct {
	ID        string    `json:"id"`                                    // 唯一标识
	UserName  string    `json:"user_name" binding:"required"`          // 用户名
	RealName  string    `json:"real_name" binding:"required"`          // 真实姓名
	Password  string    `json:"password"`                              // 密码
	Phone     string    `json:"phone"`                                 // 手机号
	Email     string    `json:"email"`                                 // 邮箱
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 用户状态(1:启用 2:停用)
	Creator   string    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UserRoles UserRoles `json:"user_roles" binding:"required,gt=0"`    // 角色授权
}

func GetRootUser() *User {
	user := config.C.Root
	return &User{
		ID:       user.UserName,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: hash.MD5String(user.Password),
	}
}

func CheckIsRootUser(userID string) bool {
	return GetRootUser().ID == userID
}

func (a *User) String() string {
	return json.MarshalToString(a)
}

func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

type UserQueryParam struct {
	PaginationParam
	UserName   string   `form:"userName"`   // 用户名
	QueryValue string   `form:"queryValue"` // 模糊查询
	Status     int      `form:"status"`     // 用户状态(1:启用 2:停用)
	RoleIDs    []string `form:"-"`          // 角色ID列表
}

type UserQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

type UserRole struct {
	ID     string `json:"id"`      // 唯一标识
	UserID string `json:"user_id"` // 用户ID
	RoleID string `json:"role_id"` // 角色ID
}

type UserRoleQueryParam struct {
	PaginationParam
	UserID  string   // 用户ID
	UserIDs []string // 用户ID列表
}

type UserRoleQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *PaginationResult
}

func (uqr UserQueryResult) ToShowResult(mUserRoles map[string]UserRoles, mRoles map[string]*Role) *UserShowQueryResult {
	return &UserShowQueryResult{
		PageResult: uqr.PageResult,
		Data:       uqr.Data.ToUserShows(mUserRoles, mRoles),
	}
}

func (us Users) ToUserShows(mUserRoles map[string]UserRoles, mRoles map[string]*Role) UserShows {
	list := make(UserShows, len(us))
	for i, item := range us {
		showItem := new(UserShow)
		structure.Copy(item, showItem)
		for _, roleID := range mUserRoles[item.ID].ToRoleIDs() {
			if v, ok := mRoles[roleID]; ok {
				showItem.Roles = append(showItem.Roles, v)
			}
		}
		list[i] = showItem
	}

	return list
}

type Users []*User
type UserRoles []*UserRole
type UserShows []*UserShow

func (a Users) ToIDS() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}

func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

type UserShow struct {
	ID        string    `json:"id"`         // 唯一标识
	UserName  string    `json:"user_name"`  // 用户名
	RealName  string    `json:"real_name"`  // 真实姓名
	Phone     string    `json:"phone"`      // 手机号
	Email     string    `json:"email"`      // 邮箱
	Status    int       `json:"status"`     // 用户状态(1:启用 2:停用)
	CreatedAt time.Time `json:"created_at"` // 创建时间
	Roles     []*Role   `json:"roles"`      // 授权角色列表
}

type UserShowQueryResult struct {
	Data       UserShows
	PageResult *PaginationResult
}
