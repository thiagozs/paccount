package models

// OprType model structure
type OprType struct {
	ID          uint64 `json:"id" gorm:"primary_key"`
	Description string `json:"description" gorm:"type:varchar(200);"`
	CreatedAt   int32  `json:"created_at" gorm:"type:int;not null;"`
	UpdatedAt   int32  `json:"updated_at" gorm:"type:int;not null;"`
}

// TableName convention gorm ocr
func (o OprType) TableName() string {
	return "operation_type"
}
