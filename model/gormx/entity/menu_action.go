package entity

import (
	"context"
	"fusion-gin-admin/lib/util/structure"
	"fusion-gin-admin/schema"
	"github.com/jinzhu/gorm"
)

func GetMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(MenuAction))
}

type SchemaMenuAction schema.MenuAction

func (a SchemaMenuAction) ToMenuAction() *MenuAction {
	item := new(MenuAction)
	structure.Copy(a, item)
	return item
}

type MenuAction struct {
	ID     string `gorm:"column:id;primary_key;size:36;"`
	MenuID string `gorm:"column:menu_id;size:36;index;default:'';not null;"` // 菜单ID
	Code   string `gorm:"column:code;size:100;default:'';not null;"`         // 动作编号
	Name   string `gorm:"column:name;size:100;default:'';not null;"`         // 动作名称
}

func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	structure.Copy(a, item)
	return item
}

type MenuActions []*MenuAction

func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}
