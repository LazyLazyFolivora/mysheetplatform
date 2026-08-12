package handler

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/repository"
	"github.com/sheet-platform/backend/internal/service"
)

type SheetHandler struct {
	sheetService  *service.SheetService
	db            *gorm.DB
	userRepo      *repository.UserRepo
	sheetFileRepo *repository.SheetFileRepo
	dlRecordRepo  *repository.DownloadRecordRepo
	orderRepo     *repository.OrderRepo
	cfg           *config.Config
}

type SheetHandlerParams struct {
	fx.In
	SheetService *service.SheetService
	DB           *gorm.DB
	UserRepo     *repository.UserRepo
	FileRepo     *repository.SheetFileRepo
	DLRecordRepo *repository.DownloadRecordRepo
	OrderRepo    *repository.OrderRepo
	Cfg          *config.Config
}

func NewSheetHandler(p SheetHandlerParams) *SheetHandler {
	return &SheetHandler{
		sheetService:  p.SheetService,
		db:            p.DB,
		userRepo:      p.UserRepo,
		sheetFileRepo: p.FileRepo,
		dlRecordRepo:  p.DLRecordRepo,
		orderRepo:     p.OrderRepo,
		cfg:           p.Cfg,
	}
}

func (h *SheetHandler) RegisterRoutes(public, auth, admin *gin.RouterGroup) {
	public.GET("/sheet-music", h.List)
	public.GET("/sheet-music/:id", h.Detail)
	auth.GET("/sheet-music/:id/download", h.Download)
	public.GET("/tags", h.ListTags)
	admin.POST("/sheet-music", h.Create)
	admin.PUT("/sheet-music/:id", h.Update)
	admin.DELETE("/sheet-music/:id", h.Delete)
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

	purchased := false
	if userIDVal, ok := c.Get("user_id"); ok {
		if userID, ok := userIDVal.(uint); ok {
			paid, err := h.orderRepo.HasPaidOrder(userID, uint(id))
			if err == nil {
				purchased = paid
			}
		}
	}
	detail.IsPurchased = purchased

	// 未购买时不返回付费文件路径，避免前端误用/泄露高清版
	if !purchased && len(detail.Files) > 0 {
		filtered := make([]model.SheetFile, 0, len(detail.Files))
		for _, f := range detail.Files {
			if f.FileType != "paid" {
				filtered = append(filtered, f)
			}
		}
		detail.Files = filtered
	}

	c.JSON(http.StatusOK, response.Success(detail))
}

type createSheetReq struct {
	Title          string               `json:"title" binding:"required,max=200"`
	TrackName      string               `json:"track_name" binding:"max=200"`
	Composer       string               `json:"composer" binding:"max=100"`
	Arranger       string               `json:"arranger" binding:"max=100"`
	Transcriber    string               `json:"transcriber" binding:"max=100"`
	Description    string               `json:"description"`
	CoverURL       string               `json:"cover_url" binding:"max=500"`
	ExternalLink   string               `json:"external_link" binding:"max=500"`
	Price          float64              `json:"price" binding:"min=0"`
	DownloadPoints int                  `json:"download_points" binding:"min=0"`
	Tags           []string             `json:"tags"`
	PageSync       model.PageSyncPoints `json:"page_sync"`
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
		PageSync:       r.PageSync,
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

	// version=free|paid，默认 free；绝不能因已购买就把免费请求升级成付费文件
	version := strings.ToLower(c.DefaultQuery("version", "free"))
	if version != "free" && version != "paid" {
		c.JSON(http.StatusBadRequest, response.Error(400, "version 只能是 free 或 paid"))
		return
	}

	if version == "paid" {
		paid, err := h.orderRepo.HasPaidOrder(userID, sheetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, "查询订单失败"))
			return
		}
		if !paid {
			c.JSON(http.StatusForbidden, response.Error(403, "请先购买高清版"))
			return
		}

		file, err := h.sheetFileRepo.FindBySheetIDAndType(sheetID, "paid")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, response.Error(404, "付费文件暂未上传"))
				return
			}
			c.JSON(http.StatusInternalServerError, response.Error(500, "查询文件失败"))
			return
		}

		if err := h.sheetFileRepo.IncrementDownload(file.ID); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, "更新下载次数失败"))
			return
		}

		h.serveSheetFile(c, sheet.Title, file.FilePath)
		return
	}

	// 免费版：始终返回 free 文件，与是否已购买无关
	file, err := h.sheetFileRepo.FindBySheetIDAndType(sheetID, "free")
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

	// 已购买高清版时，免费版下载也不再扣积分
	skipPoints := false
	if requiredPoints > 0 {
		if paid, _ := h.orderRepo.HasPaidOrder(userID, sheetID); paid {
			skipPoints = true
		}
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		pointsSpent := 0
		if requiredPoints > 0 && !skipPoints {
			if err := h.userRepo.DeductPoints(userID, requiredPoints); err != nil {
				return fmt.Errorf("积分不足，需要 %d 积分", requiredPoints)
			}
			pointsSpent = requiredPoints
		}

		if err := h.sheetFileRepo.IncrementDownload(file.ID); err != nil {
			return fmt.Errorf("更新下载次数失败: %w", err)
		}

		record := &model.DownloadRecord{
			UserID:       userID,
			SheetFileID:  file.ID,
			SheetMusicID: sheetID,
			PointsSpent:  pointsSpent,
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

	h.serveSheetFile(c, sheet.Title, file.FilePath)
}

func (h *SheetHandler) serveSheetFile(c *gin.Context, title, filePath string) {
	diskPath := filepath.Join(h.cfg.Upload.Dir, strings.TrimPrefix(filePath, "/uploads/"))
	ext := strings.ToLower(filepath.Ext(filePath))
	downloadName := title
	if downloadName == "" {
		downloadName = "sheet"
	}
	c.FileAttachment(diskPath, downloadName+ext)
}
