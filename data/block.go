package data

type BlockRepository interface {
	GetCurrentBlock() (int64, error)
	UpdateCurrentBlock(block int64) error
}

func (s *InMemoryStorage) GetCurrentBlock() (int64, error) {
	return s.currentBlock, nil
}

func (s *InMemoryStorage) UpdateCurrentBlock(block int64) error {
	s.currentBlock = block

	return nil
}
