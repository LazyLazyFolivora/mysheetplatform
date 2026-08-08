package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"go.uber.org/fx"
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

type OrderServiceParams struct {
	fx.In
	DB        *gorm.DB
	SheetRepo *repository.SheetMusicRepo
	OrderRepo *repository.OrderRepo
	Cfg       *config.Config
	Logger    *zap.Logger
}

func NewOrderService(p OrderServiceParams) *OrderService {
	return &OrderService{
		db:        p.DB,
		sheetRepo: p.SheetRepo,
		orderRepo: p.OrderRepo,
		cfg:       p.Cfg,
		logger:    p.Logger,
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

// HandleAlipayNotify 处理支付宝异步通知：验签 → 校验商户/金额 → 幂等更新订单状态。
// params 是通知里的全部表单参数（含 sign）。
func (s *OrderService) HandleAlipayNotify(params map[string]string) error {
	if err := pkg.VerifyAlipayNotify(params, s.cfg.Alipay.AlipayPublicKey); err != nil {
		s.logger.Warn("alipay notify verify failed",
			zap.String("out_trade_no", params["out_trade_no"]),
			zap.Error(err))
		return errors.New("签名验证失败")
	}

	if appID := params["app_id"]; appID != "" && appID != s.cfg.Alipay.AppID {
		return errors.New("app_id 不匹配")
	}

	tradeStatus := params["trade_status"]
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		// 其他状态（WAIT_BUYER_PAY / TRADE_CLOSED）直接确认，避免支付宝重试
		return nil
	}

	orderNo := params["out_trade_no"]
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != "pending" {
		// 已处理过，幂等返回成功
		return nil
	}

	notifyAmount, err := strconv.ParseFloat(params["total_amount"], 64)
	if err != nil || math.Abs(notifyAmount-order.Amount) > 0.001 {
		s.logger.Warn("alipay notify amount mismatch",
			zap.String("out_trade_no", orderNo),
			zap.String("notify_amount", params["total_amount"]),
			zap.Float64("order_amount", order.Amount))
		return errors.New("金额不匹配")
	}

	return s.orderRepo.MarkPaid(orderNo, params["trade_no"])
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
