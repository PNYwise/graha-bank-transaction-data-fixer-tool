package internal

import (
	"gorm.io/gorm"
)


type SaleReturnEntity struct {
	ID                  uint          `gorm:"primaryKey;autoIncrement"`
	Code                string        `gorm:"unique;not null"`
	Note                string        `gorm:"type:text;not null;default:'-'"`
	Date                string     `gorm:"type:date;not null"`
	PpnInValue          int64        `gorm:"column:ppn_in_value;default:0"`
	PpnInPercent        int64        `gorm:"column:ppn_in_percent;default:0"`
	Total               int64         `gorm:"not null"`
	TotalIncludingPpn   int64         `gorm:"column:total_including_ppn;not null;default:0"`
	TotalNotIncludingPpn int64        `gorm:"column:total_not_including_ppn;not null;default:0"`
	IsLastMonthProcessed bool         `gorm:"column:is_last_month_processed;not null;default:false"`
}


func (SaleReturnEntity) TableName() string {
	return "sale_returns"
}

type ISaleReturnRepository interface {
	FindSaleReturnWithOutBankTransaction() (*[]SaleReturnEntity, error)
}

type saleReturnRepository struct {
	db *gorm.DB
}

func NewSaleReturnRepository(db *gorm.DB) ISaleReturnRepository {
	return &saleReturnRepository{
		db,
	}
}

func (s *saleReturnRepository) FindSaleReturnWithOutBankTransaction() (*[]SaleReturnEntity, error) {
	saleReturns := new([]SaleReturnEntity)
	err := s.db.
		Table("sale_reutrns sr").
		Select("sr.*").
		Joins("LEFT JOIN bank_transactions bt ON bt.sale_return_id = sr.id").
		Where("bt.id IS NULL").
		Find(&saleReturns).Error

	if err != nil {
		return nil, err
	}
	return saleReturns, nil
}