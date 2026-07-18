package model

import "time"

type Role struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(200)" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRole struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_role;not null" json:"user_id"`
	RoleID    uint      `gorm:"uniqueIndex:idx_user_role;not null" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RoleModule struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	RoleID    uint      `gorm:"index;not null" json:"role_id"`
	ModuleID  uint      `gorm:"index;not null" json:"module_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Role) TableName() string     { return "role" }
func (UserRole) TableName() string  { return "user_role" }
func (RoleModule) TableName() string { return "role_module" }
