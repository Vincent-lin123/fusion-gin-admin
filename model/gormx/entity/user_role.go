package entity

import (
	"context"
	"fusion-gin-admin/lib/util/structure"
	"fusion-gin-admin/schema"
	"github.com/jinzhu/gorm"
)

func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(UserRole))
}

type SchemaUserRole schema.UserRole

func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	structure.Copy(a, item)
	return item
}

type UserRole struct {
	ID     string `gorm:"column:id;primary_key;size:36;"`
	UserID string `gorm:"column:user_id;size:36;index;default:'';not null;"` // 用户内码
	RoleID string `gorm:"column:role_id;size:36;index;default:'';not null;"` // 角色内码
}

func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	structure.Copy(a, item)
	return item
}

type UserRoles []*UserRole

func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}
