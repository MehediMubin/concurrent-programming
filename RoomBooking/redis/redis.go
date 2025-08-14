package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	room = "room-101"
)

func bookRoom(locker *redislock.Client, user string, wg *sync.WaitGroup) {
	defer wg.Done()

	lock, err := locker.Obtain(ctx, room, 5 * time.Second, nil)
	if err == redislock.ErrNotObtained {
		fmt.Printf("%s could not book %s (already locked)\n", user, room)
		return
	} else if err != nil {
		fmt.Printf("Lock error: %v", err)
		return
	}
	defer lock.Release(ctx)

	fmt.Printf("%s successfully booked %s\n", user, room)
	time.Sleep(2 * time.Second)
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer client.Close()

	locker := redislock.New(client)
	var wg sync.WaitGroup

	users := []string{"Mehedi", "Hasan", "Mubin"}
	for _, user := range users {
		wg.Add(1)
		go bookRoom(locker, user, &wg)
	}

	wg.Wait()
}