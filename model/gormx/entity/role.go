package entity

import (
	"context"
	"fusion-gin-admin/lib/util/structure"
	"fusion-gin-admin/schema"
	"github.com/jinzhu/gorm"
	"time"
)

func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Role))
}

type SchemaRole schema.Role

func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	structure.Copy(a, item)
	return item
}

type Role struct {
	ID        string     `gorm:"column:id;primary_key;size:36;"`
	Name      string     `gorm:"column:name;size:100;index;default:'';not null;"` // 角色名称
	Sequence  int        `gorm:"column:sequence;index;default:0;not null;"`       // 排序值
	Memo      *string    `gorm:"column:memo;size:1024;"`                          // 备注
	Status    int        `gorm:"column:status;index;default:0;not null;"`         // 状态(1:启用 2:禁用)
	Creator   string     `gorm:"column:creator;size:36;"`                         // 创建者
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	structure.Copy(a, item)
	return item
}

type Roles []*Role

func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}
