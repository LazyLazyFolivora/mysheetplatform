package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
)

type ContactMessageRepo struct {
	db *gorm.DB
}

func NewContactMessageRepo(db *gorm.DB) *ContactMessageRepo {
	return &ContactMessageRepo{db: db}
}

func (r *ContactMessageRepo) Create(msg *model.ContactMessage) error {
	return r.db.Create(msg).Error
}

func (r *ContactMessageRepo) List(page, pageSize int) ([]model.ContactMessage, int64, error) {
	return pkg.Paginate[model.ContactMessage](r.db, page, pageSize, "created_at DESC")
}

func (r *ContactMessageRepo) CountUnread() (int64, error) {
	var count int64
	err := r.db.Model(&model.ContactMessage{}).Where("read = ?", false).Count(&count).Error
	return count, err
}

func (r *ContactMessageRepo) MarkRead(id uint) error {
	return r.db.Model(&model.ContactMessage{}).Where("id = ?", id).Update("read", true).Error
}

func (r *ContactMessageRepo) Delete(id uint) error {
	return r.db.Delete(&model.ContactMessage{}, id).Error
}
