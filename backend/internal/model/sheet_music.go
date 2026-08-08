package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PageSyncPoint 曲谱同步点：播放到 Time 秒时翻到第 Page 页
type PageSyncPoint struct {
	Time float64 `json:"time"`
	Page int     `json:"page"`
}

// PageSyncPoints 以 JSON 文本存储在数据库，序列化为 JSON 数组返回给前端
type PageSyncPoints []PageSyncPoint

func (p PageSyncPoints) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}
	return json.Marshal(p)
}

func (p *PageSyncPoints) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("无法解析 page_sync 字段: %T", value)
	}
	if len(data) == 0 {
		*p = nil
		return nil
	}
	return json.Unmarshal(data, p)
}

type SheetMusic struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

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

	// 曲谱同步点（可选）：为空表示该乐谱不启用播放翻页同步
	PageSync PageSyncPoints `gorm:"type:jsonb" json:"page_sync,omitempty"`
}

func (SheetMusic) TableName() string {
	return "sheet_music"
}
