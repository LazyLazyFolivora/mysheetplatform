package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Type string `gorm:"type:varchar(20);default:genre" json:"type"`
}

type SheetTag struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	SheetMusicID uint      `gorm:"uniqueIndex:idx_sheet_tag;not null" json:"sheet_music_id"`
	TagID        uint      `gorm:"uniqueIndex:idx_sheet_tag;not null" json:"tag_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Tag) TableName() string      { return "tag" }
func (SheetTag) TableName() string  { return "sheet_tag" }
