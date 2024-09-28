package database

import (
	"fmt"
	"time"

	"loan-service/config"

	retry "github.com/avast/retry-go/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//nolint:gomnd
func GetDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		config.Data.DBHost,
		config.Data.DBPort,
		config.Data.DBUser,
		config.Data.DBName,
		config.Data.DBPassword,
	)

	var db *gorm.DB
	var err error
	err = retry.Do(
		func() error {
			var connErr error
			db, connErr = gorm.Open(postgres.Open(dsn), nil)
			return connErr
		},
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(1*time.Second),
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}
