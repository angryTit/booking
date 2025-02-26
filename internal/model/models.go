package model

import "time"

type Order struct {
	OrderId   string
	HotelId   string
	RoomClass string
	CheckIn   time.Time
	CheckOut  time.Time
	Rooms     uint32
	UserEmail string
}
