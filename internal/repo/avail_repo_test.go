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

func TestInMemoryAvailability_BookRooms(t *testing.T) {
	t.Run("successful booking", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID := "hotel1"
		roomClass := "standard"
		bookingKey := BookingKey{HotelID: hotelID, RoomClass: roomClass}

		availMap := make(map[time.Time]uint32)
		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap[checkIn] = 5
		availMap[checkIn.AddDate(0, 0, 1)] = 5

		repo.roomAvailabilityMap.Store(bookingKey, availMap)

		order := model.Order{
			HotelId:   hotelID,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     2,
		}

		err := repo.BookRooms(context.Background(), []model.Order{order})
		assert.NoError(t, err, "Failed to book rooms, expected success")

		val, ok := repo.roomAvailabilityMap.Load(bookingKey)
		require.True(t, ok, "Failed to find availability data after booking")

		updatedAvailMap := val.(map[time.Time]uint32)
		assert.Equal(t, uint32(3), updatedAvailMap[checkIn], "Expected 3 available rooms on %v", checkIn)
		assert.Equal(t, uint32(3), updatedAvailMap[checkIn.AddDate(0, 0, 1)], "Expected 3 available rooms on %v", checkIn.AddDate(0, 0, 1))
	})

	t.Run("not enough rooms for booking", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID := "hotel1"
		roomClass := "standard"
		bookingKey := BookingKey{HotelID: hotelID, RoomClass: roomClass}

		availMap := make(map[time.Time]uint32)
		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap[checkIn] = 2
		availMap[checkIn.AddDate(0, 0, 1)] = 2

		repo.roomAvailabilityMap.Store(bookingKey, availMap)

		order := model.Order{
			HotelId:   hotelID,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     3,
		}

		err := repo.BookRooms(context.Background(), []model.Order{order})
		assert.Equal(t, model.ErrInsufficientQuota, err, "Expected ErrInsufficientQuota")

		val, ok := repo.roomAvailabilityMap.Load(bookingKey)
		require.True(t, ok, "Failed to find availability data after booking attempt")

		updatedAvailMap := val.(map[time.Time]uint32)
		assert.Equal(t, uint32(2), updatedAvailMap[checkIn], "Expected 2 available rooms on %v", checkIn)
	})

	t.Run("booking multiple orders", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID1 := "hotel1"
		hotelID2 := "hotel2"
		roomClass := "standard"

		bookingKey1 := BookingKey{HotelID: hotelID1, RoomClass: roomClass}
		bookingKey2 := BookingKey{HotelID: hotelID2, RoomClass: roomClass}

		availMap1 := make(map[time.Time]uint32)
		availMap2 := make(map[time.Time]uint32)

		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap1[checkIn] = 5
		availMap1[checkIn.AddDate(0, 0, 1)] = 5

		availMap2[checkIn] = 3
		availMap2[checkIn.AddDate(0, 0, 1)] = 3

		repo.roomAvailabilityMap.Store(bookingKey1, availMap1)
		repo.roomAvailabilityMap.Store(bookingKey2, availMap2)

		order1 := model.Order{
			HotelId:   hotelID1,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     2,
		}

		order2 := model.Order{
			HotelId:   hotelID2,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     1,
		}

		err := repo.BookRooms(context.Background(), []model.Order{order1, order2})
		assert.NoError(t, err, "Expected successful booking")

		val1, ok := repo.roomAvailabilityMap.Load(bookingKey1)
		require.True(t, ok, "Failed to find availability data for hotel1")

		val2, ok := repo.roomAvailabilityMap.Load(bookingKey2)
		require.True(t, ok, "Failed to find availability data for hotel2")

		updatedAvailMap1 := val1.(map[time.Time]uint32)
		updatedAvailMap2 := val2.(map[time.Time]uint32)

		assert.Equal(t, uint32(3), updatedAvailMap1[checkIn], "Expected 3 available rooms for hotel1")
		assert.Equal(t, uint32(2), updatedAvailMap2[checkIn], "Expected 2 available rooms for hotel2")
	})

	t.Run("booking fails when one order has insufficient quota", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID1 := "hotel1"
		hotelID2 := "hotel2"
		roomClass := "standard"

		bookingKey1 := BookingKey{HotelID: hotelID1, RoomClass: roomClass}
		bookingKey2 := BookingKey{HotelID: hotelID2, RoomClass: roomClass}

		availMap1 := make(map[time.Time]uint32)
		availMap2 := make(map[time.Time]uint32)

		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap1[checkIn] = 5
		availMap1[checkIn.AddDate(0, 0, 1)] = 5

		availMap2[checkIn] = 2
		availMap2[checkIn.AddDate(0, 0, 1)] = 2

		repo.roomAvailabilityMap.Store(bookingKey1, availMap1)
		repo.roomAvailabilityMap.Store(bookingKey2, availMap2)

		order1 := model.Order{
			HotelId:   hotelID1,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     2,
		}

		order2 := model.Order{
			HotelId:   hotelID2,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     3,
		}

		err := repo.BookRooms(context.Background(), []model.Order{order1, order2})
		assert.Equal(t, model.ErrInsufficientQuota, err, "Expected ErrInsufficientQuota")

		val1, ok := repo.roomAvailabilityMap.Load(bookingKey1)
		require.True(t, ok, "Failed to find availability data for hotel1")

		val2, ok := repo.roomAvailabilityMap.Load(bookingKey2)
		require.True(t, ok, "Failed to find availability data for hotel2")

		updatedAvailMap1 := val1.(map[time.Time]uint32)
		updatedAvailMap2 := val2.(map[time.Time]uint32)

		assert.Equal(t, uint32(5), updatedAvailMap1[checkIn], "Expected 5 available rooms for hotel1")
		assert.Equal(t, uint32(2), updatedAvailMap2[checkIn], "Expected 2 available rooms for hotel2")

		assert.Equal(t, uint32(2), updatedAvailMap2[checkIn], "Expected 2 available rooms for hotel2")
	})
}

