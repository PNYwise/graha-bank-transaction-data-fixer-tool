package internal

import (
	"gorm.io/gorm"
)

type SaleEntity struct {
	ID                   uint   `gorm:"primaryKey;autoIncrement"`
	SaleCode             string `gorm:"column:sale_code;unique"`
	Date                 string `gorm:"type:date"`
	Total                int64  `gorm:"not null;type:bigint"`
	TotalIncludingPpn    int64  `gorm:"column:total_including_ppn;not null;default:0;type:bigint"`
	TotalNotIncludingPpn int64  `gorm:"column:total_not_including_ppn;not null;default:0;type:bigint"`
	DeliveryFee          int64  `gorm:"column:delivery_fee;type:bigint;default:0"`
	DownPayment          int64  `gorm:"column:down_payment;type:bigint;default:0"`
	PpnInValue           int64  `gorm:"column:ppn_in_value;default:0"`
	PpnInPercent         int64  `gorm:"column:ppn_in_percent;default:0"`
	IsSaleOrder          bool   `gorm:"column:is_sale_order;not null;default:false"`
}

func (SaleEntity) TableName() string {
	return "sales"
}

type ISaleRepository interface {
	FindSaleWithOutBankTransaction() (*[]SaleEntity, error)
}

type saleRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) ISaleRepository {
	return &saleRepository{
		db,
	}
}

func (s *saleRepository) FindSaleWithOutBankTransaction() (*[]SaleEntity, error) {
	sales := new([]SaleEntity)
	err := s.db.
		Table("sales s").
		Select("s.*").
		Joins("LEFT JOIN bank_transactions bt ON bt.sale_id = s.id").
		Where("bt.id IS NULL").
		Where("s.is_processed = ?", true).
		Find(&sales).Error

	if err != nil {
		return nil, err
	}
	return sales, nil
}
