package model

import "time"

type SheetOrder struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	OrderNo      string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_no"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	SheetMusicID uint      `gorm:"index;not null" json:"sheet_music_id"`
	Amount       float64   `json:"amount"`
	Status       string    `gorm:"type:varchar(20);default:pending" json:"status"`
	AlipayTradeNo string   `gorm:"type:varchar(64)" json:"alipay_trade_no"`
	PaidAt       *time.Time `json:"paid_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SheetOrder) TableName() string {
	return "sheet_order"
}
