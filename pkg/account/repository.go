package account

type AccRepository interface {
	Create(acc Entity) (*Entity, error)
	Find(uit uint64) (Entity, error)
	Update(acc Entity) error
}
