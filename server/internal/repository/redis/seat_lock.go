package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UpstashRedisClient struct {
	url   string
	token string
}

type UpstashResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

type SeatLockRepository struct {
	client *UpstashRedisClient
}

func NewSeatLockRepository(url, token string) *SeatLockRepository {
	return &SeatLockRepository{
		client: &UpstashRedisClient{
			url:   url,
			token: token,
		},
	}
}

// LockSeat locks a seat for a user with expiration
func (r *SeatLockRepository) LockSeat(ctx context.Context, eventID, seat string, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	return r.client.set(ctx, key, userID, expiration)
}

// UnlockSeat unlocks a seat
func (r *SeatLockRepository) UnlockSeat(ctx context.Context, eventID, seat string) error {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	return r.client.del(ctx, key)
}

// IsSeatLocked checks if a seat is locked and by whom
func (r *SeatLockRepository) IsSeatLocked(ctx context.Context, eventID, seat string) (string, error) {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	userID, err := r.client.get(ctx, key)
	if err != nil && err.Error() == "key not found" {
		return "", nil // not locked
	}
	return userID, err
}

// ExtendLock extends the lock expiration
func (r *SeatLockRepository) ExtendLock(ctx context.Context, eventID, seat string, expiration time.Duration) error {
	key := fmt.Sprintf("seat_lock:%s:%s", eventID, seat)
	return r.client.expire(ctx, key, expiration)
}

// UpstashRedisClient methods
func (c *UpstashRedisClient) set(ctx context.Context, key, value string, expiration time.Duration) error {
	var url string
	if expiration > 0 {
		// SET key value EX seconds
		url = fmt.Sprintf("%s/set/%s/%s/ex/%d", c.url, key, value, int(expiration.Seconds()))
	} else {
		// SET key value
		url = fmt.Sprintf("%s/set/%s/%s", c.url, key, value)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upstash request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var upstashResp UpstashResponse
	if err := json.NewDecoder(resp.Body).Decode(&upstashResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if upstashResp.Error != "" {
		return fmt.Errorf("upstash error: %s", upstashResp.Error)
	}

	return nil
}

func (c *UpstashRedisClient) get(ctx context.Context, key string) (string, error) {
	url := fmt.Sprintf("%s/get/%s", c.url, key)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var upstashResp UpstashResponse
	if err := json.NewDecoder(resp.Body).Decode(&upstashResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if upstashResp.Error != "" {
		return "", fmt.Errorf("upstash error: %s", upstashResp.Error)
	}

	// For GET, result can be string or null
	if upstashResp.Result == nil {
		return "", fmt.Errorf("key not found")
	}

	result, ok := upstashResp.Result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type")
	}

	return result, nil
}

func (c *UpstashRedisClient) del(ctx context.Context, key string) error {
	url := fmt.Sprintf("%s/del/%s", c.url, key)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upstash request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var upstashResp UpstashResponse
	if err := json.NewDecoder(resp.Body).Decode(&upstashResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if upstashResp.Error != "" {
		return fmt.Errorf("upstash error: %s", upstashResp.Error)
	}

	return nil
}

func (c *UpstashRedisClient) expire(ctx context.Context, key string, expiration time.Duration) error {
	url := fmt.Sprintf("%s/expire/%s/%d", c.url, key, int(expiration.Seconds()))

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upstash request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var upstashResp UpstashResponse
	if err := json.NewDecoder(resp.Body).Decode(&upstashResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if upstashResp.Error != "" {
		return fmt.Errorf("upstash error: %s", upstashResp.Error)
	}

	return nil
}
