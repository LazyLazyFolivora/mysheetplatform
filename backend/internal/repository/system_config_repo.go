package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type SystemConfigRepo struct {
	db *gorm.DB
}

func NewSystemConfigRepo(db *gorm.DB) *SystemConfigRepo {
	return &SystemConfigRepo{db: db}
}

func (r *SystemConfigRepo) List() ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := r.db.Order("config_key ASC").Find(&configs).Error
	return configs, err
}

func (r *SystemConfigRepo) GetByKey(key string) (*model.SystemConfig, error) {
	var cfg model.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&cfg).Error
	return &cfg, err
}

func (r *SystemConfigRepo) ExistsByKey(key string) (bool, error) {
	var count int64
	err := r.db.Model(&model.SystemConfig{}).Where("config_key = ?", key).Count(&count).Error
	return count > 0, err
}

func (r *SystemConfigRepo) Create(cfg *model.SystemConfig) error {
	return r.db.Create(cfg).Error
}

func (r *SystemConfigRepo) Update(key string, updates map[string]interface{}) error {
	return r.db.Model(&model.SystemConfig{}).Where("config_key = ?", key).Updates(updates).Error
}

func (r *SystemConfigRepo) Delete(key string) error {
	return r.db.Where("config_key = ?", key).Delete(&model.SystemConfig{}).Error
}
