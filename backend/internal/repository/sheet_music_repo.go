package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type SheetMusicRepo struct {
	db *gorm.DB
}

func NewSheetMusicRepo(db *gorm.DB) *SheetMusicRepo {
	return &SheetMusicRepo{db: db}
}

func (r *SheetMusicRepo) FindByID(id uint) (*model.SheetMusic, error) {
	var sheet model.SheetMusic
	err := r.db.First(&sheet, id).Error
	return &sheet, err
}

func (r *SheetMusicRepo) List(page, pageSize int, keyword string, tagIDs []uint) ([]model.SheetMusic, int64, error) {
	var sheets []model.SheetMusic
	var total int64

	query := r.db.Model(&model.SheetMusic{}).Where("status = 1")

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR composer LIKE ? OR arranger LIKE ?", like, like, like)
	}

	// Count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply tag filter if specified
	if len(tagIDs) > 0 {
		subQuery := r.db.Table("sheet_tag").
			Select("sheet_music_id").
			Where("tag_id IN ?", tagIDs).
			Group("sheet_music_id").
			Having("COUNT(DISTINCT tag_id) = ?", len(tagIDs))
		query = query.Where("id IN (?)", subQuery)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&sheets).Error; err != nil {
		return nil, 0, err
	}

	return sheets, total, nil
}

func (r *SheetMusicRepo) Create(sheet *model.SheetMusic) error {
	return r.db.Create(sheet).Error
}

func (r *SheetMusicRepo) Update(sheet *model.SheetMusic) error {
	return r.db.Save(sheet).Error
}

func (r *SheetMusicRepo) Delete(id uint) error {
	return r.db.Delete(&model.SheetMusic{}, id).Error
}

func (r *SheetMusicRepo) IncrementView(id uint) error {
	return r.db.Model(&model.SheetMusic{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
