package models

import "time"

type Option struct {
	ID         uint `json:"id" gorm:"primary_key"`
	QuestionId uint
	Option     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
