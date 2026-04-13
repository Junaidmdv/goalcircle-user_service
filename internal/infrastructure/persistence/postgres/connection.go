package postgres

import (
	"fmt"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type postgressDB struct {
	DB *gorm.DB
}

func NewDatabase(db *config.PostgresConfig) (*postgressDB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", db.Host, db.User, db.Password, db.DBName, db.Port)
	psqlInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("connection error:%v", err)
	}

	return &postgressDB{
		DB: psqlInstance,
	}, nil
}
