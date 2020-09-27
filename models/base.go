package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {

	/* setting connection to database */
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
	}

	/* getting env variables */

	DbUsername := os.Getenv("DB_USERNAME")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbHost := os.Getenv("DB_HOST")

	dbURL := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbUsername, DbName, DbPassword)

	fmt.Println(dbURL)

	conn, err := gorm.Open("mysql", dbURL)

	if err != nil {
		fmt.Println("connection error", err)
	}
	db = conn

	db.Debug().AutoMigrate(&User{})
}

func getDB() *gorm.DB {
	/* get gORM db */
	return db
}
