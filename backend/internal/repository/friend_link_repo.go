package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
)

type FriendLinkRepo struct {
	db *gorm.DB
}

func NewFriendLinkRepo(db *gorm.DB) *FriendLinkRepo {
	return &FriendLinkRepo{db: db}
}

func (r *FriendLinkRepo) Create(link *model.FriendLink) error {
	return r.db.Create(link).Error
}

func (r *FriendLinkRepo) GetByID(id uint) (*model.FriendLink, error) {
	var link model.FriendLink
	err := r.db.First(&link, id).Error
	return &link, err
}

func (r *FriendLinkRepo) List(page, pageSize int) ([]model.FriendLink, int64, error) {
	return pkg.Paginate[model.FriendLink](r.db, page, pageSize, "sort_order ASC, id ASC")
}

func (r *FriendLinkRepo) ListEnabled() ([]model.FriendLink, error) {
	var links []model.FriendLink
	err := r.db.Where("is_enabled = ?", true).Order("sort_order ASC, id ASC").Find(&links).Error
	return links, err
}

func (r *FriendLinkRepo) Update(link *model.FriendLink) error {
	return r.db.Save(link).Error
}

func (r *FriendLinkRepo) Delete(id uint) error {
	return r.db.Delete(&model.FriendLink{}, id).Error
}
