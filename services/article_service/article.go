package article_service

import (
	"awesomeProject/models"
	"errors"
	"sync"
)

type service struct {
	sync.RWMutex
	articles map[int]models.Article
}

var db service

func init() {
	db = service{
		articles: make(map[int]models.Article),
	}
}

func AddArticle(a *models.Article) error {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.articles[a.Id]; ok {
		return errors.New("article with that ID already exists")
	}

	db.articles[a.Id] = *a
	return nil
}

func Get(id int) (*models.Article, error) {
	db.Lock()
	defer db.Unlock()

	a, ok := db.articles[id]
	if !ok {
		return nil, errors.New("no article with that ID exists")
	}

	return &a, nil
}
