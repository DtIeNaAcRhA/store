package model

import "time"

type Item struct {
	ID              int
	UserID          int
	Title           string
	DescriptionItem string
	ImagePath       string
	Price           float64
	CreatedAt       time.Time
	AuthorLogin     string
}
