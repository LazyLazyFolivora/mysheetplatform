package model

import (
	"time"

	"gorm.io/gorm"
)

type SheetMusic struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Title          string  `gorm:"type:varchar(200);not null" json:"title"`
	TrackName      string  `gorm:"type:varchar(200)" json:"track_name"`
	Composer       string  `gorm:"type:varchar(100)" json:"composer"`
	Arranger       string  `gorm:"type:varchar(100)" json:"arranger"`
	Transcriber    string  `gorm:"type:varchar(100)" json:"transcriber"`
	Description    string  `gorm:"type:text" json:"description"`
	CoverURL       string  `gorm:"type:varchar(500)" json:"cover_url"`
	ExternalLink   string  `gorm:"type:varchar(500)" json:"external_link"`
	Price          float64 `gorm:"type:numeric(10,2);default:0" json:"price"`
	DownloadPoints int     `gorm:"default:0" json:"download_points"`
	Status         int     `gorm:"default:1" json:"status"`
	ViewCount      int     `gorm:"default:0" json:"view_count"`
	UserID         uint    `gorm:"index" json:"user_id"`
}

func (SheetMusic) TableName() string {
	return "sheet_music"
}
