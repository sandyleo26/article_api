package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	r.HandleFunc("/articles", PostArticleHandler)
	r.HandleFunc("/articles/{id}", GetArticleHandler)
	r.HandleFunc("/tags/{tagName}/{date}", GetTagHandler)
	http.Handle("/", r)

	// start web server
	port := ":4321"
	log.Println("Starting web server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

//PostArticleHandler create article
func PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "POST") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	articleRequest := new(article.ArticleRequest)
	if err := json.NewDecoder(r.Body).Decode(articleRequest); err != nil {
		http.Error(w, "Error when parsing request", http.StatusBadRequest)
		return
	}

	newArticle := article.NewArticle(1, articleRequest)
	db := database.OpenDB()
	result := db.Debug().Create(newArticle)
	if result.Error != nil {
		http.Error(w, "Error when creating article", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newArticle)
}

//GetArticleHandler get article by id
func GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "GET") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Error when parsing parameters id", http.StatusBadRequest)
		return
	}

	var articleFound article.Article
	db := database.OpenDB()
	result := db.Debug().Where(&article.Article{ID: id}).First(&articleFound)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error when retrieving article (%d)", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articleFound)
}

//GetTagHandler get tag info given a date
func GetTagHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "GET") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	tagName := vars["tagName"]
	if len(tagName) == 0 {
		http.Error(w, fmt.Sprintf("Error when parsing parameter tagName"), http.StatusBadRequest)
		return
	}

	date, dateErr := time.Parse("20060102", vars["date"])
	if dateErr != nil {
		http.Error(w, fmt.Sprintf("Error when parsing parameter date"), http.StatusBadRequest)
		return
	}

	datePlusOneDay := date.Add(time.Hour * 24)
	var articlesFound []article.Article
	db := database.OpenDB()
	result := db.Debug().Where("created_at BETWEEN ? AND ?", date, datePlusOneDay).Order("created_at").Find(&articlesFound)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error when retrieving articles on (%v)", date), http.StatusInternalServerError)
		return
	}

	var relatedTagsMap = make(map[string]bool)
	var filteredArticles []article.Article
	for _, articleFound := range articlesFound {
		if strings.Contains(articleFound.Tags, tagName) {
			filteredArticles = append(filteredArticles, articleFound)
			tags := strings.Split(articleFound.Tags, ",")
			for _, tag := range tags {
				if !strings.EqualFold(tag, tagName) {
					relatedTagsMap[tag] = true
				}
			}
		}
	}

	var relatedTags = make([]string, 0)
	for t := range relatedTagsMap {
		relatedTags = append(relatedTags, t)
	}

	var lastTenArticles = make([]string, 0)
	for _, articleFound := range filteredArticles {
		if len(lastTenArticles) != 10 {
			lastTenArticles = append(lastTenArticles, strconv.Itoa(articleFound.ID))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&article.TagResponse{
		Tag:         tagName,
		Count:       len(filteredArticles),
		Articles:    lastTenArticles,
		RelatedTags: relatedTags,
	})
}