func TestInMemoryAvailability_Concurrency(t *testing.T) {
	t.Run("concurrent booking different hotels", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID1 := "hotel1"
		hotelID2 := "hotel2"
		roomClass := "standard"

		bookingKey1 := BookingKey{HotelID: hotelID1, RoomClass: roomClass}
		bookingKey2 := BookingKey{HotelID: hotelID2, RoomClass: roomClass}

		availMap1 := make(map[time.Time]uint32)
		availMap2 := make(map[time.Time]uint32)

		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap1[checkIn] = 10
		availMap1[checkIn.AddDate(0, 0, 1)] = 10

		availMap2[checkIn] = 10
		availMap2[checkIn.AddDate(0, 0, 1)] = 10

		repo.roomAvailabilityMap.Store(bookingKey1, availMap1)
		repo.roomAvailabilityMap.Store(bookingKey2, availMap2)

		numGoroutines := 9
		var wg sync.WaitGroup
		wg.Add(numGoroutines * 2)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				order := model.Order{
					HotelId:   hotelID1,
					RoomClass: roomClass,
					CheckIn:   checkIn,
					CheckOut:  checkOut,
					Rooms:     1,
				}

				err := repo.BookRooms(context.Background(), []model.Order{order})
				assert.NoError(t, err, "Error booking hotel1")
			}()
		}

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				order := model.Order{
					HotelId:   hotelID2,
					RoomClass: roomClass,
					CheckIn:   checkIn,
					CheckOut:  checkOut,
					Rooms:     1,
				}

				err := repo.BookRooms(context.Background(), []model.Order{order})
				assert.NoError(t, err, "Error booking hotel2")
			}()
		}

		wg.Wait()

		val1, _ := repo.roomAvailabilityMap.Load(bookingKey1)
		val2, _ := repo.roomAvailabilityMap.Load(bookingKey2)

		updatedAvailMap1 := val1.(map[time.Time]uint32)
		updatedAvailMap2 := val2.(map[time.Time]uint32)

		assert.Equal(t, uint32(1), updatedAvailMap1[checkIn], "Expected 1 available room for hotel1")
		assert.Equal(t, uint32(1), updatedAvailMap2[checkIn], "Expected 1 available room for hotel2")
	})

	t.Run("concurrent booking one hotel", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID := "hotel1"
		roomClass := "standard"
		bookingKey := BookingKey{HotelID: hotelID, RoomClass: roomClass}

		availMap := make(map[time.Time]uint32)
		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap[checkIn] = 5
		availMap[checkIn.AddDate(0, 0, 1)] = 5

		repo.roomAvailabilityMap.Store(bookingKey, availMap)

		numGoroutines := 10
		successCount := 0
		var mu sync.Mutex
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				order := model.Order{
					HotelId:   hotelID,
					RoomClass: roomClass,
					CheckIn:   checkIn,
					CheckOut:  checkOut,
					Rooms:     1,
				}

				err := repo.BookRooms(context.Background(), []model.Order{order})
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		assert.Equal(t, 5, successCount, "Expected 5 successful bookings")

		val, _ := repo.roomAvailabilityMap.Load(bookingKey)
		updatedAvailMap := val.(map[time.Time]uint32)

		assert.Equal(t, uint32(0), updatedAvailMap[checkIn], "Expected 0 available rooms")
	})

	t.Run("concurrent booking with different dates", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID := "hotel1"
		roomClass := "standard"
		bookingKey := BookingKey{HotelID: hotelID, RoomClass: roomClass}

		availMap := make(map[time.Time]uint32)
		baseDate := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)

		for i := 0; i < 10; i++ {
			date := baseDate.AddDate(0, 0, i)
			availMap[date] = 5
		}

		repo.roomAvailabilityMap.Store(bookingKey, availMap)

		var wg sync.WaitGroup
		wg.Add(5)

		for i := 0; i < 5; i++ {
			checkIn := baseDate.AddDate(0, 0, i*2)
			checkOut := checkIn.AddDate(0, 0, 2)

			go func(in, out time.Time) {
				defer wg.Done()
				order := model.Order{
					HotelId:   hotelID,
					RoomClass: roomClass,
					CheckIn:   in,
					CheckOut:  out,
					Rooms:     1,
				}

				err := repo.BookRooms(context.Background(), []model.Order{order})
				assert.NoError(t, err, "Error booking")
			}(checkIn, checkOut)
		}

		wg.Wait()

		val, _ := repo.roomAvailabilityMap.Load(bookingKey)
		updatedAvailMap := val.(map[time.Time]uint32)

		expectedAvailability := make(map[time.Time]uint32)
		for i := 0; i < 10; i++ {
			date := baseDate.AddDate(0, 0, i)
			expectedAvailability[date] = 5
		}

		for i := 0; i < 5; i++ {
			checkIn := baseDate.AddDate(0, 0, i*2)
			checkOut := checkIn.AddDate(0, 0, 2)

			for j := 0; j < 2; j++ {
				date := checkIn.AddDate(0, 0, j)
				if date.Before(checkOut) {
					expectedAvailability[date]--
				}
			}
		}

		for i := 0; i < 10; i++ {
			date := baseDate.AddDate(0, 0, i)
			expected := expectedAvailability[date]

			assert.Equal(t, expected, updatedAvailMap[date], "For date %v expected %d available rooms", date, expected)
		}
	})
}

