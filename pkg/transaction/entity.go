package transaction

// Entity model structure
type Entity struct {
	ID          uint64  `json:"id" gorm:"primary_key"`
	AccountID   uint64  `json:"account_id" gorm:"type:int;not null;"`
	OperationID uint64  `json:"operation_id" gorm:"type:int;not null;"`
	Amount      float64 `json:"amount" gorm:"type:real;not null;"`
	Balance     float64 `json:"balance" gorm:"type:real;not null;"`
	CreatedAt   int32   `json:"created_at" gorm:"type:int;not null;"`
	UpdatedAt   int32   `json:"updated_at" gorm:"type:int;not null;"`
}

// TableName convention gorm ocr
func (o Entity) TableName() string {
	return "transaction"
}
