package oprtype

import (
	"paccount/pkg/account"
)

type OprTypeRepository interface {
	CalcLimit(oprID uint64, value float64,
		acc account.Entity) (account.Entity, float64)
	ValidateID(oprID uint64) bool
}
