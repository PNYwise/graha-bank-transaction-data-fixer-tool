package internal

type BankTransactionItemEntity struct {
	ID                      uint    `gorm:"primaryKey;autoIncrement"`
	CardNumber              *string `gorm:"column:card_number"`
	Note                    string  `gorm:"type:text;not null;default:'-'"`
	SubAmount               float64 `gorm:"column:sub_amount;not null"`
	BankTransactionID       uint    `gorm:"not null"`
	MasterBankTransactionID uint    `gorm:"not null"`
}

func (BankTransactionItemEntity) TableName() string {
	return "bank_transaction_items"
}
