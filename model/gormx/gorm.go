package gormx

import (
	"fusion-gin-admin/config"
	"fusion-gin-admin/lib/logger"
	"fusion-gin-admin/model/gormx/entity"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

func NewDB(c *Config) (*gorm.DB, func(), error) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return c.TablePrefix + defaultTableName
	}

	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		return nil, nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	cleanFunc := func() {
		err := db.Close()
		if err != nil {
			logger.Errorf("Gorm db close error: %s", err.Error())
		}
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, cleanFunc, err
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, cleanFunc, nil
}

//数据库自动映射
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(entity.MenuAction),
		new(entity.MenuActionResource),
		new(entity.Menu),
		new(entity.RoleMenu),
		new(entity.Role),
		new(entity.UserRole),
		new(entity.User),
	).Error
}
