package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/service"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) RegisterRoutes(public, auth, admin *gin.RouterGroup) {
	admin.POST("/sheet-music/upload-pdf", h.UploadSheetPDF)
	admin.POST("/sheet-music/upload-audio", h.UploadAudio)
	admin.POST("/upload/image", h.UploadImage)
}

func (h *FileHandler) UploadSheetPDF(c *gin.Context) {
	sheetID, err := strconv.ParseUint(c.PostForm("sheet_music_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "sheet_music_id 参数错误"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请选择要上传的文件"))
		return
	}

	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, response.Error(400, "文件大小不能超过 50MB"))
		return
	}

	fileType := c.PostForm("file_type")
	if fileType == "" {
		c.JSON(http.StatusBadRequest, response.Error(400, "请指定 file_type（free 或 paid）"))
		return
	}
	if fileType != "free" && fileType != "paid" {
		c.JSON(http.StatusBadRequest, response.Error(400, "file_type 只能是 free 或 paid"))
		return
	}

	sheetFile, err := h.fileService.UploadSheetPDF(uint(sheetID), fileType, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(sheetFile))
}

func (h *FileHandler) UploadAudio(c *gin.Context) {
	sheetID, err := strconv.ParseUint(c.PostForm("sheet_music_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "sheet_music_id 参数错误"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请选择要上传的文件"))
		return
	}

	if file.Size > 200*1024*1024 {
		c.JSON(http.StatusBadRequest, response.Error(400, "音频文件大小不能超过 200MB"))
		return
	}

	audioFile, err := h.fileService.UploadAudio(uint(sheetID), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(audioFile))
}

func (h *FileHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请选择要上传的文件"))
		return
	}

	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, response.Error(400, "图片大小不能超过 10MB"))
		return
	}

	url, err := h.fileService.UploadImage(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{"url": url}))
}
