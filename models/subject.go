package models

import "time"

type Subject struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `gorm:"unique"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
