package transaction

import (
	"errors"
	"paccount/database"
	"paccount/pkg/account"
	"paccount/pkg/oprtype"
	"time"
)

type repo struct {
	DB      database.IGormRepo
	Account account.AccRepository
	Oprtype oprtype.OprTypeRepository
}

func New(db database.IGormRepo) TxRepository {
	ac := account.New(db)
	op := oprtype.New(db)
	return repo{DB: db, Account: ac, Oprtype: op}
}

func (o repo) Create(tx Entity) (Entity, error) {
	timeStamp := time.Now().Unix()
	model := Entity{
		Amount:      tx.Amount,
		OperationID: tx.OperationID,
		AccountID:   tx.AccountID,
		Balance:     tx.Balance,
		CreatedAt:   int32(timeStamp),
		UpdatedAt:   int32(timeStamp),
	}

	if err := o.DB.Create(&model); err != nil {
		return Entity{}, err
	}

	return model, nil
}

func (o repo) Update(tx Entity) (Entity, error) {
	timeStamp := time.Now()
	tx.UpdatedAt = int32(timeStamp.Unix())

	if err := o.DB.Update(tx); err != nil {
		return Entity{}, err
	}

	return tx, nil
}

func (o repo) GetAccount(id uint64) (account.Entity, error) {
	return o.Account.Find(id)
}

func (o repo) UpdateAccount(acc account.Entity) error {
	return o.Account.Update(acc)
}

func (o repo) AccountHasLimit(acc account.Entity, input Entity) bool {
	if input.OperationID < oprtype.PAGAMENTO && input.Amount > acc.Limit {
		return false
	}
	return true
}

func (o repo) GetAllTxsWithoutBalance(accID uint64) ([]Entity, error) {
	var txs []Entity
	if err := o.DB.GetDB().Table(Entity{}.TableName()).
		Where("account_id = ? AND amount < ? AND balance < ? AND operation_id in (?)",
			accID, 0, 0, []uint64{1, 2, 3}).
		Order("created_at asc").
		Find(&txs).Error; err != nil {
		return txs, err
	}
	return txs, nil
}

func (o repo) GetAllTxsWithCredit(accID uint64) ([]Entity, error) {
	var txs []Entity
	if err := o.DB.GetDB().Table(Entity{}.TableName()).
		Where("account_id = ? AND balance > ? AND operation_id in (?)",
			accID, 0, []uint64{4}).
		Order("created_at asc").
		Find(&txs).Error; err != nil {
		return txs, err
	}
	return txs, nil
}

func (o repo) GetAllTxs(accID uint64) ([]Entity, error) {
	var txs []Entity
	if err := o.DB.GetDB().Table(Entity{}.TableName()).
		Where("account_id = ?", accID).
		Order("created_at asc").
		Find(&txs).Error; err != nil {
		return txs, err
	}
	return txs, nil
}

func (o repo) ProcessTx(input Entity) (Entity, error) {

	if !o.Oprtype.ValidateID(input.OperationID) {
		return Entity{}, errors.New("invalid operation id")
	}

	acc, err := o.GetAccount(input.AccountID)
	if err != nil {
		return Entity{}, err
	}

	if !o.AccountHasLimit(acc, input) {
		return Entity{}, errors.New("no limit avaliable")
	}

	acc, amount := o.Oprtype.CalcLimit(input.OperationID, input.Amount, acc)
	if err := o.UpdateAccount(acc); err != nil {
		return Entity{}, err
	}
	input.Balance = amount
	input.Amount = amount

	if input.OperationID == oprtype.PAGAMENTO {
		accu, err := o.UpdateBalance(input)
		if err != nil {
			return Entity{}, err
		}
		input = accu
	}

	txr, err := o.Create(input)
	if err != nil {
		return Entity{}, errors.New("can't create a transaction")
	}

	return txr, nil
}

func (o repo) UpdateBalance(input Entity) (Entity, error) {
	//TODO : improve this method
	type ResultError struct {
		balance float64
		err     error
	}

	updateCredit := func(balance float64) ResultError {

		result := make(chan ResultError)
		go func() {
			defer close(result)

			txs, err := o.GetAllTxsWithCredit(input.AccountID)
			if err != nil {
				result <- ResultError{0.0, err}
				return
			}

			lastBalance := 0.0
			for _, t := range txs {
				lastBalance += t.Balance
				t.Balance = 0
				if _, err := o.Update(t); err != nil {
					result <- ResultError{0.0, err}
					break
				}
			}
			result <- ResultError{lastBalance, nil}

		}()
		return <-result
	}

	updateTxs := func(balance float64) ResultError {

		result := make(chan ResultError)
		go func() {
			defer close(result)

			txs, err := o.GetAllTxsWithoutBalance(input.AccountID)
			if err != nil {
				result <- ResultError{0.0, err}
				return
			}

			lastBalance := balance
			for _, t := range txs {
				if t.Amount < 0 && t.Amount < lastBalance {
					t.Balance = lastBalance
					t.UpdatedAt = int32(time.Now().Unix())
					if _, err := o.Update(t); err != nil {
						result <- ResultError{0.0, err}
						return
					}
					lastBalance += t.Amount
				}
			}
			result <- ResultError{lastBalance, nil}
		}()

		return <-result
	}

	rcredit := updateCredit(input.Balance)
	if rcredit.err != nil {
		return Entity{}, rcredit.err
	}
	input.Balance += rcredit.balance

	rupdate := updateTxs(input.Balance)
	if rupdate.err != nil {
		return Entity{}, rupdate.err
	}
	input.Balance = rupdate.balance

	return input, nil

}
