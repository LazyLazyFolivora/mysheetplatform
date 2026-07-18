package service

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
)

type AdminService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAdminService(db *gorm.DB, logger *zap.Logger) *AdminService {
	return &AdminService{db: db, logger: logger}
}

type DashboardStats struct {
	UserCount      int64 `json:"user_count"`
	SheetCount     int64 `json:"sheet_count"`
	OrderCount     int64 `json:"order_count"`
	TodayPV        int64 `json:"today_pv"`
	PendingMsgs    int64 `json:"pending_messages"`
}

func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats
	_, stats.TodayPV = pkg.GetPvStats(s.db)
	s.db.Model(&model.User{}).Count(&stats.UserCount)
	s.db.Model(&model.SheetMusic{}).Count(&stats.SheetCount)
	s.db.Model(&model.SheetOrder{}).Count(&stats.OrderCount)
	s.db.Model(&model.ContactMessage{}).Where("read = ?", false).Count(&stats.PendingMsgs)
	return &stats, nil
}

func (s *AdminService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := s.db.Model(&model.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (s *AdminService) BanUser(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", 0).Error
}

func (s *AdminService) UnbanUser(id uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", 1).Error
}

func (s *AdminService) ListContactMessages(page, pageSize int) ([]model.ContactMessage, int64, error) {
	var msgs []model.ContactMessage
	var total int64
	query := s.db.Model(&model.ContactMessage{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&msgs).Error
	return msgs, total, err
}

func (s *AdminService) MarkMessageRead(id uint) error {
	return s.db.Model(&model.ContactMessage{}).Where("id = ?", id).Update("read", true).Error
}

func (s *AdminService) DeleteMessage(id uint) error {
	return s.db.Delete(&model.ContactMessage{}, id).Error
}

func (s *AdminService) CreateContactMessage(msg *model.ContactMessage) error {
	return s.db.Create(msg).Error
}

func (s *AdminService) ListFriendLinks(page, pageSize int) ([]model.FriendLink, int64, error) {
	var links []model.FriendLink
	var total int64
	query := s.db.Model(&model.FriendLink{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("sort_order ASC, id ASC").Find(&links).Error
	return links, total, err
}

// ListEnabledFriendLinks 前台展示，只返回启用的链接
func (s *AdminService) ListEnabledFriendLinks() ([]model.FriendLink, error) {
	var links []model.FriendLink
	err := s.db.Where("is_enabled = ?", true).Order("sort_order ASC, id ASC").Find(&links).Error
	return links, err
}

func (s *AdminService) GetFriendLink(id uint) (*model.FriendLink, error) {
	var link model.FriendLink
	err := s.db.First(&link, id).Error
	return &link, err
}

func (s *AdminService) CreateFriendLink(link *model.FriendLink) error {
	return s.db.Create(link).Error
}

func (s *AdminService) UpdateFriendLink(link *model.FriendLink) error {
	return s.db.Save(link).Error
}

func (s *AdminService) DeleteFriendLink(id uint) error {
	return s.db.Delete(&model.FriendLink{}, id).Error
}

func (s *AdminService) ListBannedIPs(page, pageSize int) ([]model.BannedIp, int64, error) {
	var ips []model.BannedIp
	var total int64
	query := s.db.Model(&model.BannedIp{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&ips).Error
	return ips, total, err
}

func (s *AdminService) CreateBannedIP(ip *model.BannedIp) error {
	return s.db.Create(ip).Error
}

func (s *AdminService) DeleteBannedIP(id uint) error {
	return s.db.Delete(&model.BannedIp{}, id).Error
}

func (s *AdminService) ListSystemConfigs() ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := s.db.Order("config_key ASC").Find(&configs).Error
	return configs, err
}

// UpdateSystemConfig 按键更新配置，不存在则创建（与原 Java 版 setConfigValue 一致）
func (s *AdminService) UpdateSystemConfig(key, value, remark string) error {
	var cfg model.SystemConfig
	err := s.db.Where("config_key = ?", key).First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.db.Create(&model.SystemConfig{ConfigKey: key, ConfigVal: value, Remark: remark}).Error
		}
		return err
	}
	updates := map[string]interface{}{"config_val": value}
	if remark != "" {
		updates["remark"] = remark
	}
	return s.db.Model(&cfg).Updates(updates).Error
}

func (s *AdminService) CreateSystemConfig(cfg *model.SystemConfig) error {
	var count int64
	s.db.Model(&model.SystemConfig{}).Where("config_key = ?", cfg.ConfigKey).Count(&count)
	if count > 0 {
		return errors.New("配置键已存在")
	}
	return s.db.Create(cfg).Error
}

func (s *AdminService) DeleteSystemConfig(key string) error {
	return s.db.Where("config_key = ?", key).Delete(&model.SystemConfig{}).Error
}
