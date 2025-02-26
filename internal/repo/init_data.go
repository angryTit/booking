package repo

import (
	"log"
	"time"
)

func (r *inMemoryAvailability) InitTestData() {
	hotelID := "reddison"
	roomClass := "lux"

	key := BookingKey{HotelID: hotelID, RoomClass: roomClass}

	availMap := make(map[time.Time]uint32)

	now := time.Now()
	for i := 0; i < 30; i++ {
		date := time.Date(now.Year(), now.Month(), now.Day()+i, 0, 0, 0, 0, time.UTC)
		var quota uint32 = 5
		if i%5 == 0 {
			quota = 2
		}
		if i%7 == 0 {
			quota = 1
		}
		if i%10 == 0 {
			quota = 0
		}
		availMap[date] = quota
	}

	r.roomAvailabilityMap.Store(key, availMap)

	log.Printf("InitTestData: %v\n", availMap)
}
