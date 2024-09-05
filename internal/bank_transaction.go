package internal

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BankTransactionEntity struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	BankNumber string         `gorm:"not null" json:"bank_number"`
	Date       string         `gorm:"not null"`
	Note       string         `gorm:"type:text;default:'-';not null"`
	Amount     float64        `gorm:"not null"`
	IsMain     bool           `gorm:"default:false;not null" json:"is_main"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	PurchaseID uint           `gorm:"not null"`
	Purchase   PurchaseEntity `gorm:"foreignKey:PurchaseID"`
}

func (BankTransactionEntity) TableName() string {
	return "bank_transactions"
}

type IBankTransactionRepository interface {
	FindBankTransactionWithDateNotEqualPurchaseDate() (*[]BankTransactionEntity, error)
	FindBankTransactionWithDateNotEqualPurchaseReturnDate() (*[]BankTransactionEntity, error)
	FindLastCode(date string, prefix string) (string, error)
	Update(bankTransaction *BankTransactionEntity)
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
		Where("bank_transactions.id = 381").
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

func (b *bankTransactionRepository) FindLastCode(date string, prefix string) (string, error) {

	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", fmt.Errorf("GetNewCode invalid date format: %v", err)
	}
	formattedDate := parsedDate.Format("20060102")

	bankTransaction := new(BankTransactionEntity)
	err = b.db.Table("bank_transactions").
		Where("bank_transactions.bank_number like ?", prefix+formattedDate+"%").
		Order("bank_transactions.date DESC").
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
