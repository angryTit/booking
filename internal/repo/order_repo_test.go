package repo

import (
	"applicationDesignTest/internal/model"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryOrder_AddOrder(t *testing.T) {
	t.Run("successful adding one order", func(t *testing.T) {
		repo := NewInMemoryOrder()
		ctx := context.Background()

		userEmail := "test@example.com"
		order := model.Order{
			UserEmail: userEmail,
			HotelId:   "hotel1",
			RoomClass: "standard",
			CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			Rooms:     2,
		}

		err := repo.AddOrder(ctx, []model.Order{order})
		assert.NoError(t, err, "error adding order")

		repo.ordersLock.RLock()
		orders, exists := repo.usersOrders[userEmail]
		repo.ordersLock.RUnlock()

		assert.True(t, exists, "order not added for user %s", userEmail)
		assert.Len(t, orders, 1, "expected 1 order, got: %d", len(orders))
		assert.Equal(t, order.HotelId, orders[0].HotelId, "incorrect hotel")
		assert.Equal(t, order.RoomClass, orders[0].RoomClass, "incorrect room class")
	})

	t.Run("adding multiple orders for one user", func(t *testing.T) {
		repo := NewInMemoryOrder()
		ctx := context.Background()

		userEmail := "test@example.com"
		order1 := model.Order{
			UserEmail: userEmail,
			HotelId:   "hotel1",
			RoomClass: "standard",
			CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			Rooms:     2,
		}

		order2 := model.Order{
			UserEmail: userEmail,
			HotelId:   "hotel2",
			RoomClass: "deluxe",
			CheckIn:   time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			CheckOut:  time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC),
			Rooms:     1,
		}

		err := repo.AddOrder(ctx, []model.Order{order1, order2})
		assert.NoError(t, err, "error adding orders")

		repo.ordersLock.RLock()
		orders, exists := repo.usersOrders[userEmail]
		repo.ordersLock.RUnlock()

		assert.True(t, exists, "orders not added for user %s", userEmail)
		assert.Len(t, orders, 2, "expected 2 orders, got: %d", len(orders))

		foundHotel1 := false
		foundHotel2 := false

		for _, o := range orders {
			if o.HotelId == "hotel1" && o.RoomClass == "standard" {
				foundHotel1 = true
			}
			if o.HotelId == "hotel2" && o.RoomClass == "deluxe" {
				foundHotel2 = true
			}
		}

		assert.True(t, foundHotel1, "order for hotel1 not found")
		assert.True(t, foundHotel2, "order for hotel2 not found")
	})

	t.Run("adding orders for different users", func(t *testing.T) {
		repo := NewInMemoryOrder()
		ctx := context.Background()

		userEmail1 := "user1@example.com"
		userEmail2 := "user2@example.com"

		order1 := model.Order{
			UserEmail: userEmail1,
			HotelId:   "hotel1",
			RoomClass: "standard",
			CheckIn:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			CheckOut:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			Rooms:     2,
		}

		order2 := model.Order{
			UserEmail: userEmail2,
			HotelId:   "hotel2",
			RoomClass: "deluxe",
			CheckIn:   time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			CheckOut:  time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC),
			Rooms:     1,
		}

		err := repo.AddOrder(ctx, []model.Order{order1, order2})
		assert.NoError(t, err, "error adding orders")

		repo.ordersLock.RLock()
		orders1, exists1 := repo.usersOrders[userEmail1]
		orders2, exists2 := repo.usersOrders[userEmail2]
		repo.ordersLock.RUnlock()

		assert.True(t, exists1, "order not added for user %s", userEmail1)
		assert.True(t, exists2, "order not added for user %s", userEmail2)

		assert.Len(t, orders1, 1, "expected 1 order for user %s", userEmail1)
		assert.Len(t, orders2, 1, "expected 1 order for user %s", userEmail2)

		assert.Equal(t, "hotel1", orders1[0].HotelId, "incorrect hotel for user %s", userEmail1)
		assert.Equal(t, "hotel2", orders2[0].HotelId, "incorrect hotel for user %s", userEmail2)
	})
}

