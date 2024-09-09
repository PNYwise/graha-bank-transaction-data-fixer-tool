package internal

import "gorm.io/gorm"

type MasterBankTransactionEntity struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Code string `gorm:"unique;not null"`
	Name string `gorm:"not null"`
}

func (MasterBankTransactionEntity) TableName() string {
	return "master_bank_transactions"
}

type IMasterBankTransactionRepository interface {
	Find(codes []string) ([]MasterBankTransactionEntity, error)
}

type masterBankTransactionRepository struct {
	db *gorm.DB
}

func NewMasterBankTransactionRepository(db *gorm.DB) IMasterBankTransactionRepository {
	return &masterBankTransactionRepository{
		db,
	}
}

// Find implements IMasterBankTransactionRepository.
func (m *masterBankTransactionRepository) Find(codes []string) ([]MasterBankTransactionEntity, error) {
	var masterBankTransactions  []MasterBankTransactionEntity

	if err := m.db.Where("code IN ?", codes).Find(&masterBankTransactions).Error; err != nil {
		return nil, err
	}
	return masterBankTransactions, nil
}
