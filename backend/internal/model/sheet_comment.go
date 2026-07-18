package model

import (
	"time"

	"gorm.io/gorm"
)

type SheetComment struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	SheetMusicID uint   `gorm:"index;not null" json:"sheet_music_id"`
	UserID       uint   `gorm:"index;not null" json:"user_id"`
	Content      string `gorm:"type:text;not null" json:"content"`
	ParentID     *uint  `gorm:"index" json:"parent_id"`
	RootID       *uint  `gorm:"index" json:"root_id"`
}

func (SheetComment) TableName() string {
	return "sheet_comment"
}