func TestInMemoryOrder_Concurrency(t *testing.T) {
	t.Run("parallel adding orders for different users", func(t *testing.T) {
		repo := NewInMemoryOrder()
		ctx := context.Background()

		numUsers := 10
		ordersPerUser := 5

		var wg sync.WaitGroup
		wg.Add(numUsers)

		for i := 0; i < numUsers; i++ {
			userEmail := "user" + string(rune('0'+i)) + "@example.com"

			go func(email string) {
				defer wg.Done()

				for j := 0; j < ordersPerUser; j++ {
					order := model.Order{
						UserEmail: email,
						HotelId:   "hotel" + string(rune('0'+j)),
						RoomClass: "standard",
						CheckIn:   time.Date(2023, 7, j+1, 0, 0, 0, 0, time.UTC),
						CheckOut:  time.Date(2023, 7, j+3, 0, 0, 0, 0, time.UTC),
						Rooms:     uint32(j + 1),
					}

					err := repo.AddOrder(ctx, []model.Order{order})
					assert.NoError(t, err, "error adding order")
				}
			}(userEmail)
		}

		wg.Wait()

		repo.ordersLock.RLock()
		defer repo.ordersLock.RUnlock()

		assert.Len(t, repo.usersOrders, numUsers, "expected %d users, got: %d", numUsers, len(repo.usersOrders))

		for i := 0; i < numUsers; i++ {
			userEmail := "user" + string(rune('0'+i)) + "@example.com"
			orders, exists := repo.usersOrders[userEmail]

			assert.True(t, exists, "orders not found for user %s", userEmail)
			assert.Len(t, orders, ordersPerUser, "expected %d orders for user %s, got: %d", ordersPerUser, userEmail, len(orders))
		}
	})

	t.Run("parallel adding orders for one user", func(t *testing.T) {
		repo := NewInMemoryOrder()
		ctx := context.Background()

		userEmail := "test@example.com"
		numGoroutines := 10
		ordersPerGoroutine := 5

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()

				for j := 0; j < ordersPerGoroutine; j++ {
					order := model.Order{
						UserEmail: userEmail,
						HotelId:   "hotel" + string(rune('0'+routineID)) + string(rune('0'+j)),
						RoomClass: "standard",
						CheckIn:   time.Date(2023, 7, j+1, 0, 0, 0, 0, time.UTC),
						CheckOut:  time.Date(2023, 7, j+3, 0, 0, 0, 0, time.UTC),
						Rooms:     uint32(j + 1),
					}

					err := repo.AddOrder(ctx, []model.Order{order})
					assert.NoError(t, err, "error adding order")
				}
			}(i)
		}

		wg.Wait()

		repo.ordersLock.RLock()
		defer repo.ordersLock.RUnlock()

		orders, exists := repo.usersOrders[userEmail]
		require.True(t, exists, "orders not found for user %s", userEmail)

		expectedOrders := numGoroutines * ordersPerGoroutine
		assert.Len(t, orders, expectedOrders, "expected %d orders, got: %d", expectedOrders, len(orders))

		hotelIds := make(map[string]bool)
		for _, order := range orders {
			assert.False(t, hotelIds[order.HotelId], "duplicate order for hotel %s", order.HotelId)
			hotelIds[order.HotelId] = true
		}
	})

	t.Run("parallel adding multiple orders in one call", func(t *testing.T) {
		repo := NewInMemoryOrder()
		ctx := context.Background()

		numGoroutines := 10
		ordersPerBatch := 5

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()

				userEmail := "user" + string(rune('0'+routineID)) + "@example.com"
				orders := make([]model.Order, ordersPerBatch)

				for j := 0; j < ordersPerBatch; j++ {
					orders[j] = model.Order{
						UserEmail: userEmail,
						HotelId:   "hotel" + string(rune('0'+j)),
						RoomClass: "standard",
						CheckIn:   time.Date(2023, 7, j+1, 0, 0, 0, 0, time.UTC),
						CheckOut:  time.Date(2023, 7, j+3, 0, 0, 0, 0, time.UTC),
						Rooms:     uint32(j + 1),
					}
				}

				err := repo.AddOrder(ctx, orders)
				assert.NoError(t, err, "error adding orders batch")
			}(i)
		}

		wg.Wait()

		repo.ordersLock.RLock()
		defer repo.ordersLock.RUnlock()

		assert.Len(t, repo.usersOrders, numGoroutines, "expected %d users, got: %d", numGoroutines, len(repo.usersOrders))

		for i := 0; i < numGoroutines; i++ {
			userEmail := "user" + string(rune('0'+i)) + "@example.com"
			orders, exists := repo.usersOrders[userEmail]

			assert.True(t, exists, "orders not found for user %s", userEmail)
			assert.Len(t, orders, ordersPerBatch, "expected %d orders for user %s, got: %d", ordersPerBatch, userEmail, len(orders))
		}
	})
}
