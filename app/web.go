package app

import (
	"fusion-gin-admin/config"
	"fusion-gin-admin/middleware"
	"fusion-gin-admin/router"
	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	prefixes := r.Prefixes()

	// Trace ID
	app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Copy body
	app.Use(middleware.CopyBodyMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Access logger
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Recover
	app.Use(middleware.RecoveryMiddleware())

	// CORS
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// GZIP
	if config.C.GZIP.Enable {
		app.Use(gzip.Gzip(gzip.BestCompression,
			gzip.WithExcludedExtensions(config.C.GZIP.ExcludedExtentions),
			gzip.WithExcludedPaths(config.C.GZIP.ExcludedPaths),
		))
	}

	// Router register
	r.Register(app)

	// Swagger
	if config.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Website
	if dir := config.C.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(prefixes...)))
	}

	return app
}
