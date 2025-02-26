package repo

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/pkg/utils"
	"context"
	"sort"
	"sync"
	"time"
)

type BookingKey struct {
	HotelID   string
	RoomClass string
}

type inMemoryAvailability struct {
	hotelLocks sync.Map
	//should be initialized before use
	//roomAvailabilityMap map[BookingIdEntity]map[time.Time]uint32
	roomAvailabilityMap sync.Map
}

func NewinMemoryAvailability() *inMemoryAvailability {
	return &inMemoryAvailability{}
}

func (r *inMemoryAvailability) getHotelLock(hotelId string) *sync.RWMutex {
	lock, _ := r.hotelLocks.LoadOrStore(hotelId, &sync.RWMutex{})
	return lock.(*sync.RWMutex)
}

func (r *inMemoryAvailability) checkAvailability(order model.Order) (bool, map[time.Time]uint32, error) {
	key := BookingKey{HotelID: order.HotelId, RoomClass: order.RoomClass}
	val, ok := r.roomAvailabilityMap.Load(key)
	if !ok {
		return false, nil, model.ErrNotFound
	}
	availMap := val.(map[time.Time]uint32)

	for _, date := range utils.GetDatesExclusive(order.CheckIn, order.CheckOut) {
		dateQuota, ok := availMap[date]
		if !ok || dateQuota < order.Rooms {
			return false, nil, nil
		}
	}

	return true, availMap, nil
}

func (r *inMemoryAvailability) CheckRoomsAvailability(ctx context.Context, orders []model.Order) (bool, error) {
	sortedHotelIds := r.getSortedHotelIds(orders)

	hotelIdsMap := make(map[string]bool)

	//read lock hotels
	for _, hotelId := range sortedHotelIds {
		if !hotelIdsMap[hotelId] {
			hotelLock := r.getHotelLock(hotelId)
			//use read lock
			hotelLock.RLock()
			defer hotelLock.RUnlock()
			hotelIdsMap[hotelId] = true
		}
	}

	//check availability
	for _, order := range orders {
		available, _, err := r.checkAvailability(order)
		if err != nil {
			return false, err
		}
		if !available {
			return false, nil
		}
	}
	return true, nil
}

func (r *inMemoryAvailability) getSortedHotelIds(orders []model.Order) []string {
	hotelIds := make([]string, len(orders))
	for _, order := range orders {
		hotelIds = append(hotelIds, order.HotelId)
	}
	sort.Strings(hotelIds)
	return hotelIds
}

func (r *inMemoryAvailability) BookRooms(ctx context.Context, orders []model.Order) error {
	sortedHotelIds := r.getSortedHotelIds(orders)

	hotelIdsMap := make(map[string]bool)

	//lock hotels
	for _, hotelId := range sortedHotelIds {
		if !hotelIdsMap[hotelId] {
			hotelLock := r.getHotelLock(hotelId)
			//use write lock
			hotelLock.Lock()
			defer hotelLock.Unlock()
			hotelIdsMap[hotelId] = true
		}
	}

	//check availability
	forBooked := make(map[BookingKey]map[time.Time]uint32)
	for _, order := range orders {
		available, availMap, err := r.checkAvailability(order)
		if err != nil {
			return err
		}
		if !available {
			return model.ErrInsufficientQuota
		}
		key := BookingKey{HotelID: order.HotelId, RoomClass: order.RoomClass}
		forBooked[key] = availMap
	}

	//prepare to book rooms
	for _, order := range orders {
		availMap := forBooked[BookingKey{HotelID: order.HotelId, RoomClass: order.RoomClass}]
		for _, date := range utils.GetDatesExclusive(order.CheckIn, order.CheckOut) {
			currentQuota := availMap[date]
			availMap[date] = currentQuota - order.Rooms
		}
	}

	//book rooms
	for bookKey, availMap := range forBooked {
		r.roomAvailabilityMap.Store(bookKey, availMap)
	}

	return nil
}
