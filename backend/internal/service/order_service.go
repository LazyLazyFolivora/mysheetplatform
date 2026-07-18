package service

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/repository"
)

type OrderService struct {
	db        *gorm.DB
	sheetRepo *repository.SheetMusicRepo
	orderRepo *repository.OrderRepo
	cfg       *config.Config
	logger    *zap.Logger
}

func NewOrderService(db *gorm.DB, sheetRepo *repository.SheetMusicRepo, orderRepo *repository.OrderRepo, cfg *config.Config, logger *zap.Logger) *OrderService {
	return &OrderService{
		db:        db,
		sheetRepo: sheetRepo,
		orderRepo: orderRepo,
		cfg:       cfg,
		logger:    logger,
	}
}

type CreateOrderReq struct {
	UserID       uint `json:"user_id"`
	SheetMusicID uint `json:"sheet_music_id"`
}

func (s *OrderService) CreateOrder(req *CreateOrderReq) (*model.SheetOrder, error) {
	sheet, err := s.sheetRepo.FindByID(req.SheetMusicID)
	if err != nil {
		return nil, errors.New("乐谱不存在")
	}

	if sheet.Price <= 0 {
		return nil, errors.New("该乐谱无需购买")
	}

	orderNo := fmt.Sprintf("SHEET%d%s", req.SheetMusicID, time.Now().Format("20060102150405"))

	order := &model.SheetOrder{
		OrderNo:      orderNo,
		UserID:       req.UserID,
		SheetMusicID: req.SheetMusicID,
		Amount:       sheet.Price,
		Status:       "pending",
	}

	if err := s.orderRepo.Create(order); err != nil {
		s.logger.Error("create order failed", zap.Error(err))
		return nil, errors.New("创建订单失败")
	}

	return order, nil
}

func (s *OrderService) ListOrders(userID uint, page, pageSize int) ([]model.SheetOrder, int64, error) {
	return s.orderRepo.ListByUser(userID, page, pageSize)
}

func (s *OrderService) HandleAlipayNotify(orderNo, tradeNo, status string) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != "pending" {
		return nil
	}
	return s.orderRepo.UpdateStatus(orderNo, status)
}

func (s *OrderService) GetPayForm(orderNo string, userID uint) (string, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return "", errors.New("订单不存在")
	}
	if order.UserID != userID {
		return "", errors.New("无权操作该订单")
	}
	if order.Status != "pending" {
		return "", errors.New("订单状态不允许支付")
	}

	sheet, err := s.sheetRepo.FindByID(order.SheetMusicID)
	if err != nil {
		return "", errors.New("乐谱不存在")
	}

	amount := fmt.Sprintf("%.2f", order.Amount)
	subject := sheet.Title
	if subject == "" {
		subject = "乐谱购买"
	}

	return pkg.BuildPayForm(&s.cfg.Alipay, orderNo, amount, subject, subject, fmt.Sprintf("%d", order.SheetMusicID))
}

func (s *OrderService) GetOrderStatus(orderNo string) (string, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return "", errors.New("订单不存在")
	}
	return order.Status, nil
}
