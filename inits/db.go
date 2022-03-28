package inits

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"hitokoto-go/global"
)

func DB() error {
	var err error
	var gormConfig gorm.Config
	if global.Config.IsProdMode {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	global.DB, err = gorm.Open(postgres.Open(global.Config.PGConnString), &gormConfig)
	return err
}
