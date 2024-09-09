package internal

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseReturnEntity struct {
	ID                   uint           `gorm:"primaryKey;autoIncrement"`
	Code                 string         `gorm:"unique;not null"`
	InventoryCode        *string        `gorm:"column:inventory_code"` // Nullable, hence pointer
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

func (PurchaseReturnEntity) TableName() string {
	return "purchase_returns"
}

type IPurchaseReturnRepository interface {
	FindPurchaseReturnWithOutBankTransaction() (*[]PurchaseReturnEntity, error)
}

type purchaseReturnRepository struct {
	db *gorm.DB
}

func NewPurchaseReturnRepository(db *gorm.DB) IPurchaseReturnRepository {
	return &purchaseReturnRepository{
		db,
	}
}

// FindAll implements IMemberRepository.
func (p *purchaseReturnRepository) FindPurchaseReturnWithOutBankTransaction() (*[]PurchaseReturnEntity, error) {
	purchaseReturns := new([]PurchaseReturnEntity)
	err := p.db.
		Table("purchase_returns pr").
		Select("pr.*").
		Joins("LEFT JOIN bank_transactions bt ON bt.purchase_return_id = pr.id").
		Where("bt.id IS NULL").
		Where("pr.is_fictive = ?", false).
		Find(&purchaseReturns).Error

	if err != nil {
		return nil, err
	}
	return purchaseReturns, nil
}
