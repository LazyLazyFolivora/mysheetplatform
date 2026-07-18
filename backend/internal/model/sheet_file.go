package model

import (
	"time"

	"gorm.io/gorm"
)

type SheetFile struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	SheetMusicID  uint   `gorm:"index;not null" json:"sheet_music_id"`
	FilePath      string `gorm:"type:varchar(500);not null" json:"file_path"`
	FileName      string `gorm:"type:varchar(200)" json:"file_name"`
	FileType      string `gorm:"type:varchar(10);default:free" json:"file_type"`
	ImageFolder   string `gorm:"type:varchar(500)" json:"image_folder"`
	PageCount     int    `gorm:"default:0" json:"page_count"`
	FileSize      int64  `json:"file_size"`
	Points        int    `gorm:"default:0" json:"points"`
	DownloadCount int    `gorm:"default:0" json:"download_count"`
}

func (SheetFile) TableName() string {
	return "sheet_file"
}
