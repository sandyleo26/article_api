package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

//OpenDB open database
func OpenDB() *gorm.DB {

	if db != nil {
		if err := db.DB().Ping(); err != nil {
			db.DB().Close()
		} else {
			return db
		}
	}

	dbHost := "localhost"
	dbPort := "5432"
	dbName := "blueco"
	dbUser := "postgres"
	dbPass := "pass"
	dbSSL := "disable"
	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL)
	driver := "postgres"

	var err error
	db, err = gorm.Open(driver, connString)
	if err != nil {
		fmt.Println("Failed to connect database "+connString+". Error: %v", err)
		panic("OpenDB error")
	}

	db.LogMode(true)
	log.Println("database connected!")
	return db
}
