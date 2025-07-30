package main

import (
	"fmt"
	"sync"
	"time"
)

type Booking struct {
	UserId int
	StartDate time.Time
	EndDate time.Time
}

type Hotel struct {
	mu sync.Mutex
	reservations map[int][]Booking
}

func NewHotel() *Hotel {
	return &Hotel{
		reservations: make(map[int][]Booking),
	}
}

func datesOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

func (h *Hotel) BookRoom(roomId, userId int, startDate, endDate time.Time) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	bookings := h.reservations[roomId]
	for _, booking := range bookings {
		if datesOverlap(startDate, endDate, booking.StartDate, booking.EndDate) {
			return fmt.Errorf("room %d already booked for overlapping dates", roomId)
		}
	}

	// No overlap, book the room
	h.reservations[roomId] = append(h.reservations[roomId], Booking{
		UserId: userId,
		StartDate: startDate,
		EndDate: endDate,
	})
	return nil
}

func main() {
	hotel := NewHotel()

	start := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC)

	var wg sync.WaitGroup

	for user := 1; user <= 100; user++ {
		wg.Add(1)
		go func(userId int) {
			defer wg.Done()
			err := hotel.BookRoom(101, userId, start, end)
			if err != nil {
				fmt.Printf("User %d: failed to book room 101: %s\n", userId, err)
			} else {
				fmt.Printf("User %d: successfully booked room 101\n", userId)
			}
		}(user)
	}
	wg.Wait()
}