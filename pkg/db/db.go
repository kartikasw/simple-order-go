package db

import (
	"fmt"

	cfg "simple-order-go/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DSN = "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"

func InitDB(d cfg.Database) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(fmt.Sprintf(
			DSN,
			d.Host,
			d.User,
			d.Password,
			d.Name,
			d.Port,
			d.SslMode,
			d.Timezone,
		)),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
		},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
