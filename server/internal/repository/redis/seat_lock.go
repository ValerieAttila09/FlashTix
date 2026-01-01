package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type SeatLockRepository struct {
	client *redis.Client
}

func NewSeatLockRepository(client *redis.Client) *SeatLockRepository {
	return &SeatLockRepository{client: client}
}

// LockSeat locks a seat for a user with expiration
func (r *SeatLockRepository) LockSeat(ctx context.Context, eventID, seat string, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	return r.client.Set(ctx, key, userID, expiration).Err()
}

// UnlockSeat unlocks a seat
func (r *SeatLockRepository) UnlockSeat(ctx context.Context, eventID, seat string) error {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	return r.client.Del(ctx, key).Err()
}

// IsSeatLocked checks if a seat is locked and by whom
func (r *SeatLockRepository) IsSeatLocked(ctx context.Context, eventID, seat string) (string, error) {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	userID, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // not locked
	}
	return userID, err
}

// ExtendLock extends the lock expiration
func (r *SeatLockRepository) ExtendLock(ctx context.Context, eventID, seat string, expiration time.Duration) error {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	return r.client.Expire(ctx, key, expiration).Err()
}
