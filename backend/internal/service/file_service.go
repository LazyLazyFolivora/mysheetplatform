package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os/exec"
	"path/filepath"
	"strconv"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/repository"
)

type FileService struct {
	sheetFileRepo *repository.SheetFileRepo
	cfg           *config.Config
	logger        *zap.Logger
	db            *gorm.DB
}

type FileServiceParams struct {
	fx.In
	SheetFileRepo *repository.SheetFileRepo
	Cfg           *config.Config
	Logger        *zap.Logger
	DB            *gorm.DB
}

func NewFileService(p FileServiceParams) *FileService {
	return &FileService{
		sheetFileRepo: p.SheetFileRepo,
		cfg:           p.Cfg,
		logger:        p.Logger,
		db:            p.DB,
	}
}

func (s *FileService) UploadSheetPDF(sheetMusicID uint, fileType string, file *multipart.FileHeader) (*model.SheetFile, error) {
	if fileType != "free" && fileType != "paid" {
		return nil, errors.New("file_type 只能是 free 或 paid")
	}

	saved, err := pkg.SaveUploadedFile(file, s.cfg.Upload.Dir, pkg.FileTypeSheetPDF, s.cfg.Upload.MaxSize)
	if err != nil {
		return nil, err
	}

	// 同类型只保留最新一份，避免历史错误数据（例如全被标成 free）干扰下载
	if err := s.sheetFileRepo.SoftDeleteBySheetIDAndType(sheetMusicID, fileType); err != nil {
		return nil, fmt.Errorf("清理旧文件记录失败: %w", err)
	}

	sheetFile := &model.SheetFile{
		SheetMusicID: sheetMusicID,
		FilePath:     uploadURL(saved.RelativePath),
		FileName:     saved.FileName,
		FileType:     fileType,
		FileSize:     saved.FileSize,
		PageCount:    0,
	}

	if err := s.sheetFileRepo.Create(sheetFile); err != nil {
		return nil, fmt.Errorf("保存乐谱文件记录失败: %w", err)
	}

	s.logger.Info("sheet pdf uploaded",
		zap.Uint("sheet_music_id", sheetMusicID),
		zap.String("file_type", fileType),
		zap.String("file_name", saved.FileName),
		zap.Uint("file_id", sheetFile.ID),
	)

	// 异步把 PDF 切分成逐页图片，完成后回填 image_folder 和 page_count
	if fileType == "free" {
		go s.convertSheetPDF(sheetFile.ID, saved.FilePath)
	}

	return sheetFile, nil
}

func (s *FileService) convertSheetPDF(fileID uint, pdfPath string) {
	if _, err := exec.LookPath("pdftoppm"); err != nil {
		s.logger.Warn("未找到 pdftoppm，跳过 PDF 切图", zap.Uint("file_id", fileID))
		return
	}

	imageDir := filepath.Join(s.cfg.Upload.Dir, "sheets_images", strconv.FormatUint(uint64(fileID), 10))
	pageCount, err := pkg.ConvertPDFToImages(pdfPath, imageDir)
	if err != nil {
		s.logger.Error("PDF 切图失败", zap.Uint("file_id", fileID), zap.Error(err))
		return
	}

	updates := map[string]interface{}{
		"image_folder": fmt.Sprintf("/uploads/sheets_images/%d/", fileID),
		"page_count":   pageCount,
	}
	if err := s.db.Model(&model.SheetFile{}).Where("id = ?", fileID).Updates(updates).Error; err != nil {
		s.logger.Error("回填 PDF 切图信息失败", zap.Uint("file_id", fileID), zap.Error(err))
		return
	}
	s.logger.Info("PDF 切图完成", zap.Uint("file_id", fileID), zap.Int("page_count", pageCount))
}

func (s *FileService) UploadAudio(sheetMusicID uint, file *multipart.FileHeader) (*model.AudioFile, error) {
	saved, err := pkg.SaveUploadedFile(file, s.cfg.Upload.Dir, pkg.FileTypeAudio, s.cfg.Upload.MaxSize)
	if err != nil {
		return nil, err
	}

	audioFile := &model.AudioFile{
		SheetMusicID: sheetMusicID,
		OriginalURL:  uploadURL(saved.RelativePath),
		FileSize:     saved.FileSize,
		Quality:      "standard",
	}

	if err := s.db.Create(audioFile).Error; err != nil {
		return nil, fmt.Errorf("保存音频文件记录失败: %w", err)
	}

	// 异步转 HLS，完成后回填 hls_url；失败则前端回退播放原始文件
	go s.transcodeAudio(audioFile.ID, saved.FilePath)

	return audioFile, nil
}

func (s *FileService) transcodeAudio(audioID uint, inputPath string) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		s.logger.Warn("未找到 ffmpeg，跳过 HLS 转码", zap.Uint("audio_id", audioID))
		return
	}

	hlsDir := filepath.Join(s.cfg.Upload.Dir, "hls", strconv.FormatUint(uint64(audioID), 10))
	if err := pkg.TranscodeToHLS(inputPath, hlsDir); err != nil {
		s.logger.Error("HLS 转码失败", zap.Uint("audio_id", audioID), zap.Error(err))
		return
	}

	updates := map[string]interface{}{
		"hls_url":  fmt.Sprintf("/uploads/hls/%d/playlist.m3u8", audioID),
		"duration": pkg.ProbeDuration(inputPath),
	}
	if err := s.db.Model(&model.AudioFile{}).Where("id = ?", audioID).Updates(updates).Error; err != nil {
		s.logger.Error("回填 HLS 地址失败", zap.Uint("audio_id", audioID), zap.Error(err))
		return
	}
	s.logger.Info("HLS 转码完成", zap.Uint("audio_id", audioID))
}

func (s *FileService) UploadImage(file *multipart.FileHeader) (string, error) {
	saved, err := pkg.SaveUploadedFile(file, s.cfg.Upload.Dir, pkg.FileTypeImage, s.cfg.Upload.MaxSize)
	if err != nil {
		return "", err
	}
	return uploadURL(saved.RelativePath), nil
}

// uploadURL 将磁盘相对路径转成 /uploads 静态路由可访问的 URL 路径
func uploadURL(relativePath string) string {
	return "/uploads/" + relativePath
}

func (s *FileService) GetSheetFiles(sheetMusicID uint) ([]model.SheetFile, error) {
	file, err := s.sheetFileRepo.FindBySheetID(sheetMusicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return []model.SheetFile{*file}, nil
}

func (s *FileService) GetAudioFiles(sheetMusicID uint) ([]model.AudioFile, error) {
	var files []model.AudioFile
	err := s.db.Where("sheet_music_id = ?", sheetMusicID).Find(&files).Error
	return files, err
}
