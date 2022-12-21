package models

import "time"

type Question struct {
	ID         uint `json:"id" gorm:"primary_key"`
	SubjectId  uint
	Type       uint
	Attachment string `json:"attachment" gorm:"default:null"`
	Body       string
	Answer     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
