package repository

import (
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type SheetFileRepo struct {
	db *gorm.DB
}

func NewSheetFileRepo(db *gorm.DB) *SheetFileRepo {
	return &SheetFileRepo{db: db}
}

func (r *SheetFileRepo) FindBySheetID(sheetMusicID uint) (*model.SheetFile, error) {
	var file model.SheetFile
	err := r.db.Where("sheet_music_id = ?", sheetMusicID).First(&file).Error
	return &file, err
}

func (r *SheetFileRepo) FindBySheetIDAndType(sheetMusicID uint, fileType string) (*model.SheetFile, error) {
	var file model.SheetFile
	err := r.db.Where("sheet_music_id = ? AND file_type = ?", sheetMusicID, fileType).First(&file).Error
	return &file, err
}

func (r *SheetFileRepo) ListBySheetID(sheetMusicID uint) ([]model.SheetFile, error) {
	var files []model.SheetFile
	err := r.db.Where("sheet_music_id = ?", sheetMusicID).Order("created_at ASC").Find(&files).Error
	return files, err
}

func (r *SheetFileRepo) ListAudioBySheetID(sheetMusicID uint) ([]model.AudioFile, error) {
	var files []model.AudioFile
	err := r.db.Where("sheet_music_id = ?", sheetMusicID).Order("created_at ASC").Find(&files).Error
	return files, err
}

func (r *SheetFileRepo) Create(file *model.SheetFile) error {
	return r.db.Create(file).Error
}

func (r *SheetFileRepo) IncrementDownload(id uint) error {
	return r.db.Model(&model.SheetFile{}).Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}
