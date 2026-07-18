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

func (r *OrderRepo) UpdateStatus(orderNo, status string) error {
	return r.db.Model(&model.SheetOrder{}).
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{
			"status":  status,
			"paid_at": time.Now(),
		}).Error
}
