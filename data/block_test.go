package data

import "testing"

func TestBlockRepository(t *testing.T) {
	repo := NewInMemoryStorage()

	// Assert CurrentBlock initial state
	block, err := repo.GetCurrentBlock()
	if err != nil {
		t.Errorf("Error = %v; want nil", err)
	}
	if block != 0 {
		t.Errorf("Block = %d; want 0", block)
	}

	// Assert UpdateCurrentBlock no error
	err = repo.UpdateCurrentBlock(5000000)
	if err != nil {
		t.Errorf("Error = %v; want nil", err)
	}

	// Assert CurrentBlock updated
	block, err = repo.GetCurrentBlock()
	if err != nil {
		t.Errorf("Error = %v; want nil", err)
	}
	if block != 5000000 {
		t.Errorf("Block = %d; want 5000000", block)
	}
}
