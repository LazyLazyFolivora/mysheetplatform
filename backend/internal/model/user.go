package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username  string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string  `gorm:"type:varchar(255);not null" json:"-"`
	// email 为空时存 NULL 而非 ''，避免唯一索引冲突（Postgres 允许多个 NULL）
	Email     *string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Nickname  string  `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string  `gorm:"type:varchar(255)" json:"avatar"`
	Points    int     `gorm:"default:0" json:"points"`
	IsAdmin   bool    `gorm:"default:false" json:"is_admin"`
	Status    int     `gorm:"default:1" json:"status"`
}

func (User) TableName() string {
	return "user"
}
