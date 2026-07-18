package model

import (
	"time"

	"gorm.io/gorm"
)

type ContactMessage struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name    string `gorm:"type:varchar(100);not null" json:"name"`
	Email   string `gorm:"type:varchar(100)" json:"email"`
	Subject string `gorm:"type:varchar(200)" json:"subject"`
	Content string `gorm:"type:text;not null" json:"content"`
	Read    bool   `gorm:"default:false" json:"read"`
}

type FriendLink struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	URL         string    `gorm:"type:varchar(500);not null" json:"url"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	Logo        string    `gorm:"type:varchar(500)" json:"logo"`
	IsEnabled   bool      `gorm:"default:true" json:"is_enabled"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DownloadRecord struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	SheetFileID  uint      `gorm:"index;not null" json:"sheet_file_id"`
	SheetMusicID uint      `gorm:"index" json:"sheet_music_id"`
	PointsSpent  int       `json:"points_spent"`
	CreatedAt    time.Time `json:"created_at"`
}

type SitePvHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Date      string    `gorm:"type:date;uniqueIndex;not null" json:"date"`
	PageViews int64     `gorm:"default:0" json:"page_views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ContactMessage) TableName() string  { return "contact_message" }
func (FriendLink) TableName() string      { return "friend_link" }
func (DownloadRecord) TableName() string   { return "download_record" }
func (SitePvHistory) TableName() string    { return "site_pv_history" }
