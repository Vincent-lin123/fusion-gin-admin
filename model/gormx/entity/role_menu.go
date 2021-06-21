package entity

import (
	"context"
	"fusion-gin-admin/lib/util/structure"
	"fusion-gin-admin/schema"
	"github.com/jinzhu/gorm"
)

func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(RoleMenu))
}

type SchemaRoleMenu schema.RoleMenu

func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := new(RoleMenu)
	structure.Copy(a, item)
	return item
}

type RoleMenu struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`   // 角色ID
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`   // 菜单ID
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"` // 动作ID
}

func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	structure.Copy(a, item)
	return item
}

type RoleMenus []*RoleMenu

func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}
