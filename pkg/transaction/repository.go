package transaction

import (
	"paccount/pkg/account"
)

type TxRepository interface {
	Create(tx Entity) (Entity, error)
	Update(tx Entity) (Entity, error)
	GetAccount(id uint64) (account.Entity, error)
	UpdateAccount(acc account.Entity) error
	AccountHasLimit(acc account.Entity, input Entity) bool
	GetAllTxsWithCredit(accID uint64) ([]Entity, error)
	GetAllTxsWithoutBalance(accID uint64) ([]Entity, error)
	ProcessTx(input Entity) (Entity, error)
	GetAllTxs(accID uint64) ([]Entity, error)
	UpdateBalance(input Entity) (Entity, error)
}
