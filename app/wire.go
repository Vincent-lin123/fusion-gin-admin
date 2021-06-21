package app

import (
	"fusion-gin-admin/api"
	"fusion-gin-admin/model/gormx/model"
	"fusion-gin-admin/module/adapter"
	"fusion-gin-admin/router"
	"fusion-gin-admin/service"
	"github.com/google/wire"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		// mock.MockSet,
		InitGormDB,
		model.ModelSet,
		InitAuth,
		InitCasbin,
		InitGinEngine,
		service.ServiceSet,
		api.APISet,
		router.RouterSet,
		adapter.CasbinAdapterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
