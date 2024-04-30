package data

import (
	"testing"
)

func TestSubscriberRepository(t *testing.T) {
	repo := NewInMemoryStorage()

	// Assert AddSubscriber no error
	err := repo.AddSubscriber("0xfe9e3cdD06Ab5ebb65b408E8503dD67CDe74478E")
	if err != nil {
		t.Errorf("Error = %v; want nil", err)
	}

	t.Run("subscriber_exists_case_insensitve", func(t *testing.T) {
		exists, err := repo.SubscriberExists("0xfe9e3cdd06ab5ebb65b408e8503dd67CDE74478E")
		if err != nil {
			t.Errorf("Error = %v; want nil", err)
		}
		if exists != true {
			t.Errorf("Exists = %t; want true", exists)
		}
	})

	t.Run("subscriber_not_exists", func(t *testing.T) {
		exists, err := repo.SubscriberExists("0x000000000000")
		if err != nil {
			t.Errorf("Error = %v; want nil", err)
		}
		if exists != false {
			t.Errorf("Exists = %t; want false", exists)
		}
	})

	t.Run("get_subscribers", func(t *testing.T) {
		subs, err := repo.GetSubscribers()
		if err != nil {
			t.Errorf("Error = %v; want nil", err)
		}
		if len(subs) != 1 {
			t.Errorf("Len = %d; want 1", len(subs))
		}
		if subs[0] != "0xfe9e3cdd06ab5ebb65b408e8503dd67cde74478e" {
			t.Errorf("Address = %s; want 0xfe9e3cdd06ab5ebb65b408e8503dd67cde74478e", subs[0])
		}
	})
}
