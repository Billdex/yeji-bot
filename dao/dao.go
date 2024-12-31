package dao

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"yeji-bot/config"
)

var DB *gorm.DB

func InitDao() error {
	err := initDatabase(config.AppConfig.DB.DSN)
	if err != nil {
		return err
	}

	return nil
}

func initDatabase(dsn string) (err error) {
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // enable color
			},
		),
	})
	if err != nil {
		return errors.Wrap(err, "gorm.Open fail")
	}

	db, err := DB.DB()
	if err != nil {
		return errors.Wrap(err, "load gorm.DB() fail")
	}
	db.SetMaxIdleConns(8)
	db.SetMaxOpenConns(16)
	db.SetConnMaxLifetime(8 * time.Hour)

	return nil
}

func ReloadAllCache(ctx context.Context) error {
	for _, fn := range []func(context.Context) error{
		ReloadChefs,
		ReloadRecipes,
		ReloadComboRecipes,
		ReloadEquips,
		ReloadGifts,
		ReloadSkills,
		ReloadQuests,
		ReloadTarots,
	} {
		err := fn(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
