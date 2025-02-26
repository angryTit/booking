package repo

import (
	"applicationDesignTest/internal/model"
	"context"
	"sync"
)

type inMemoryOrder struct {
	ordersLock  sync.RWMutex
	usersOrders map[string][]model.Order
}

func NewInMemoryOrder() *inMemoryOrder {
	return &inMemoryOrder{
		usersOrders: make(map[string][]model.Order),
	}
}

func (r *inMemoryOrder) AddOrder(ctx context.Context, orders []model.Order) error {
	r.ordersLock.Lock()
	defer r.ordersLock.Unlock()
	for _, order := range orders {
		uOrders := r.usersOrders[order.UserEmail]
		uOrders = append(uOrders, order)
		r.usersOrders[order.UserEmail] = uOrders
	}
	return nil
}
