package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//PostArticleHandler create article
func PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "POST") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	articleRequest := new(ArticleRequest)
	if err := json.NewDecoder(r.Body).Decode(articleRequest); err != nil {
		http.Error(w, "Error when parsing request", http.StatusBadRequest)
		return
	}

	newArticle, newArticleErr := AddArticle(1, articleRequest)
	if newArticleErr != nil {
		http.Error(w, "Error when creating article", http.StatusInternalServerError)
		return
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

	articleFound, articleFoundErr := GetArticle(id)
	if articleFoundErr != nil {
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

	tagResponse, tagResponseErr := GetTag(tagName, date)
	if tagResponseErr != nil {
		http.Error(w, fmt.Sprintf("Error when retrieving articles on (%v)", date), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tagResponse)
}
