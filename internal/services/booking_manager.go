package services

import (
	"applicationDesignTest/internal/model"
	"context"
)

type OrderService interface {
	CreateOrders(ctx context.Context, orders []model.Order) ([]model.Order, error)
}

type RoomService interface {
	CheckRoomsAvailability(ctx context.Context, orders []model.Order) (bool, error)
	BookRooms(ctx context.Context, orders []model.Order) error
}

type BookingManager interface {
	HandleOrders(ctx context.Context, orders []model.Order) ([]model.Order, error)
}

type bookingManager struct {
	orderService OrderService
	roomService  RoomService
}

func NewBookingManager(orderService OrderService, roomService RoomService) BookingManager {
	return &bookingManager{orderService: orderService, roomService: roomService}
}

func (bm *bookingManager) HandleOrders(ctx context.Context, orders []model.Order) ([]model.Order, error) {
	ok, err := bm.roomService.CheckRoomsAvailability(ctx, orders)
	if err != nil {
		//log error
		return nil, err
	}
	if !ok {
		//log error
		return nil, model.ErrInsufficientQuota
	}

	err = bm.roomService.BookRooms(ctx, orders)
	if err != nil {
		//log error
		return nil, err
	}

	orders, err = bm.orderService.CreateOrders(ctx, orders)
	if err != nil {
		//log error
		return nil, err
	}

	//log result
	return orders, nil
}
