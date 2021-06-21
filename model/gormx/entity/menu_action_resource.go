package entity

import (
	"context"
	"fusion-gin-admin/lib/util/structure"
	"fusion-gin-admin/schema"
	"github.com/jinzhu/gorm"
)

func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(MenuActionResource))
}

type SchemaMenuActionResource schema.MenuActionResource

func (a SchemaMenuActionResource) ToMenuActionResource() *MenuActionResource {
	item := new(MenuActionResource)
	structure.Copy(a, item)
	return item
}

type MenuActionResource struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"` // 菜单动作ID
	Method   string `gorm:"column:method;size:100;default:'';not null;"`         // 资源请求方式(支持正则)
	Path     string `gorm:"column:path;size:100;default:'';not null;"`           // 资源请求路径（支持/:id匹配）
}

func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	structure.Copy(a, item)
	return item
}

type MenuActionResources []*MenuActionResource

func (a MenuActionResources) ToSchemaMenuActionResources() []*schema.MenuActionResource {
	list := make([]*schema.MenuActionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuActionResource()
	}
	return list
}
