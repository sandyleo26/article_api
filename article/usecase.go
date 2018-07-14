package article

import (
	"strconv"
	"strings"
	"time"

	"github.com/sandyleo26/article_api/database"
)

//AddArticle add article
func AddArticle(userID int, articleRequest *ArticleRequest) (*Article, error) {
	newArticle := NewArticle(1, articleRequest)
	db := database.OpenDB()
	result := db.Debug().Create(newArticle)
	return newArticle, result.Error
}

//GetArticle get article by id
func GetArticle(articleID int) (*Article, error) {
	var articleFound Article
	db := database.OpenDB()
	result := db.Debug().Where(&Article{ID: articleID}).First(&articleFound)
	return &articleFound, result.Error
}

//GetTag
func GetTag(tagName string, date time.Time) (*TagResponse, error) {
	datePlusOneDay := date.Add(time.Hour * 24)
	var articlesFound []Article
	db := database.OpenDB()
	result := db.Debug().Where("created_at BETWEEN ? AND ?", date, datePlusOneDay).Order("created_at").Find(&articlesFound)
	if result.Error != nil {
		return nil, result.Error
	}

	var relatedTagsMap = make(map[string]bool)
	var filteredArticles []Article
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

	return &TagResponse{
		Tag:         tagName,
		Count:       len(filteredArticles),
		Articles:    lastTenArticles,
		RelatedTags: relatedTags,
	}, nil
}
