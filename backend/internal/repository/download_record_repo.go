package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type DownloadRecordRepo struct {
	db *gorm.DB
}

func NewDownloadRecordRepo(db *gorm.DB) *DownloadRecordRepo {
	return &DownloadRecordRepo{db: db}
}

func (r *DownloadRecordRepo) Create(record *model.DownloadRecord) error {
	return r.db.Create(record).Error
}
