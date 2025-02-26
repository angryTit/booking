package services

import (
	"applicationDesignTest/internal/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrders(ctx context.Context, orders []model.Order) ([]model.Order, error) {
	args := m.Called(ctx, orders)
	return args.Get(0).([]model.Order), args.Error(1)
}

type MockRoomService struct {
	mock.Mock
}

func (m *MockRoomService) CheckRoomsAvailability(ctx context.Context, orders []model.Order) (bool, error) {
	args := m.Called(ctx, orders)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoomService) BookRooms(ctx context.Context, orders []model.Order) error {
	args := m.Called(ctx, orders)
	return args.Error(0)
}

func TestBookingManager_HandleOrders(t *testing.T) {
	mockDate := func() time.Time {
		return time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	}

	t.Run("successful booking", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		mockRoomService := new(MockRoomService)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   mockDate(),
				CheckOut:  mockDate().Add(24 * time.Hour),
				Rooms:     2,
			},
		}

		expectedOrders := []model.Order{
			{
				OrderId:   "order1",
				UserEmail: "user1@example.com",
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   mockDate(),
				CheckOut:  mockDate().Add(24 * time.Hour),
				Rooms:     2,
			},
		}

		mockRoomService.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(true, nil)
		mockRoomService.On("BookRooms", mock.Anything, testOrders).Return(nil)
		mockOrderService.On("CreateOrders", mock.Anything, testOrders).Return(expectedOrders, nil)

		bookingManager := NewBookingManager(mockOrderService, mockRoomService)

		result, err := bookingManager.HandleOrders(context.Background(), testOrders)

		assert.NoError(t, err, "expected no error")
		assert.Equal(t, expectedOrders, result, "expected orders to be created")

		mockRoomService.AssertExpectations(t)
		mockOrderService.AssertExpectations(t)
	})

	t.Run("rooms not available", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		mockRoomService := new(MockRoomService)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   mockDate(),
				CheckOut:  mockDate().Add(24 * time.Hour),
				Rooms:     10,
			},
		}

		mockRoomService.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(false, nil)

		bookingManager := NewBookingManager(mockOrderService, mockRoomService)

		result, err := bookingManager.HandleOrders(context.Background(), testOrders)

		assert.Equal(t, model.ErrInsufficientQuota, err, "expected ErrInsufficientQuota error")
		assert.Nil(t, result, "expected no orders to be returned")

		mockRoomService.AssertExpectations(t)
	})

	t.Run("error checking availability", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		mockRoomService := new(MockRoomService)

		testOrders := []model.Order{
			{
				HotelId:   "nonexistent_hotel",
				RoomClass: "standard",
				CheckIn:   mockDate(),
				CheckOut:  mockDate().Add(24 * time.Hour),
				Rooms:     2,
			},
		}

		mockRoomService.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(false, model.ErrNotFound)

		bookingManager := NewBookingManager(mockOrderService, mockRoomService)

		result, err := bookingManager.HandleOrders(context.Background(), testOrders)

		assert.Equal(t, model.ErrNotFound, err, "expected ErrNotFound error")
		assert.Nil(t, result, "expected no orders to be returned")

		mockRoomService.AssertExpectations(t)
	})

	t.Run("error booking rooms", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		mockRoomService := new(MockRoomService)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   mockDate(),
				CheckOut:  mockDate().Add(24 * time.Hour),
				Rooms:     2,
			},
		}

		bookingError := errors.New("booking error")

		mockRoomService.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(true, nil)
		mockRoomService.On("BookRooms", mock.Anything, testOrders).Return(bookingError)

		bookingManager := NewBookingManager(mockOrderService, mockRoomService)

		result, err := bookingManager.HandleOrders(context.Background(), testOrders)

		assert.Equal(t, bookingError, err, "expected booking error")
		assert.Nil(t, result, "expected no orders to be returned")

		mockRoomService.AssertExpectations(t)
	})

	t.Run("error creating orders", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		mockRoomService := new(MockRoomService)

		testOrders := []model.Order{
			{
				HotelId:   "hotel1",
				RoomClass: "standard",
				CheckIn:   mockDate(),
				CheckOut:  mockDate().Add(24 * time.Hour),
				Rooms:     2,
			},
		}

		createError := errors.New("create order error")

		mockRoomService.On("CheckRoomsAvailability", mock.Anything, testOrders).Return(true, nil)
		mockRoomService.On("BookRooms", mock.Anything, testOrders).Return(nil)
		mockOrderService.On("CreateOrders", mock.Anything, testOrders).Return([]model.Order{}, createError)

		bookingManager := NewBookingManager(mockOrderService, mockRoomService)

		result, err := bookingManager.HandleOrders(context.Background(), testOrders)

		assert.Equal(t, createError, err, "expected create order error")
		assert.Nil(t, result, "expected no orders to be returned")

		mockRoomService.AssertExpectations(t)
		mockOrderService.AssertExpectations(t)
	})
}
