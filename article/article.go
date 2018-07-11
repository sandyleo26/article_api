package article

import (
	"strings"
	"time"
)

//Article type definition
type Article struct {
	ID        int        `gorm:"primary_key"`
	CreatedAt time.Time  `json:"date"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	UserID    uint       `json:"-"`
	Title     string
	Body      string
	Tags      string
}

type ArticleRequest struct {
	Title string
	Body  string
	Tags  []string
}

//NewArticle create new article from request
func NewArticle(userID uint, request *ArticleRequest) *Article {
	return &Article{
		UserID: userID,
		Title:  request.Title,
		Body:   request.Body,
		Tags:   strings.Join(request.Tags[:], ","),
	}
}

func (*Article) TableName() string {
	return "article"
}
