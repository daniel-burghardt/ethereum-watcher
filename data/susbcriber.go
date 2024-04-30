package data

import (
	"slices"
	"strings"
)

type SubscriberRepository interface {
	AddSubscriber(address string) error
	GetSubscribers() ([]string, error)
	SubscriberExists(address string) (bool, error)
}

func (s *InMemoryStorage) AddSubscriber(address string) error {
	address = strings.ToLower(address)

	if !slices.Contains(s.subscribers, address) {
		s.subscribers = append(s.subscribers, address)
	}

	return nil
}

func (s *InMemoryStorage) GetSubscribers() ([]string, error) {
	return s.subscribers, nil
}

func (s *InMemoryStorage) SubscriberExists(address string) (bool, error) {
	address = strings.ToLower(address)
	return slices.Contains(s.subscribers, address), nil
}