func TestInMemoryAvailability_CheckRoomsAvailability(t *testing.T) {
	t.Run("rooms are available", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID := "hotel1"
		roomClass := "standard"
		bookingKey := BookingKey{HotelID: hotelID, RoomClass: roomClass}

		availMap := make(map[time.Time]uint32)
		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap[checkIn] = 5
		availMap[checkIn.AddDate(0, 0, 1)] = 5

		repo.roomAvailabilityMap.Store(bookingKey, availMap)

		order := model.Order{
			HotelId:   hotelID,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     3,
		}

		available, err := repo.CheckRoomsAvailability(context.Background(), []model.Order{order})
		assert.NoError(t, err, "Expected no error")
		assert.True(t, available, "Expected rooms to be available")
	})

	t.Run("rooms are not available", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID := "hotel1"
		roomClass := "standard"
		bookingKey := BookingKey{HotelID: hotelID, RoomClass: roomClass}

		availMap := make(map[time.Time]uint32)
		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap[checkIn] = 2
		availMap[checkIn.AddDate(0, 0, 1)] = 2

		repo.roomAvailabilityMap.Store(bookingKey, availMap)

		order := model.Order{
			HotelId:   hotelID,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     3,
		}

		available, err := repo.CheckRoomsAvailability(context.Background(), []model.Order{order})
		assert.NoError(t, err, "Expected no error")
		assert.False(t, available, "Expected rooms to be not available")
	})

	t.Run("checking multiple orders", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID1 := "hotel1"
		hotelID2 := "hotel2"
		roomClass := "standard"

		bookingKey1 := BookingKey{HotelID: hotelID1, RoomClass: roomClass}
		bookingKey2 := BookingKey{HotelID: hotelID2, RoomClass: roomClass}

		availMap1 := make(map[time.Time]uint32)
		availMap2 := make(map[time.Time]uint32)

		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap1[checkIn] = 5
		availMap1[checkIn.AddDate(0, 0, 1)] = 5

		availMap2[checkIn] = 3
		availMap2[checkIn.AddDate(0, 0, 1)] = 3

		repo.roomAvailabilityMap.Store(bookingKey1, availMap1)
		repo.roomAvailabilityMap.Store(bookingKey2, availMap2)

		order1 := model.Order{
			HotelId:   hotelID1,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     2,
		}

		order2 := model.Order{
			HotelId:   hotelID2,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     2,
		}

		available, err := repo.CheckRoomsAvailability(context.Background(), []model.Order{order1, order2})
		assert.NoError(t, err, "Expected no error")
		assert.True(t, available, "Expected rooms to be available for both orders")

		order3 := model.Order{
			HotelId:   hotelID1,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     2,
		}

		order4 := model.Order{
			HotelId:   hotelID2,
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     4,
		}

		available, err = repo.CheckRoomsAvailability(context.Background(), []model.Order{order3, order4})
		assert.NoError(t, err, "Expected no error")
		assert.False(t, available, "Expected rooms to be not available for one of the orders")
	})

	t.Run("checking non-existent hotel or room class", func(t *testing.T) {
		repo := NewinMemoryAvailability()

		hotelID := "hotel1"
		roomClass := "standard"
		bookingKey := BookingKey{HotelID: hotelID, RoomClass: roomClass}

		availMap := make(map[time.Time]uint32)
		checkIn := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
		checkOut := time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)

		availMap[checkIn] = 5
		availMap[checkIn.AddDate(0, 0, 1)] = 5

		repo.roomAvailabilityMap.Store(bookingKey, availMap)

		order1 := model.Order{
			HotelId:   "nonexistent_hotel",
			RoomClass: roomClass,
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     1,
		}

		available, err := repo.CheckRoomsAvailability(context.Background(), []model.Order{order1})
		assert.Equal(t, model.ErrNotFound, err, "Expected ErrNotFound error")
		assert.False(t, available, "Expected rooms to be not available for non-existent hotel")

		order2 := model.Order{
			HotelId:   hotelID,
			RoomClass: "nonexistent_class",
			CheckIn:   checkIn,
			CheckOut:  checkOut,
			Rooms:     1,
		}

		available, err = repo.CheckRoomsAvailability(context.Background(), []model.Order{order2})
		assert.Equal(t, model.ErrNotFound, err, "Expected ErrNotFound error")
		assert.False(t, available, "Expected rooms to be not available for non-existent room class")
	})
}
