package items

import (
	"fmt"
	"sync"
)

type StrictItem struct {
	mu      sync.RWMutex
	data    map[string]interface{}
	allowed map[string]struct{}
}

func NewStrictItem(allowedFields []string) *StrictItem {
	a := make(map[string]struct{})
	for _, field := range allowedFields {
		a[field] = struct{}{}
	}
	return &StrictItem{
		data:    make(map[string]interface{}),
		allowed: a,
	}
}

func (s *StrictItem) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.allowed[key]; !ok {
		panic(fmt.Sprintf("字段 '%s' 不在预定义字段中", key))
	}
	s.data[key] = value
}

func (s *StrictItem) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *StrictItem) All() map[string]interface{} {
	return s.data
}
