package model

import "time"

type SheetLike struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	SheetMusicID uint      `gorm:"uniqueIndex:idx_user_sheet;not null" json:"sheet_music_id"`
	UserID       uint      `gorm:"uniqueIndex:idx_user_sheet;not null" json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (SheetLike) TableName() string {
	return "sheet_like"
}
