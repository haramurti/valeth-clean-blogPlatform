package config

import (
	"fmt"
	"os"
	"valeth-clean-blogPlatform/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {

	dsn := os.Getenv("DB_DSN")
	if dsn == " " {
		fmt.Println("error : empty dsn env")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Post{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}
