package main

import (
	"bot/internal/config"
	"bot/internal/storage/postgres/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad(3)

	dsn := cfg.DB.ConnString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = runMigrations(db)
	if err != nil {
		panic(err)
	}
}

func runMigrations(db *gorm.DB) error {
	ok := db.Migrator().HasTable(&models.User{})
	if !ok {
		err := db.Migrator().CreateTable(&models.User{})
		if err != nil {
			return err
		}
		fmt.Println("user table created")
	}

	_ = dropColumns(db)

	err := db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return err
	}

	//users := getUsers(db)
	//for _, user := range users {
	//	_ = addSchedule(db, user.TelegramID)
	//}

	return nil
}

func dropColumns(db *gorm.DB) error {
	return nil
}

func getUsers(db *gorm.DB) []models.User {
	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		fmt.Println("Failed to get users", result.Error)
		return []models.User{}
	}

	return users
}
