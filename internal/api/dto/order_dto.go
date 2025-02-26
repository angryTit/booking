package dto

import (
	"applicationDesignTest/internal/model"
	"time"
)

type OrderRequest struct {
	HotelID   string    `json:"hotel_id" binding:"required"`
	RoomClass string    `json:"room_class" binding:"required"`
	UserEmail string    `json:"email" binding:"required,email"`
	CheckIn   time.Time `json:"check_in" binding:"required"`
	CheckOut  time.Time `json:"check_out" binding:"required"`
	Rooms     uint32    `json:"rooms" binding:"required,min=1"`
}

type OrderResponse struct {
	OrderID   string    `json:"order_id"`
	HotelID   string    `json:"hotel_id"`
	RoomClass string    `json:"room_class"`
	UserEmail string    `json:"email"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
	Rooms     uint32    `json:"rooms"`
}

func (r *OrderRequest) ToModel() model.Order {
	return model.Order{
		HotelId:   r.HotelID,
		RoomClass: r.RoomClass,
		UserEmail: r.UserEmail,
		CheckIn:   r.CheckIn,
		CheckOut:  r.CheckOut,
		Rooms:     r.Rooms,
	}
}

func FromModel(order model.Order) OrderResponse {
	return OrderResponse{
		OrderID:   order.OrderId,
		HotelID:   order.HotelId,
		RoomClass: order.RoomClass,
		UserEmail: order.UserEmail,
		CheckIn:   order.CheckIn,
		CheckOut:  order.CheckOut,
		Rooms:     order.Rooms,
	}
}

func FromModelList(orders []model.Order) []OrderResponse {
	result := make([]OrderResponse, len(orders))
	for i, order := range orders {
		result[i] = FromModel(order)
	}
	return result
}
