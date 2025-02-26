package services

import (
	"applicationDesignTest/internal/model"
	"context"
)

type OrderRepository interface {
	AddOrder(ctx context.Context, orders []model.Order) error
}

type orderService struct {
	orderRepo OrderRepository
}

func NewOrderService(orderRepo OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

func (s *orderService) CreateOrders(ctx context.Context, orders []model.Order) ([]model.Order, error) {
	err := s.orderRepo.AddOrder(ctx, orders)
	if err != nil {
		//mb wrap error, log error
		return nil, err
	}
	//log result
	return orders, nil
}
