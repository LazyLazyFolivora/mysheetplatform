package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/smartwalle/alipay/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/repository"
)

type OrderService struct {
	db           *gorm.DB
	sheetRepo    *repository.SheetMusicRepo
	orderRepo    *repository.OrderRepo
	alipayClient *alipay.Client
	cfg          *config.Config
	logger       *zap.Logger
}

type OrderServiceParams struct {
	fx.In
	DB           *gorm.DB
	SheetRepo    *repository.SheetMusicRepo
	OrderRepo    *repository.OrderRepo
	AlipayClient *alipay.Client
	Cfg          *config.Config
	Logger       *zap.Logger
}

func NewOrderService(p OrderServiceParams) *OrderService {
	return &OrderService{
		db:           p.DB,
		sheetRepo:    p.SheetRepo,
		orderRepo:    p.OrderRepo,
		alipayClient: p.AlipayClient,
		cfg:          p.Cfg,
		logger:       p.Logger,
	}
}

type CreateOrderReq struct {
	UserID       uint `json:"user_id"`
	SheetMusicID uint `json:"sheet_music_id"`
}

func (s *OrderService) CreateOrder(req *CreateOrderReq) (*model.SheetOrder, string, error) {
	sheet, err := s.sheetRepo.FindByID(req.SheetMusicID)
	if err != nil {
		return nil, "", errors.New("乐谱不存在")
	}

	if sheet.Price <= 0 {
		return nil, "", errors.New("该乐谱无需购买")
	}

	if has, _ := s.orderRepo.HasPendingOrder(req.UserID, req.SheetMusicID); has {
		return nil, "", errors.New("已有待支付订单，请先完成支付")
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
		return nil, "", errors.New("创建订单失败")
	}

	amount := fmt.Sprintf("%.2f", order.Amount)
	subject := sheet.Title
	if subject == "" {
		subject = "乐谱购买"
	}
	returnURL := strings.ReplaceAll(s.cfg.Alipay.ReturnUrl, "{id}", fmt.Sprintf("%d", req.SheetMusicID))

	payURL, err := s.alipayClient.TradePagePay(alipay.TradePagePay{
		Trade: alipay.Trade{
			Subject:     subject,
			OutTradeNo:  orderNo,
			TotalAmount: amount,
			ProductCode: "FAST_INSTANT_TRADE_PAY",
			Body:        subject,
			NotifyURL:   s.cfg.Alipay.NotifyUrl,
			ReturnURL:   returnURL,
		},
	})
	if err != nil {
		s.logger.Error("build pay url failed", zap.Error(err))
		return nil, "", errors.New("生成支付链接失败")
	}

	return order, payURL.String(), nil
}

func (s *OrderService) ListOrders(userID uint, page, pageSize int) ([]model.SheetOrder, int64, error) {
	return s.orderRepo.ListByUser(userID, page, pageSize)
}

func (s *OrderService) HandleAlipayNotify(values url.Values) error {
	notification, err := s.alipayClient.DecodeNotification(context.Background(), values)
	if err != nil {
		s.logger.Warn("alipay notify verify failed",
			zap.String("out_trade_no", values.Get("out_trade_no")),
			zap.Error(err))
		return errors.New("签名验证失败")
	}

	if notification.AppId != "" && notification.AppId != s.cfg.Alipay.AppID {
		return errors.New("app_id 不匹配")
	}

	if notification.TradeStatus != alipay.TradeStatusSuccess &&
		notification.TradeStatus != alipay.TradeStatusFinished {
		return nil
	}

	order, err := s.orderRepo.GetByOrderNo(notification.OutTradeNo)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != "pending" {
		return nil
	}

	notifyAmount, err := strconv.ParseFloat(notification.TotalAmount, 64)
	if err != nil || math.Abs(notifyAmount-order.Amount) > 0.001 {
		s.logger.Warn("alipay notify amount mismatch",
			zap.String("out_trade_no", notification.OutTradeNo),
			zap.String("notify_amount", notification.TotalAmount),
			zap.Float64("order_amount", order.Amount))
		return errors.New("金额不匹配")
	}

	return s.orderRepo.MarkPaid(notification.OutTradeNo, notification.TradeNo)
}

func (s *OrderService) GetPayURL(orderNo string, userID uint) (string, error) {
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
	returnURL := strings.ReplaceAll(s.cfg.Alipay.ReturnUrl, "{id}", fmt.Sprintf("%d", order.SheetMusicID))

	payURL, err := s.alipayClient.TradePagePay(alipay.TradePagePay{
		Trade: alipay.Trade{
			Subject:     subject,
			OutTradeNo:  orderNo,
			TotalAmount: amount,
			ProductCode: "FAST_INSTANT_TRADE_PAY",
			Body:        subject,
			NotifyURL:   s.cfg.Alipay.NotifyUrl,
			ReturnURL:   returnURL,
		},
	})
	if err != nil {
		return "", errors.New("生成支付链接失败")
	}

	return payURL.String(), nil
}

func (s *OrderService) GetOrderStatus(orderNo string, userID uint) (string, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return "", errors.New("订单不存在")
	}
	if order.UserID != userID {
		return "", errors.New("订单不存在")
	}

	if order.Status != "pending" {
		return order.Status, nil
	}

	resp, err := s.alipayClient.TradeQuery(context.Background(), alipay.TradeQuery{
		OutTradeNo: orderNo,
	})
	if err != nil {
		return order.Status, nil
	}

	if resp.TradeStatus == alipay.TradeStatusSuccess ||
		resp.TradeStatus == alipay.TradeStatusFinished {
		if err := s.orderRepo.MarkPaid(orderNo, resp.TradeNo); err != nil {
			s.logger.Warn("mark paid failed in GetOrderStatus",
				zap.String("order_no", orderNo),
				zap.Error(err))
		}
		return "paid", nil
	}

	return order.Status, nil
}
