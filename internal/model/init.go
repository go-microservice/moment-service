package model

import (
	"fmt"
	"sync"

	"github.com/go-eagle/eagle/pkg/config"
	"github.com/go-eagle/eagle/pkg/storage/orm"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Once sync.Once
)

// Init init db
func Init() (*gorm.DB, func()) {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("load orm conf err: %v", err))
	}

	DB, cleanup := orm.NewMySQL(cfg)
	return DB, cleanup
}

// GetDB get a db instance
func GetDB() *gorm.DB {
	Once.Do(func() {
		DB, _ = Init()
	})
	return DB
}

// loadConf load gorm config
func loadConf() (ret *orm.Config, err error) {
	var cfg orm.Config
	if err := config.Load("database", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
