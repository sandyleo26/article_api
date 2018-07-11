package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sandyleo26/article_api/article"
)

func main() {
	fmt.Println("Hello World!")

	db := OpenDB()
	defer db.Close()

	articleRequest := article.ArticleRequest{
		Title: "some title",
		Body:  "some body",
		Tags:  []string{"debug", "abc"},
	}
	article := article.NewArticle(1, &articleRequest)
	result := db.Debug().Create(article)
	fmt.Println(result)
}

//OpenDB open database
func OpenDB() *gorm.DB {
	dbHost := "localhost"
	dbPort := "5432"
	dbName := "blueco"
	dbUser := "postgres"
	dbPass := "pass"
	dbSSL := "disable"

	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL)

	driver := "postgres"

	db, err := gorm.Open(driver, connString)
	if err != nil {
		fmt.Println("Failed to connect database "+connString+". Error: %v", err)
		panic("OpenDB error")
	}

	db.LogMode(true)
	return db
}
