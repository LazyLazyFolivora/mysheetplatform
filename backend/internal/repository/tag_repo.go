package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) FindBySheetID(sheetMusicID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Table("tag").
		Joins("JOIN sheet_tag ON sheet_tag.tag_id = tag.id").
		Where("sheet_tag.sheet_music_id = ?", sheetMusicID).
		Find(&tags).Error
	return tags, err
}

func (r *TagRepo) ListByType(tagType string) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Where("type = ?", tagType).Find(&tags).Error
	return tags, err
}

func (r *TagRepo) ListAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *TagRepo) FindOrCreate(name, tagType string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ? AND type = ?", name, tagType).First(&tag).Error
	if err == nil {
		return &tag, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	tag = model.Tag{Name: name, Type: tagType}
	err = r.db.Create(&tag).Error
	return &tag, err
}

// SetSheetTags 全量替换乐谱的标签关联
func (r *TagRepo) SetSheetTags(sheetMusicID uint, tagIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("sheet_music_id = ?", sheetMusicID).Delete(&model.SheetTag{}).Error; err != nil {
			return err
		}
		for _, tagID := range tagIDs {
			if err := tx.Create(&model.SheetTag{SheetMusicID: sheetMusicID, TagID: tagID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// MapBySheetIDs 批量查询多个乐谱的标签，按乐谱 ID 分组返回
func (r *TagRepo) MapBySheetIDs(sheetIDs []uint) (map[uint][]model.Tag, error) {
	result := make(map[uint][]model.Tag)
	if len(sheetIDs) == 0 {
		return result, nil
	}

	type sheetTagRow struct {
		SheetMusicID uint   `gorm:"column:sheet_music_id"`
		ID           uint   `gorm:"column:id"`
		Name         string `gorm:"column:name"`
		Type         string `gorm:"column:type"`
	}
	var rows []sheetTagRow
	err := r.db.Table("tag").
		Select("tag.id, tag.name, tag.type, sheet_tag.sheet_music_id").
		Joins("JOIN sheet_tag ON sheet_tag.tag_id = tag.id").
		Where("sheet_tag.sheet_music_id IN ? AND tag.deleted_at IS NULL", sheetIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.SheetMusicID] = append(result[row.SheetMusicID], model.Tag{
			ID:   row.ID,
			Name: row.Name,
			Type: row.Type,
		})
	}
	return result, nil
}
