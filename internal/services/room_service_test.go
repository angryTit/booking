package services

import (
	"applicationDesignTest/internal/model"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAvailabilityRepository struct {
	mock.Mock
}

func (m *MockAvailabilityRepository) CheckRoomsAvailability(ctx context.Context, orders []model.Order) (bool, error) {
	args := m.Called(ctx, orders)
	return args.Bool(0), args.Error(1)
}

func (m *MockAvailabilityRepository) BookRooms(ctx context.Context, orders []model.Order) error {
	args := m.Called(ctx, orders)
	return args.Error(0)
}

func TestRoomService_CheckRoomsAvailability(t *testing.T) {
	t.Run("rooms available", func(t *testing.T) {
		mockRepo := new(MockAvailabilityRepository)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
				Rooms:     2,
			},
		}

		mockRepo.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(true, nil)

		service := NewRoomService(mockRepo)

		result, err := service.CheckRoomsAvailability(context.Background(), testOrders)

		assert.NoError(t, err, "expected no error")
		assert.True(t, result, "expected rooms to be available")

		mockRepo.AssertExpectations(t)
	})

	t.Run("rooms not available", func(t *testing.T) {
		mockRepo := new(MockAvailabilityRepository)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
				Rooms:     10,
			},
		}

		mockRepo.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(false, nil)

		service := NewRoomService(mockRepo)

		result, err := service.CheckRoomsAvailability(context.Background(), testOrders)

		assert.NoError(t, err, "expected no error")
		assert.False(t, result, "expected rooms to be not available")

		mockRepo.AssertExpectations(t)
	})

	t.Run("error checking availability", func(t *testing.T) {
		mockRepo := new(MockAvailabilityRepository)

		testOrders := []model.Order{
			{
				HotelId:   "nonexistent_hotel",
				RoomClass: "standard",
				CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
				Rooms:     2,
			},
		}

		mockRepo.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(false, model.ErrNotFound)

		service := NewRoomService(mockRepo)

		result, err := service.CheckRoomsAvailability(context.Background(), testOrders)

		assert.Equal(t, model.ErrNotFound, err, "expected ErrNotFound error")
		assert.False(t, result, "expected false when error occurs")

		mockRepo.AssertExpectations(t)
	})
}

func TestRoomService_BookRooms(t *testing.T) {
	t.Run("successful booking", func(t *testing.T) {
		mockRepo := new(MockAvailabilityRepository)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
				Rooms:     2,
			},
		}

		mockRepo.On("BookRooms", mock.Anything, testOrders).Return(nil)

		service := NewRoomService(mockRepo)

		err := service.BookRooms(context.Background(), testOrders)

		assert.NoError(t, err, "expected successful booking without errors")

		mockRepo.AssertExpectations(t)
	})

	t.Run("booking error", func(t *testing.T) {
		mockRepo := new(MockAvailabilityRepository)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
				Rooms:     10,
			},
		}

		expectedError := model.ErrInsufficientQuota

		mockRepo.On("BookRooms", mock.Anything, testOrders).Return(expectedError)

		service := NewRoomService(mockRepo)

		err := service.BookRooms(context.Background(), testOrders)

		assert.Equal(t, expectedError, err, "expected ErrInsufficientQuota error")

		mockRepo.AssertExpectations(t)
	})
}
