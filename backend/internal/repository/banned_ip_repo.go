package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
)

type BannedIPRepo struct {
	db *gorm.DB
}

func NewBannedIPRepo(db *gorm.DB) *BannedIPRepo {
	return &BannedIPRepo{db: db}
}

func (r *BannedIPRepo) Create(ip *model.BannedIp) error {
	return r.db.Create(ip).Error
}

func (r *BannedIPRepo) List(page, pageSize int) ([]model.BannedIp, int64, error) {
	return pkg.Paginate[model.BannedIp](r.db, page, pageSize, "created_at DESC")
}

func (r *BannedIPRepo) Delete(id uint) error {
	return r.db.Delete(&model.BannedIp{}, id).Error
}
