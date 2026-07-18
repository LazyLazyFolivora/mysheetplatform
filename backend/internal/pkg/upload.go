package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileType string

const (
	FileTypeImage    FileType = "images"
	FileTypeAudio    FileType = "audio"
	FileTypeSheetPDF FileType = "sheets_pdf"
)

var allowedExtensions = map[FileType][]string{
	FileTypeImage:    {".jpg", ".jpeg", ".png", ".gif", ".webp"},
	FileTypeAudio:    {".mp3", ".wav", ".flac", ".m4a"},
	FileTypeSheetPDF: {".pdf"},
}

type SavedFile struct {
	FilePath     string `json:"file_path"`
	RelativePath string `json:"relative_path"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	Extension    string `json:"extension"`
}

func SaveUploadedFile(file *multipart.FileHeader, uploadDir string, fileType FileType, maxSizeMB int) (*SavedFile, error) {
	if maxSizeMB > 0 && file.Size > int64(maxSizeMB)<<20 {
		return nil, fmt.Errorf("文件大小超过限制（最大 %dMB）", maxSizeMB)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed, ok := allowedExtensions[fileType]
	if !ok {
		return nil, fmt.Errorf("不支持的文件类型: %s", fileType)
	}

	valid := false
	for _, e := range allowed {
		if ext == e {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("文件格式不支持，允许的格式: %s", strings.Join(allowed, ", "))
	}

	targetDir := filepath.Join(uploadDir, string(fileType))
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 时间戳加随机后缀，避免同一毫秒内的上传互相覆盖
	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		return nil, fmt.Errorf("生成文件名失败: %w", err)
	}
	newName := fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(), hex.EncodeToString(suffix), ext)
	targetPath := filepath.Join(targetDir, newName)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	relPath := filepath.Join(string(fileType), newName)

	return &SavedFile{
		FilePath:     targetPath,
		RelativePath: strings.ReplaceAll(relPath, "\\", "/"),
		FileName:     file.Filename,
		FileSize:     file.Size,
		Extension:    ext,
	}, nil
}
