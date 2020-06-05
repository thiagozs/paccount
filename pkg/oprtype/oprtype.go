package oprtype

import (
	"paccount/database"
	"paccount/pkg/account"
	"time"
)

type repo struct {
	db database.IGormRepo
}

func New(db database.IGormRepo) OprTypeRepository {
	return repo{db}
}

func (o repo) CalcLimit(oprID uint64, value float64,
	acc account.Entity) (account.Entity, float64) {

	switch oprID {
	case COMPRA_AVISTA, COMPRA_PARCELADO, SAQUE:
		// -
		if value > 0 {
			value = value * (-1)
		}
	case PAGAMENTO:
		// +
		if value < 0 {
			value = value * (-1)
		}
	}

	acc.Limit += value

	timeStamp := time.Now()
	acc.UpdatedAt = int32(timeStamp.Unix())

	return acc, value
}

func (o repo) ValidateID(oprID uint64) bool {
	switch oprID {
	case COMPRA_AVISTA, COMPRA_PARCELADO, SAQUE, PAGAMENTO:
	default:
		return false
	}
	return true
}
