package db

import (
	"fmt"
	"log"
	"time"

	"github.com/DanielStefanK/stream-bingo/config"
	"github.com/DanielStefanK/stream-bingo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var config = config.GetConfig()
	var dbConfig = config.Database

	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC"
	dsn = fmt.Sprintf(dsn, dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	// Auto migrate models
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to auto migrate models: %v", err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
