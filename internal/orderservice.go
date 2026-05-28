package main

import (
	"errors"
)

// Order представляет собой структуру заказа.
type Order struct {
	ID    int
	Items []string
}

// OrderService предоставляет методы для работы с заказами.
type OrderService struct {
	orders map[int]Order // Хранение заказов по ID
	nextID int           // Следующий доступный ID для нового заказа
}

// NewOrderService создает новый экземпляр OrderService.
func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[int]Order),
		nextID: 1,
	}
}

// CreateOrder создает новый заказ и возвращает его.
func (s *OrderService) CreateOrder(items []string) Order {
	order := Order{
		ID:    s.nextID,
		Items: items,
	}
	s.orders[s.nextID] = order
	s.nextID++
	return order
}

// GetOrder возвращает заказ по ID.
func (s *OrderService) GetOrder(id int) (Order, error) {
	order, exists := s.orders[id]
	if !exists {
		return Order{}, errors.New("order not found")
	}
	return order, nil
}

// DeleteOrder удаляет заказ по ID.
func (s *OrderService) DeleteOrder(id int) error {
	if _, exists := s.orders[id]; !exists {
		return errors.New("order not found")
	}
	delete(s.orders, id)
	return nil
}
