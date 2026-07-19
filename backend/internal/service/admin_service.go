package service

import (
	"errors"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/repository"
)

type AdminService struct {
	db              *gorm.DB
	logger          *zap.Logger
	userRepo        *repository.UserRepo
	contactMsgRepo  *repository.ContactMessageRepo
	friendLinkRepo  *repository.FriendLinkRepo
	bannedIPRepo    *repository.BannedIPRepo
	sysConfigRepo   *repository.SystemConfigRepo
}

type AdminServiceParams struct {
	fx.In
	DB             *gorm.DB
	Logger         *zap.Logger
	UserRepo       *repository.UserRepo
	ContactMsgRepo *repository.ContactMessageRepo
	FriendLinkRepo *repository.FriendLinkRepo
	BannedIPRepo   *repository.BannedIPRepo
	SysConfigRepo  *repository.SystemConfigRepo
}

func NewAdminService(p AdminServiceParams) *AdminService {
	return &AdminService{
		db:              p.DB,
		logger:          p.Logger,
		userRepo:        p.UserRepo,
		contactMsgRepo:  p.ContactMsgRepo,
		friendLinkRepo:  p.FriendLinkRepo,
		bannedIPRepo:    p.BannedIPRepo,
		sysConfigRepo:   p.SysConfigRepo,
	}
}

// ===== Dashboard =====

type DashboardStats struct {
	UserCount   int64 `json:"user_count"`
	SheetCount  int64 `json:"sheet_count"`
	OrderCount  int64 `json:"order_count"`
	TodayPV     int64 `json:"today_pv"`
	PendingMsgs int64 `json:"pending_messages"`
}

func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats
	_, stats.TodayPV = pkg.GetPvStats(s.db)
	s.db.Model(&model.User{}).Count(&stats.UserCount)
	s.db.Model(&model.SheetMusic{}).Count(&stats.SheetCount)
	s.db.Model(&model.SheetOrder{}).Count(&stats.OrderCount)
	stats.PendingMsgs, _ = s.contactMsgRepo.CountUnread()
	return &stats, nil
}

// ===== Users =====

func (s *AdminService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

func (s *AdminService) BanUser(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", 0).Error
}

func (s *AdminService) UnbanUser(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", 1).Error
}

// ===== Contact Messages =====

func (s *AdminService) ListContactMessages(page, pageSize int) ([]model.ContactMessage, int64, error) {
	return s.contactMsgRepo.List(page, pageSize)
}

func (s *AdminService) MarkMessageRead(id uint) error {
	return s.contactMsgRepo.MarkRead(id)
}

func (s *AdminService) DeleteMessage(id uint) error {
	return s.contactMsgRepo.Delete(id)
}

func (s *AdminService) CreateContactMessage(msg *model.ContactMessage) error {
	return s.contactMsgRepo.Create(msg)
}

// ===== Friend Links =====

type CreateFriendLinkReq struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	URL         string `json:"url" binding:"required,url,max=500"`
	Description string `json:"description" binding:"max=500"`
	Logo        string `json:"logo" binding:"max=500"`
	IsEnabled   *bool  `json:"is_enabled"`
	SortOrder   int    `json:"sort_order" binding:"min=0"`
}

func (s *AdminService) CreateFriendLink(req CreateFriendLinkReq) (*model.FriendLink, error) {
	enabled := true
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}
	link := &model.FriendLink{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Logo:        req.Logo,
		IsEnabled:   enabled,
		SortOrder:   req.SortOrder,
	}
	if err := s.friendLinkRepo.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}

type UpdateFriendLinkReq struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	URL         string `json:"url" binding:"required,url,max=500"`
	Description string `json:"description" binding:"max=500"`
	Logo        string `json:"logo" binding:"max=500"`
	IsEnabled   *bool  `json:"is_enabled"`
	SortOrder   int    `json:"sort_order" binding:"min=0"`
}

func (s *AdminService) UpdateFriendLink(id uint, req UpdateFriendLinkReq) error {
	link, err := s.friendLinkRepo.GetByID(id)
	if err != nil {
		return err
	}
	link.Name = req.Name
	link.URL = req.URL
	link.Description = req.Description
	link.Logo = req.Logo
	if req.IsEnabled != nil {
		link.IsEnabled = *req.IsEnabled
	}
	link.SortOrder = req.SortOrder
	return s.friendLinkRepo.Update(link)
}

func (s *AdminService) DeleteFriendLink(id uint) error {
	return s.friendLinkRepo.Delete(id)
}

func (s *AdminService) ListFriendLinks(page, pageSize int) ([]model.FriendLink, int64, error) {
	return s.friendLinkRepo.List(page, pageSize)
}

func (s *AdminService) ListEnabledFriendLinks() ([]model.FriendLink, error) {
	return s.friendLinkRepo.ListEnabled()
}

// ===== Banned IPs =====

func (s *AdminService) ListBannedIPs(page, pageSize int) ([]model.BannedIp, int64, error) {
	return s.bannedIPRepo.List(page, pageSize)
}

func (s *AdminService) CreateBannedIP(ip *model.BannedIp) error {
	return s.bannedIPRepo.Create(ip)
}

func (s *AdminService) DeleteBannedIP(id uint) error {
	return s.bannedIPRepo.Delete(id)
}

// ===== System Configs =====

func (s *AdminService) ListSystemConfigs() ([]model.SystemConfig, error) {
	return s.sysConfigRepo.List()
}

func (s *AdminService) UpdateSystemConfig(key, value, remark string) error {
	_, err := s.sysConfigRepo.GetByKey(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.sysConfigRepo.Create(&model.SystemConfig{ConfigKey: key, ConfigVal: value, Remark: remark})
		}
		return err
	}
	updates := map[string]interface{}{"config_val": value}
	if remark != "" {
		updates["remark"] = remark
	}
	return s.sysConfigRepo.Update(key, updates)
}

func (s *AdminService) CreateSystemConfig(cfg *model.SystemConfig) error {
	exists, err := s.sysConfigRepo.ExistsByKey(cfg.ConfigKey)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("配置键已存在")
	}
	return s.sysConfigRepo.Create(cfg)
}

func (s *AdminService) DeleteSystemConfig(key string) error {
	return s.sysConfigRepo.Delete(key)
}
