package data

type BlockRepository interface {
	GetCurrentBlock() (int, error)
}

func (s *InMemoryStorage) GetCurrentBlock() (int, error) {
	return s.currentBlock, nil
}
