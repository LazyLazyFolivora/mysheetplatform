package pkg

import "gorm.io/gorm"

func Paginate[T any](db *gorm.DB, page, pageSize int, order string) ([]T, int64, error) {
	var items []T
	var total int64
	query := db.Model(new(T))
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order(order).Find(&items).Error
	return items, total, err
}
