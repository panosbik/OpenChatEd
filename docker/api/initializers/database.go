package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"OpenChatEd/models"
)

func ConnectDB(config *Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err = db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate the database")
	}
	log.Println("Successfully run the migrations")
	// Turn off the logger and prevent GORM from logging any SQL queries

	db.Logger = logger.Default.LogMode(logger.Silent)
	log.Println("Successfully connect to the database")
	return db
}
