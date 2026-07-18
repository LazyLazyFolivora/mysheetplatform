package model

import "time"

type Permission struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Code      string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	URL       string    `gorm:"type:varchar(200)" json:"url"`
	Method    string    `gorm:"type:varchar(10)" json:"method"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PermissionModule struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	Path      string    `gorm:"type:varchar(200)" json:"path"`
	Icon      string    `gorm:"type:varchar(50)" json:"icon"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PermissionModuleRelation struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ModuleID  uint      `gorm:"index;not null" json:"module_id"`
	PermID    uint      `gorm:"index;not null" json:"perm_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Permission) TableName() string                 { return "permission" }
func (PermissionModule) TableName() string            { return "permission_module" }
func (PermissionModuleRelation) TableName() string     { return "permission_module_relation" }
