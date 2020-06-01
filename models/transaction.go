package models

// Transaction model structure
type Transaction struct {
	ID          uint64  `json:"id" gorm:"primary_key"`
	AccountID   uint64  `json:"account_id" gorm:"type:int;not null;"`
	OperationID uint64  `json:"operation_id" gorm:"type:int;not null;"`
	Amount      float64 `json:"amount" gorm:"type:real;not null;"`
	CreatedAt   int32   `json:"created_at" gorm:"type:int;not null;"`
}

// TableName convention gorm ocr
func (o Transaction) TableName() string {
	return "transaction"
}
