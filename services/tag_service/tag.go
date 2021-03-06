package tag_service

import (
	"awesomeProject/models"
	"sync"
	"time"
)

type service struct {
	sync.RWMutex
	tags map[string]map[time.Time]models.Articles
}

var db service

func init() {
	db = service{
		tags: make(map[string]map[time.Time]models.Articles),
	}
}

// AddArticle returns an error because in a real-world scenario,
// where DynamoDB is used things could go wrong and error might
// need to be returned
func AddArticle(a models.Article) error {
	db.Lock()
	defer db.Unlock()

	for _, t := range a.Tags {
		setTag(t, a)
	}

	return nil
}

func setTag(tag string, a models.Article) {
	if _, exists := db.tags[tag]; !exists {
		db.tags[tag] = make(map[time.Time]models.Articles)
	}
	setDate(tag, a)
}

func setDate(tag string, a models.Article) {
	if _, exists := db.tags[tag][a.Date]; !exists {
		db.tags[tag][a.Date] = make(models.Articles)

	}
	db.tags[tag][a.Date][a.Id] = a
}

func GetArticles(tagName string, date time.Time) models.Articles {
	db.Lock()
	defer db.Unlock()

	return db.tags[tagName][date]
}
