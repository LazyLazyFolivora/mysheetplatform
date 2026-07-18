package model

import (
	"time"

	"gorm.io/gorm"
)

type AudioFile struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	SheetMusicID uint   `gorm:"index;not null" json:"sheet_music_id"`
	OriginalURL  string `gorm:"type:varchar(500)" json:"original_url"`
	HlsURL       string `gorm:"type:varchar(500)" json:"hls_url"`
	Quality      string `gorm:"type:varchar(20);default:standard" json:"quality"`
	Duration     int    `json:"duration"`
	FileSize     int64  `json:"file_size"`
}

func (AudioFile) TableName() string {
	return "audio_file"
}
