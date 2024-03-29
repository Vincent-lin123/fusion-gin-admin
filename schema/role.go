package schema

import "time"

type Role struct {
	ID        string    `json:"id"`                                    // 唯一标识
	Name      string    `json:"name" binding:"required"`               // 角色名称
	Sequence  int       `json:"sequence"`                              // 排序值
	Memo      string    `json:"memo"`                                  // 备注
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 状态(1:启用 2:禁用)
	Creator   string    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                            // 更新时间
	RoleMenus RoleMenus `json:"role_menus" binding:"required,gt=0"`    // 角色菜单列表
}

type RoleQueryParam struct {
	PaginationParam
	IDs        []string `form:"-"`          // 唯一标识列表
	Name       string   `form:"-"`          // 角色名称
	QueryValue string   `form:"queryValue"` // 模糊查询
	UserID     string   `form:"-"`          // 用户ID
	Status     int      `form:"status"`     // 状态(1:启用 2:禁用)
}

type RoleQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

type RoleQueryResult struct {
	Data       Roles
	PageResult *PaginationResult
}

type Roles []*Role

func (r Roles) ToNames() []string {
	names := make([]string, len(r))
	for i, item := range r {
		names[i] = item.Name
	}
	return names
}

func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

type RoleMenu struct {
	ID       string `json:"id"`                           // 唯一标识
	RoleID   string `json:"role_id" binding:"required"`   // 角色ID
	MenuID   string `json:"menu_id" binding:"required"`   // 菜单ID
	ActionID string `json:"action_id" binding:"required"` // 动作ID
}

type RoleMenuQueryParam struct {
	PaginationParam
	RoleID  string   // 角色ID
	RoleIDs []string // 角色ID列表
}

type RoleMenuQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *PaginationResult
}

type RoleMenus []*RoleMenu

func (a RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range a {
		m[item.MenuID+"-"+item.ActionID] = item
	}
	return m
}

func (a RoleMenus) ToRoleIDMap() map[string]RoleMenus {
	m := make(map[string]RoleMenus)
	for _, item := range a {
		m[item.RoleID] = append(m[item.RoleID], item)
	}
	return m
}

func (a RoleMenus) ToMenuIDs() []string {
	var idList []string
	m := make(map[string]struct{})

	for _, item := range a {
		if _, ok := m[item.MenuID]; ok {
			continue
		}
		idList = append(idList, item.MenuID)
		m[item.MenuID] = struct{}{}
	}

	return idList
}

func (a RoleMenus) ToActionIDs() []string {
	idList := make([]string, len(a))
	m := make(map[string]struct{})
	for i, item := range a {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList[i] = item.ActionID
		m[item.ActionID] = struct{}{}
	}
	return idList
}
