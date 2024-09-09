package internal

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseEntity struct {
	ID                   uint           `gorm:"primaryKey;autoIncrement"`
	Code                 string         `gorm:"unique;not null"`
	Note                 string         `gorm:"type:text;default:'-'"`
	Date                 string         `gorm:"not null"`
	PPNInPercent         float64        `gorm:"not null" json:"ppn_in_percent"`
	PPNInValue           float64        `gorm:"not null" json:"ppn_in_value"`
	Total                int64          `gorm:"not null"`
	TotalIncludingPPN    int64          `gorm:"default:0;not null" json:"total_including_ppn"`
	TotalNotIncludingPPN int64          `gorm:"default:0;not null" json:"total_not_including_ppn"`
	IsFictive            bool           `gorm:"default:false;not null" json:"is_fictive"`
	CreatedAt            time.Time      `gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

func (PurchaseEntity) TableName() string {
	return "purchases"
}

type IPurchaseRepository interface {
	FindPurchaseWithOutBankTransaction() (*[]PurchaseEntity, error)
}

type purchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) IPurchaseRepository {
	return &purchaseRepository{
		db,
	}
}

func (p *purchaseRepository) FindPurchaseWithOutBankTransaction() (*[]PurchaseEntity, error) {
	purchases := new([]PurchaseEntity)
	err := p.db.
		Table("purchases p").
		Select("p.*").
		Joins("LEFT JOIN bank_transactions bt ON bt.purchase_id = p.id").
		Where("bt.id IS NULL").
		Where("p.is_fictive = ?", false).
		Find(&purchases).Error

	if err != nil {
		return nil, err
	}
	return purchases, nil
}
