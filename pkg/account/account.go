package account

import (
	"paccount/database"
	"time"
)

type repo struct {
	DB database.IGormRepo
}

func New(db database.IGormRepo) AccRepository {
	return repo{db}
}

func (a repo) Create(acc Entity) (*Entity, error) {

	timeStamp := time.Now()
	model := &Entity{
		DocNumber: acc.DocNumber,
		Limit:     acc.Limit,
		CreatedAt: int32(timeStamp.Unix()),
		UpdatedAt: int32(timeStamp.Unix()),
	}
	if err := a.DB.Create(model); err != nil {
		return &Entity{}, err
	}
	return model, nil
}

func (a repo) Find(uit uint64) (Entity, error) {
	var acc Entity
	if err := a.DB.FindOne(Entity{ID: uit}, &acc); err != nil {
		return acc, err
	}
	return acc, nil
}

func (a repo) Update(acc Entity) error {
	if err := a.DB.Update(acc); err != nil {
		return err
	}
	return nil
}
