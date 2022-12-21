package models

import "time"

type Qtype struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Description string `gorm:"unique"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
