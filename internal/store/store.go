// Package store provides thread safe in-memory key-value store
package store

import (
	"fmt"
	"sync"
)

type Store struct {
	mutex sync.RWMutex
	store map[string]any
}

var GlobalStore = &Store{
	store: make(map[string]any),
}

func (s *Store) Get(key string) (bool, any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.store) < 1 {
		return false, nil
	}

	value, exists := s.store[key]
	if !exists {
		return false, nil
	}

	return true, value
}

func (s *Store) Set(key string, value any) (bool, string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.store[key]
	if exists {
		return false, "key already exists in store"
	}

	fmt.Println("Successful")

	s.store[key] = value
	return true, ""
}
