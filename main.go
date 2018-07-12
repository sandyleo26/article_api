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

var db *gorm.DB

func main() {

	// open db
	db = OpenDB()
	defer db.Close()

	// start web server
	http.HandleFunc("/articles", PostArticleHandler)
	port := ":4321"
	log.Println("Starting web server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

//PostArticleHandler create article
func PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "POST") {
		http.Error(w, "Method not supported", 400)
		return
	}

	articleRequest := new(article.ArticleRequest)
	if err := json.NewDecoder(r.Body).Decode(articleRequest); err != nil {
		http.Error(w, "Error when parsing request", 400)
		return
	}

	newArticle := article.NewArticle(1, articleRequest)
	result := db.Debug().Create(newArticle)
	if result.Error != nil {
		http.Error(w, "Error when creating article", 500)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(newArticle)
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
	log.Println("database connected!")
	return db
}
