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
	auth.GET("/sheet-music/:id/download/free", h.DownloadFree)
	auth.GET("/sheet-music/:id/download/paid", h.DownloadPaid)
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

	if userIDVal, ok := c.Get("user_id"); ok {
		if userID, ok := userIDVal.(uint); ok {
			if paid, err := h.orderRepo.HasPaidOrder(userID, uint(id)); err == nil {
				detail.IsPurchased = paid
			}
		}
	}

	// 详情绝不返回 PDF 直链；付费文件记录不出现在 JSON 里。
	// 预览只保留免费版切图信息；是否有文件用 has_* 标记。
	preview := make([]model.SheetFile, 0, len(detail.Files))
	for _, f := range detail.Files {
		switch f.FileType {
		case "paid":
			detail.HasPaidFile = true
		case "free":
			detail.HasFreeFile = true
			f.FilePath = ""
			f.FileName = ""
			preview = append(preview, f)
		}
	}
	detail.Files = preview

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

func (h *SheetHandler) DownloadFree(c *gin.Context) {
	userID, sheetID, sheet, ok := h.parseDownloadContext(c)
	if !ok {
		return
	}

	file, err := h.sheetFileRepo.FindBySheetIDAndType(sheetID, "free")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, response.Error(404, "该乐谱暂无免费版文件"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询文件失败"))
		return
	}
	if file.FileType != "free" {
		c.JSON(http.StatusInternalServerError, response.Error(500, "免费文件数据异常"))
		return
	}

	requiredPoints := sheet.DownloadPoints
	if requiredPoints < 0 {
		requiredPoints = 0
	}

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
		return h.dlRecordRepo.Create(record)
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	h.streamSheetFile(c, sheet.Title, "免费版", file.FilePath)
}

func (h *SheetHandler) DownloadPaid(c *gin.Context) {
	userID, sheetID, sheet, ok := h.parseDownloadContext(c)
	if !ok {
		return
	}

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
			c.JSON(http.StatusNotFound, response.Error(404, "付费文件暂未上传，请联系管理员重新上传高清版"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询文件失败"))
		return
	}
	if file.FileType != "paid" {
		c.JSON(http.StatusInternalServerError, response.Error(500, "付费文件数据异常"))
		return
	}

	if err := h.sheetFileRepo.IncrementDownload(file.ID); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新下载次数失败"))
		return
	}

	h.streamSheetFile(c, sheet.Title, "高清版", file.FilePath)
}

func (h *SheetHandler) parseDownloadContext(c *gin.Context) (userID, sheetID uint, sheet *model.SheetMusic, ok bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "请先登录"))
		return 0, 0, nil, false
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return 0, 0, nil, false
	}

	sheet, err = h.sheetService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "乐谱不存在"))
		return 0, 0, nil, false
	}

	return userIDVal.(uint), uint(id), sheet, true
}

// streamSheetFile 以附件字节流返回 PDF，不暴露静态直链。
func (h *SheetHandler) streamSheetFile(c *gin.Context, title, edition, filePath string) {
	diskPath := filepath.Join(h.cfg.Upload.Dir, strings.TrimPrefix(filePath, "/uploads/"))
	ext := strings.ToLower(filepath.Ext(filePath))
	downloadName := title
	if downloadName == "" {
		downloadName = "sheet"
	}
	if edition != "" {
		downloadName = downloadName + "_" + edition
	}
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "no-store")
	c.FileAttachment(diskPath, downloadName+ext)
}
