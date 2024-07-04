package repository

import (
	"fmt"

	"github.com/Egorpalan/time-tracker/internal/models"
	"github.com/Egorpalan/time-tracker/pkg/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	var err error
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)

	DB, err = gorm.Open("postgres", dbInfo)
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.Task{})
}
