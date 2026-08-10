package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(order *model.SheetOrder) error {
	return r.db.Create(order).Error
}

func (r *OrderRepo) GetByID(id uint) (*model.SheetOrder, error) {
	var order model.SheetOrder
	err := r.db.First(&order, id).Error
	return &order, err
}

func (r *OrderRepo) GetByOrderNo(orderNo string) (*model.SheetOrder, error) {
	var order model.SheetOrder
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (r *OrderRepo) ListByUser(userID uint, page, pageSize int) ([]model.SheetOrder, int64, error) {
	var orders []model.SheetOrder
	var total int64
	query := r.db.Model(&model.SheetOrder{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

// MarkPaid 将 pending 订单标记为已支付；WHERE 带 status 条件保证并发重复通知只生效一次
func (r *OrderRepo) MarkPaid(orderNo, tradeNo string) error {
	return r.db.Model(&model.SheetOrder{}).
		Where("order_no = ? AND status = ?", orderNo, "pending").
		Updates(map[string]interface{}{
			"status":          "paid",
			"alipay_trade_no": tradeNo,
			"paid_at":         time.Now(),
		}).Error
}

func (r *OrderRepo) HasPaidOrder(userID, sheetMusicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.SheetOrder{}).
		Where("user_id = ? AND sheet_music_id = ? AND status = ?", userID, sheetMusicID, "paid").
		Count(&count).Error
	return count > 0, err
}

func (r *OrderRepo) HasPendingOrder(userID, sheetMusicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.SheetOrder{}).
		Where("user_id = ? AND sheet_music_id = ? AND status = ?", userID, sheetMusicID, "pending").
		Count(&count).Error
	return count > 0, err
}
