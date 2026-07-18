package model

import (
	"time"

	"gorm.io/gorm"
)

type SystemConfig struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ConfigKey string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"`
	ConfigVal string    `gorm:"type:text" json:"config_val"`
	Remark    string    `gorm:"type:varchar(200)" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SystemLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Action  string `gorm:"type:varchar(100);not null" json:"action"`
	Detail  string `gorm:"type:text" json:"detail"`
	UserID  uint   `gorm:"index" json:"user_id"`
	IP      string `gorm:"type:varchar(45)" json:"ip"`
}

type BannedIp struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	IP        string    `gorm:"type:varchar(45);uniqueIndex;not null" json:"ip"`
	Reason    string    `gorm:"type:varchar(200)" json:"reason"`
	BannedBy  uint      `json:"banned_by"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserLoginLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	UserAgent string    `gorm:"type:varchar(500)" json:"user_agent"`
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
}

func (SystemConfig) TableName() string  { return "system_config" }
func (SystemLog) TableName() string     { return "system_log" }
func (BannedIp) TableName() string     { return "banned_ip" }
func (UserLoginLog) TableName() string  { return "user_login_log" }
