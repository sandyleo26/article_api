package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sandyleo26/article_api/article"
)

func main() {

	http.HandleFunc("/articles", postArticleHandler)

	port := ":4321"
	log.Println("Starting web server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

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

func postArticleHandler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "POST") {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&article.Article{
			Title: "post",
		})
	} else if strings.EqualFold(r.Method, "GET") {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&article.Article{
			Title: "get",
		})
	}
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
