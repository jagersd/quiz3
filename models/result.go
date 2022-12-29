package models

import "time"

type Result struct {
	ID         uint `json:"id" gorm:"primary_key"`
	QuizId    uint
	PlayerName string
	PlayerSlug string
    IsHost     bool `gorm:"default:false"`
	Result1    uint8 `gorm:"default:0"`
	Result2    uint8 `gorm:"default:0"`
	Result3    uint8 `gorm:"default:0"`
	Result4    uint8 `gorm:"default:0"`
	Result5    uint8 `gorm:"default:0"`
	Result6    uint8 `gorm:"default:0"`
	Result7    uint8 `gorm:"default:0"`
	Result8    uint8 `gorm:"default:0"`
	Result9    uint8 `gorm:"default:0"`
	Result10   uint8 `gorm:"default:0"`
	Result11   uint8 `gorm:"default:0"`
	Result12   uint8 `gorm:"default:0"`
	Result13   uint8 `gorm:"default:0"`
	Result14   uint8 `gorm:"default:0"`
	Result15   uint8 `gorm:"default:0"`
	Result16   uint8 `gorm:"default:0"`
	Result17   uint8 `gorm:"default:0"`
	Result18   uint8 `gorm:"default:0"`
	Result19   uint8 `gorm:"default:0"`
	Result20   uint8 `gorm:"default:0"`
	Result22   uint8 `gorm:"default:0"`
	Result23   uint8 `gorm:"default:0"`
	Result24   uint8 `gorm:"default:0"`
	Result25   uint8 `gorm:"default:0"`
	Result26   uint8 `gorm:"default:0"`
	Result27   uint8 `gorm:"default:0"`
	Result28   uint8 `gorm:"default:0"`
	Result29   uint8 `gorm:"default:0"`
	Result30   uint8 `gorm:"default:0"`
	Result31   uint8 `gorm:"default:0"`
	Result32   uint8 `gorm:"default:0"`
	Result33   uint8 `gorm:"default:0"`
	Result34   uint8 `gorm:"default:0"`
	Result35   uint8 `gorm:"default:0"`
	Result36   uint8 `gorm:"default:0"`
	Result37   uint8 `gorm:"default:0"`
	Result38   uint8 `gorm:"default:0"`
	Result39   uint8 `gorm:"default:0"`
	Result40   uint8 `gorm:"default:0"`
	Total      uint  `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
