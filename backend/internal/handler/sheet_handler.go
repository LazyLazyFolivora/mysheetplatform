package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/repository"
	"github.com/sheet-platform/backend/internal/service"
)

type SheetHandler struct {
	sheetService   *service.SheetService
	db             *gorm.DB
	userRepo       *repository.UserRepo
	sheetFileRepo  *repository.SheetFileRepo
	dlRecordRepo   *repository.DownloadRecordRepo
	cfg            *config.Config
}

func NewSheetHandler(
	sheetService *service.SheetService,
	db *gorm.DB,
	userRepo *repository.UserRepo,
	sheetFileRepo *repository.SheetFileRepo,
	dlRecordRepo *repository.DownloadRecordRepo,
	cfg *config.Config,
) *SheetHandler {
	return &SheetHandler{
		sheetService:  sheetService,
		db:            db,
		userRepo:      userRepo,
		sheetFileRepo: sheetFileRepo,
		dlRecordRepo:  dlRecordRepo,
		cfg:           cfg,
	}
}

func (h *SheetHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var tagIDs []uint
	for _, raw := range c.QueryArray("tag_ids") {
		id, err := strconv.ParseUint(raw, 10, 32)
		if err == nil {
			tagIDs = append(tagIDs, uint(id))
		}
	}

	req := &service.SheetListReq{
		Page:     page,
		PageSize: pageSize,
		Keyword:  keyword,
		TagIDs:   tagIDs,
	}

	resp, err := h.sheetService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))
}

func (h *SheetHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	detail, err := h.sheetService.Detail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(detail))
}

type createSheetReq struct {
	Title          string   `json:"title" binding:"required,max=200"`
	TrackName      string   `json:"track_name" binding:"max=200"`
	Composer       string   `json:"composer" binding:"max=100"`
	Arranger       string   `json:"arranger" binding:"max=100"`
	Transcriber    string   `json:"transcriber" binding:"max=100"`
	Description    string   `json:"description"`
	CoverURL       string   `json:"cover_url" binding:"max=500"`
	ExternalLink   string   `json:"external_link" binding:"max=500"`
	Price          float64  `json:"price" binding:"min=0"`
	DownloadPoints int      `json:"download_points" binding:"min=0"`
	Tags           []string `json:"tags"`
}

func (r *createSheetReq) toModel() *model.SheetMusic {
	return &model.SheetMusic{
		Title:          r.Title,
		TrackName:      r.TrackName,
		Composer:       r.Composer,
		Arranger:       r.Arranger,
		Transcriber:    r.Transcriber,
		Description:    r.Description,
		CoverURL:       r.CoverURL,
		ExternalLink:   r.ExternalLink,
		Price:          r.Price,
		DownloadPoints: r.DownloadPoints,
	}
}

func (h *SheetHandler) Create(c *gin.Context) {
	var req createSheetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	userID, _ := c.Get("user_id")

	sheet := req.toModel()
	sheet.Status = 1
	sheet.UserID = userID.(uint)

	if err := h.sheetService.Create(sheet, req.Tags); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(sheet))
}

func (h *SheetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	var req createSheetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	sheet := req.toModel()
	sheet.ID = uint(id)

	if err := h.sheetService.Update(sheet, req.Tags); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(sheet))
}

func (h *SheetHandler) ListTags(c *gin.Context) {
	tags, err := h.sheetService.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取标签失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(tags))
}

func (h *SheetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	if err := h.sheetService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *SheetHandler) Download(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "请先登录"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	sheetID := uint(id)

	userID := userIDVal.(uint)

	sheet, err := h.sheetService.FindByID(sheetID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "乐谱不存在"))
		return
	}

	file, err := h.sheetFileRepo.FindBySheetID(sheetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, response.Error(404, "该乐谱暂无可下载文件"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询文件失败"))
		return
	}

	requiredPoints := sheet.DownloadPoints
	if requiredPoints < 0 {
		requiredPoints = 0
	}

	// 积分下载：事务里扣积分、增下载次数、记下载记录
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if requiredPoints > 0 {
			if err := h.userRepo.DeductPoints(userID, requiredPoints); err != nil {
				return fmt.Errorf("积分不足，需要 %d 积分", requiredPoints)
			}
		}

		if err := h.sheetFileRepo.IncrementDownload(file.ID); err != nil {
			return fmt.Errorf("更新下载次数失败: %w", err)
		}

		record := &model.DownloadRecord{
			UserID:       userID,
			SheetFileID:  file.ID,
			SheetMusicID: sheetID,
			PointsSpent:  requiredPoints,
		}
		if err := h.dlRecordRepo.Create(record); err != nil {
			return fmt.Errorf("记录下载失败: %w", err)
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	// 把 /uploads/sheets_pdf/xxx.pdf 转成磁盘绝对路径，直接返回文件
	diskPath := filepath.Join(h.cfg.Upload.Dir, strings.TrimPrefix(file.FilePath, "/uploads/"))
	ext := strings.ToLower(filepath.Ext(file.FilePath))
	contentType := "application/octet-stream"
	switch ext {
	case ".pdf":
		contentType = "application/pdf"
	}

	downloadName := sheet.Title
	if downloadName == "" {
		downloadName = "sheet"
	}
	c.Header("Content-Type", contentType)
	encodedName := url.PathEscape(downloadName)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s%s"; filename*=UTF-8''%s%s`, encodedName, ext, encodedName, ext))
	c.File(diskPath)
}
