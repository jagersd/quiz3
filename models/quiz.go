package models

import "time"

type Quiz struct {
	ID        uint `json:"id" gorm:"primary_key"`
	QuizSlug  string
	Started   bool `gorm:"default:false"`
	Finished  bool `gorm:"default:false"`
    Questions string 
	CreatedAt time.Time
	UpdatedAt time.Time
}
