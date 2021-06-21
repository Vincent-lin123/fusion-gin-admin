package model

import "github.com/google/wire"

var ModelSet = wire.NewSet(
	MenuActionResourceSet,
	MenuActionSet,
	MenuSet,
	RoleMenuSet,
	RoleSet,
	TransSet,
	UserRoleSet,
	UserSet,
)
