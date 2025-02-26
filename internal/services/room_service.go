package services

import (
	"applicationDesignTest/internal/model"
	"context"
)

type AvailabilityRepository interface {
	CheckRoomsAvailability(ctx context.Context, orders []model.Order) (bool, error)
	BookRooms(ctx context.Context, orders []model.Order) error
}

type roomService struct {
	availabilityRepo AvailabilityRepository
}

func NewRoomService(availabilityRepo AvailabilityRepository) RoomService {
	return &roomService{availabilityRepo: availabilityRepo}
}

func (s *roomService) CheckRoomsAvailability(ctx context.Context, orders []model.Order) (bool, error) {
	ok, err := s.availabilityRepo.CheckRoomsAvailability(ctx, orders)
	if err != nil {
		//mb wrap error, log error
		return false, err
	}
	// log result
	return ok, nil
}

func (s *roomService) BookRooms(ctx context.Context, orders []model.Order) error {
	err := s.availabilityRepo.BookRooms(ctx, orders)
	if err != nil {
		//mb wrap error, log error
		return err
	}
	//log result
	return nil
}
