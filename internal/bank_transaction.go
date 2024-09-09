package internal

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BankTransactionEntity struct {
	ID                   uint                        `gorm:"primaryKey;autoIncrement"`
	BankNumber           string                      `gorm:"not null" json:"bank_number"`
	Date                 string                      `gorm:"not null"`
	Note                 string                      `gorm:"type:text;default:'-';not null"`
	Amount               float64                     `gorm:"not null"`
	IsMain               bool                        `gorm:"default:false;not null" json:"is_main"`
	CreatedAt            time.Time                   `gorm:"autoCreateTime"`
	UpdatedAt            time.Time                   `gorm:"autoUpdateTime"`
	BankTransactionItems []BankTransactionItemEntity `gorm:"foreignKey:BankTransactionID"`
	Purchase             *PurchaseEntity             `gorm:"foreignKey:PurchaseID"`
	PurchaseReturn       *PurchaseReturnEntity       `gorm:"foreignKey:PurchaseReturnID"`
	Sale                 *SaleEntity                 `gorm:"foreignKey:SaleID"`
	SaleReturn           *SaleReturnEntity           `gorm:"foreignKey:SaleReturnID"`
	PurchaseID           *uint
	PurchaseReturnID     *uint
	SaleID               *uint
	SaleReturnID         *uint
}

func (BankTransactionEntity) TableName() string {
	return "bank_transactions"
}

type IBankTransactionRepository interface {
	FindBankTransactionWithDateNotEqualPurchaseDate() (*[]BankTransactionEntity, error)
	FindBankTransactionWithDateNotEqualPurchaseReturnDate() (*[]BankTransactionEntity, error)
	FindBankTransactionWithDateNotEqualSaleDate() (*[]BankTransactionEntity, error)
	FindBankTransactionWithDateNotEqualSaleReturnDate() (*[]BankTransactionEntity, error)
	FindLastCode(date string, prefix string) (string, error)
	Update(bankTransaction *BankTransactionEntity)
	Create(bankTransaction *BankTransactionEntity)
}

type bankTransactionRepository struct {
	db *gorm.DB
}

func NewBankTransactionRepository(db *gorm.DB) IBankTransactionRepository {
	return &bankTransactionRepository{
		db,
	}
}

func (b *bankTransactionRepository) FindBankTransactionWithDateNotEqualPurchaseDate() (*[]BankTransactionEntity, error) {
	bankTransactions := new([]BankTransactionEntity)

	err := b.db.Table("bank_transactions").
		Joins("Purchase").
		Where(`"Purchase".date <> bank_transactions.date`).
		Order("bank_transactions.purchase_id ASC").
		Find(&bankTransactions).Error

	if err != nil {
		return nil, err
	}
	return bankTransactions, nil
}

func (b *bankTransactionRepository) FindBankTransactionWithDateNotEqualPurchaseReturnDate() (*[]BankTransactionEntity, error) {
	bankTransactions := new([]BankTransactionEntity)

	err := b.db.Table("bank_transactions").
		Joins("PurchaseReturn").
		Where(`"PurchaseReturn".date <> bank_transactions.date`).
		Order("bank_transactions.purchase_return_id ASC").
		Find(&bankTransactions).Error

	if err != nil {
		return nil, err
	}
	return bankTransactions, nil
}

func (b *bankTransactionRepository) FindBankTransactionWithDateNotEqualSaleDate() (*[]BankTransactionEntity, error) {
	bankTransactions := new([]BankTransactionEntity)

	err := b.db.Table("bank_transactions").
		Joins("Sale").
		Where(`"Sale".date <> bank_transactions.date`).
		Order("bank_transactions.sale_id ASC").
		Find(&bankTransactions).Error

	if err != nil {
		return nil, err
	}
	return bankTransactions, nil
}

func (b *bankTransactionRepository) FindBankTransactionWithDateNotEqualSaleReturnDate() (*[]BankTransactionEntity, error) {
	bankTransactions := new([]BankTransactionEntity)

	err := b.db.Table("bank_transactions").
		Joins("SaleReturn").
		Where(`"SaleReturn".date <> bank_transactions.date`).
		Order("bank_transactions.sale_return_id ASC").
		Find(&bankTransactions).Error

	if err != nil {
		return nil, err
	}
	return bankTransactions, nil
}

func (b *bankTransactionRepository) FindLastCode(date string, prefix string) (string, error) {

	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", fmt.Errorf("GetNewCode invalid date format: %v", err)
	}
	formattedDate := parsedDate.Format("20060102")

	bankTransaction := new(BankTransactionEntity)
	err = b.db.Table("bank_transactions").
		Where("bank_transactions.bank_number like ?", prefix+formattedDate+"%").
		Order("bank_transactions.bank_number DESC").
		Limit(1).
		Find(&bankTransaction).Error

	if err != nil {
		return "", err
	}
	return bankTransaction.BankNumber, nil
}

func (b *bankTransactionRepository) Update(bankTransaction *BankTransactionEntity) {
	b.db.Model(&BankTransactionEntity{ID: bankTransaction.ID}).Select("BankNumber", "Date").
		Updates(BankTransactionEntity{BankNumber: bankTransaction.BankNumber, Date: bankTransaction.Date})
}

func (b *bankTransactionRepository) Create(bankTransaction *BankTransactionEntity) {
	b.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Omit(clause.Associations).Create(bankTransaction).Error; err != nil {
			return err
		}

		var allBankTransactionItems []BankTransactionItemEntity
		for _, bankTransactionItem := range bankTransaction.BankTransactionItems {
			allBankTransactionItems = append(allBankTransactionItems, BankTransactionItemEntity{
				CardNumber:              bankTransactionItem.CardNumber,
				Note:                    bankTransactionItem.Note,
				SubAmount:               bankTransactionItem.SubAmount,
				BankTransactionID:       bankTransaction.ID,
				MasterBankTransactionID: bankTransactionItem.MasterBankTransactionID,
			})
		}

		if err := tx.Create(&allBankTransactionItems).Error; err != nil {
			return err
		}
		return nil
	})
}
