package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

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
	transcoding   sync.Map // audioID -> struct{}，转码去重守卫
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
	// 全局只转一次：同一音频并发/重复触发时，只有一个真正执行 ffmpeg
	if _, loaded := s.transcoding.LoadOrStore(audioID, struct{}{}); loaded {
		return
	}
	defer s.transcoding.Delete(audioID)

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		s.logger.Warn("未找到 ffmpeg，跳过 HLS 转码", zap.Uint("audio_id", audioID))
		return
	}

	hlsDir := filepath.Join(s.cfg.Upload.Dir, "hls", strconv.FormatUint(uint64(audioID), 10))
	// 分片已存在（上次转码成功但回填 DB 失败 / 进程被杀），直接回填，不再重跑 ffmpeg
	if _, err := os.Stat(filepath.Join(hlsDir, "playlist.m3u8")); err != nil {
		if err := pkg.TranscodeToHLS(inputPath, hlsDir); err != nil {
			s.logger.Error("HLS 转码失败", zap.Uint("audio_id", audioID), zap.Error(err))
			return
		}
	}

	s.updateHLSRecord(audioID, inputPath)
	s.logger.Info("HLS 转码完成", zap.Uint("audio_id", audioID))
}

// updateHLSRecord 回填 hls_url 与时长
func (s *FileService) updateHLSRecord(audioID uint, inputPath string) {
	updates := map[string]interface{}{
		"hls_url":  fmt.Sprintf("/uploads/hls/%d/playlist.m3u8", audioID),
		"duration": pkg.ProbeDuration(inputPath),
	}
	if err := s.db.Model(&model.AudioFile{}).Where("id = ?", audioID).Updates(updates).Error; err != nil {
		s.logger.Error("回填 HLS 地址失败", zap.Uint("audio_id", audioID), zap.Error(err))
	}
}

// BackfillMissingHLS 为缺失 HLS 的历史音频补齐转码。
// 旧数据（从 Java 迁移或上线前上传）只有 original_url，前端回退播放原文件导致慢。
// 幂等：只处理 hls_url 为空的记录，重启后自动重试失败的项。
func (s *FileService) BackfillMissingHLS() {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		s.logger.Warn("未找到 ffmpeg，跳过 HLS 补齐")
		return
	}

	var audios []model.AudioFile
	if err := s.db.Where("hls_url = '' OR hls_url IS NULL").Find(&audios).Error; err != nil {
		s.logger.Error("查询待转码音频失败", zap.Error(err))
		return
	}
	if len(audios) == 0 {
		return
	}

	s.logger.Info("HLS 补齐开始", zap.Int("count", len(audios)))
	for _, a := range audios {
		inputPath := s.resolveUploadPath(a.OriginalURL)
		if inputPath == "" {
			s.logger.Warn("音频原文件地址无法解析，跳过",
				zap.Uint("audio_id", a.ID),
				zap.String("original_url", a.OriginalURL))
			continue
		}
		if _, err := os.Stat(inputPath); err != nil {
			s.logger.Warn("音频原文件不存在，跳过",
				zap.Uint("audio_id", a.ID),
				zap.String("path", inputPath))
			continue
		}
		s.transcodeAudio(a.ID, inputPath)
	}
	s.logger.Info("HLS 补齐完成")
}

// resolveUploadPath 把 /uploads 相对 URL 转成磁盘路径，兼容缺少前缀的旧数据
func (s *FileService) resolveUploadPath(urlPath string) string {
	if urlPath == "" {
		return ""
	}
	if strings.HasPrefix(urlPath, "/uploads/") {
		urlPath = strings.TrimPrefix(urlPath, "/uploads/")
	}
	// 拒绝绝对 HTTP 地址，避免把外部链接当成本地文件
	if strings.HasPrefix(urlPath, "http://") || strings.HasPrefix(urlPath, "https://") {
		return ""
	}
	return filepath.Join(s.cfg.Upload.Dir, urlPath)
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
