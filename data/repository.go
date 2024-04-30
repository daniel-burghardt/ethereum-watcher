package data

type Repository interface {
	SubscriberRepository
	TransactionRepository
	BlockRepository
}

type InMemoryStorage struct {
	transactions []Transaction
	subscribers  []string
	currentBlock int
}

// Make sure InMemoryStorage implements the Repository interface
var _ Repository = (*InMemoryStorage)(nil)

func NewInMemoryStorage() Repository {
	return &InMemoryStorage{
		transactions: []Transaction{},
		subscribers:  []string{},
		currentBlock: 0,
	}
}
