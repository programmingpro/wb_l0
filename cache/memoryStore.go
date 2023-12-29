package cache

import (
	"awesomeProject/models"
	"sync"
)

type InMemoryStore struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func NewInMemoryStore() InMemoryStore {
	return InMemoryStore{
		orders: make(map[string]models.Order),
	}
}

func (s *InMemoryStore) AddOrder(order models.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUID] = order
}

func (s *InMemoryStore) GetOrder(orderUID string) (models.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderUID]
	return order, ok
}

func (s *InMemoryStore) GetAllOrders() []models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]models.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}
