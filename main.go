package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sandyleo26/article_api/article"
	"github.com/sandyleo26/article_api/database"
)

func main() {

	// open db
	db := database.OpenDB()
	defer db.Close()

	// setup routes
	r := mux.NewRouter()
	r.HandleFunc("/articles", article.PostArticleHandler)
	r.HandleFunc("/articles/{id}", article.GetArticleHandler)
	r.HandleFunc("/tags/{tagName}/{date}", article.GetTagHandler)
	http.Handle("/", r)

	// start web server
	port := ":4321"
	log.Println("Starting web server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
