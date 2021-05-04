package models

import "time"

type Article struct {
	Id    int       `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Body  string    `json:"body"`
	Tags  []string  `json:"tags"`
}

type Articles map[int]Article
