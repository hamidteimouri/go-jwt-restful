package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	DbPort := os.Getenv("DB_PORT")

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUsername, DbPassword, DbHost, DbPort, DbName)

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
