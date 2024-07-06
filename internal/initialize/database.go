package initialize

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/wopoczynski/playground/internal/database"
)

type DBConfig struct {
	DSN string `env:"URL"`
}

func DB(cfg DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return db, nil
}

func Automigrate(ctx context.Context, db *gorm.DB) error {
	err := db.WithContext(ctx).AutoMigrate(&database.PersistingStruct{})
	if err != nil {
		return err
	}
	return nil
}
