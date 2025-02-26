package services

import (
	"applicationDesignTest/internal/model"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) AddOrder(ctx context.Context, orders []model.Order) error {
	args := m.Called(ctx, orders)
	return args.Error(0)
}

func TestOrderService_CreateOrders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockOrderRepository)

		testOrders := []model.Order{
			{
				UserEmail: "test@example.com",
				HotelId:   "hotel1",
				RoomClass: "standard",
				Rooms:     2,
			},
			{
				UserEmail: "test@example.com",
				HotelId:   "hotel2",
				RoomClass: "deluxe",
				Rooms:     1,
			},
		}

		mockRepo.On("AddOrder", mock.Anything, testOrders).Return(nil)

		service := NewOrderService(mockRepo)

		result, err := service.CreateOrders(context.Background(), testOrders)

		assert.NoError(t, err)
		assert.Equal(t, testOrders, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(MockOrderRepository)

		testOrders := []model.Order{
			{
				UserEmail: "test@example.com",
				HotelId:   "hotel1",
				RoomClass: "standard",
				Rooms:     2,
			},
		}

		expectedError := errors.New("database error")

		mockRepo.On("AddOrder", mock.Anything, testOrders).Return(expectedError)

		service := NewOrderService(mockRepo)

		result, err := service.CreateOrders(context.Background(), testOrders)

		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

}
