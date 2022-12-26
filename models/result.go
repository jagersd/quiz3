package models

import "time"

type Result struct {
	ID         uint `json:"id" gorm:"primary_key"`
	QuizId    uint
	PlayerName string
	PlayerSlug string
    IsHost     bool `gorm:"default:false"`
	Result1    uint8 `gorm:"default:null"`
	Result2    uint8 `gorm:"default:null"`
	Result3    uint8 `gorm:"default:null"`
	Result4    uint8 `gorm:"default:null"`
	Result5    uint8 `gorm:"default:null"`
	Result6    uint8 `gorm:"default:null"`
	Result7    uint8 `gorm:"default:null"`
	Result8    uint8 `gorm:"default:null"`
	Result9    uint8 `gorm:"default:null"`
	Result10   uint8 `gorm:"default:null"`
	Result11   uint8 `gorm:"default:null"`
	Result12   uint8 `gorm:"default:null"`
	Result13   uint8 `gorm:"default:null"`
	Result14   uint8 `gorm:"default:null"`
	Result15   uint8 `gorm:"default:null"`
	Result16   uint8 `gorm:"default:null"`
	Result17   uint8 `gorm:"default:null"`
	Result18   uint8 `gorm:"default:null"`
	Result19   uint8 `gorm:"default:null"`
	Result20   uint8 `gorm:"default:null"`
	Result22   uint8 `gorm:"default:null"`
	Result23   uint8 `gorm:"default:null"`
	Result24   uint8 `gorm:"default:null"`
	Result25   uint8 `gorm:"default:null"`
	Result26   uint8 `gorm:"default:null"`
	Result27   uint8 `gorm:"default:null"`
	Result28   uint8 `gorm:"default:null"`
	Result29   uint8 `gorm:"default:null"`
	Result30   uint8 `gorm:"default:null"`
	Result31   uint8 `gorm:"default:null"`
	Result32   uint8 `gorm:"default:null"`
	Result33   uint8 `gorm:"default:null"`
	Result34   uint8 `gorm:"default:null"`
	Result35   uint8 `gorm:"default:null"`
	Result36   uint8 `gorm:"default:null"`
	Result37   uint8 `gorm:"default:null"`
	Result38   uint8 `gorm:"default:null"`
	Result39   uint8 `gorm:"default:null"`
	Result40   uint8 `gorm:"default:null"`
	Total      uint  `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
