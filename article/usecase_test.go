package article

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sandyleo26/article_api/database"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	db := database.OpenDB()
	code := m.Run()
	db.Close()
	os.Exit(code)
}

func deleteArticle(article *Article) {
	if article == nil || article.ID == 0 {
		return
	}

	db := database.OpenDB()
	db.Debug().Delete(article)
}

func TestAddArticle(t *testing.T) {
	artcileRequest := &ArticleRequest{
		Title: "Australia won WC 2018",
		Body:  "This is fake news",
		Tags:  []string{"sports", "fake", "WC"},
	}

	newArticle, err := AddArticle(int(rand.Int31()), artcileRequest)
	defer deleteArticle(newArticle)

	assert.Nil(t, err)
	assert.NotZero(t, newArticle.ID)
	assert.Equal(t, artcileRequest.Title, newArticle.Title)
	assert.Equal(t, artcileRequest.Body, newArticle.Body)
	assert.ElementsMatch(t, artcileRequest.Tags, strings.Split(newArticle.Tags, ","))
}

func TestGetArticleSuccess(t *testing.T) {
	artcileRequest := &ArticleRequest{
		Title: "Australia won WC 2018",
		Body:  "This is fake news",
		Tags:  []string{"sports", "fake", "WC"},
	}

	newArticle, _ := AddArticle(int(rand.Int31()), artcileRequest)
	defer deleteArticle(newArticle)

	retrievedArticle, err := GetArticle(newArticle.ID)
	assert.Nil(t, err)
	assert.NotZero(t, retrievedArticle.ID)
	assert.Equal(t, artcileRequest.Title, retrievedArticle.Title)
	assert.Equal(t, artcileRequest.Body, retrievedArticle.Body)
	assert.ElementsMatch(t, artcileRequest.Tags, strings.Split(retrievedArticle.Tags, ","))
}

func TestGetArticleNotFound(t *testing.T) {
	retrievedArticle, err := GetArticle(int(rand.Int31()))
	assert.Equal(t, err.Error(), "record not found")
	assert.NotNil(t, retrievedArticle)
	assert.Zero(t, retrievedArticle.ID)
	assert.Zero(t, retrievedArticle.UserID)
	assert.Zero(t, retrievedArticle.Title)
	assert.Zero(t, retrievedArticle.Body)
	assert.Zero(t, retrievedArticle.Tags)
}

func TestGetTagSuccess(t *testing.T) {
	artcileRequest1 := &ArticleRequest{
		Title: "Australia won WC 2018",
		Body:  "This is fake news",
		Tags:  []string{"sports", "fake", "WC"},
	}
	artcileRequest2 := &ArticleRequest{
		Title: "Australia won WC 2018",
		Body:  "This is fake news",
		Tags:  []string{"man", "fake", "australia"},
	}

	newArticle1, _ := AddArticle(int(rand.Int31()), artcileRequest1)
	newArticle2, _ := AddArticle(int(rand.Int31()), artcileRequest2)
	defer deleteArticle(newArticle1)
	defer deleteArticle(newArticle2)

	date, _ := time.Parse("20060102", time.Now().Format("20060102"))
	tagResponse, _ := GetTag("fake", date)
	assert.NotNil(t, tagResponse)
	assert.Equal(t, 2, tagResponse.Count)
	assert.ElementsMatch(t, []string{strconv.Itoa(newArticle1.ID), strconv.Itoa(newArticle2.ID)}, tagResponse.Articles)
	assert.ElementsMatch(t, []string{"sports", "WC", "man", "australia"}, tagResponse.RelatedTags)
}

func TestGetTagNotFound(t *testing.T) {
	date, _ := time.Parse("20060102", time.Now().Format("20060102"))
	tagResponse, err := GetTag("fake", date)
	assert.Equal(t, err.Error(), "record not found")
	assert.NotNil(t, tagResponse)
	assert.Zero(t, tagResponse.Count)
	assert.Zero(t, tagResponse.Articles)
	assert.Zero(t, tagResponse.RelatedTags)
}
