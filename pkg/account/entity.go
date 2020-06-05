package account

// Account model structure
type Entity struct {
	ID        uint64  `json:"id" gorm:"primary_key"`
	Limit     float64 `json:"limit" gorm:"type:real"`
	DocNumber int32   `json:"document_number" gorm:"type:int;unique;not null;"`
	CreatedAt int32   `json:"created_at" gorm:"type:int;not null;"`
	UpdatedAt int32   `json:"updated_at" gorm:"type:int;not null;"`
}

// TableName convention gorm ocr
func (l Entity) TableName() string {
	return "account"
}
