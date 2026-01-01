package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestUpstashRedisClient(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Logf("No .env file found: %v", err)
	}

	// Skip test if environment variables are not set
	url := os.Getenv("UPSTASH_REDIS_REST_URL")
	token := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if url == "" || token == "" {
		t.Skip("UPSTASH_REDIS_REST_URL and UPSTASH_REDIS_REST_TOKEN not set")
	}

	client := &UpstashRedisClient{
		url:   url,
		token: token,
	}

	ctx := context.Background()
	testKey := "test_key_" + time.Now().Format("20060102150405")
	testValue := "test_value"

	// Test SET
	t.Run("SET", func(t *testing.T) {
		err := client.set(ctx, testKey, testValue, 60*time.Second)
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}
	})

	// Test GET
	t.Run("GET", func(t *testing.T) {
		value, err := client.get(ctx, testKey)
		if err != nil {
			t.Fatalf("Failed to get key: %v", err)
		}
		if value != testValue {
			t.Errorf("Expected %s, got %s", testValue, value)
		}
	})

	// Test DEL
	t.Run("DEL", func(t *testing.T) {
		err := client.del(ctx, testKey)
		if err != nil {
			t.Fatalf("Failed to delete key: %v", err)
		}
	})

	// Verify key is deleted
	t.Run("GET after DEL", func(t *testing.T) {
		_, err := client.get(ctx, testKey)
		if err == nil || err.Error() != "key not found" {
			t.Errorf("Expected key not found error, got: %v", err)
		}
	})
}

func TestSeatLockRepository(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Logf("No .env file found: %v", err)
	}

	// Skip test if environment variables are not set
	url := os.Getenv("UPSTASH_REDIS_REST_URL")
	token := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if url == "" || token == "" {
		t.Skip("UPSTASH_REDIS_REST_URL and UPSTASH_REDIS_REST_TOKEN not set")
	}

	repo := NewSeatLockRepository(url, token)
	ctx := context.Background()

	eventID := "test_event_" + time.Now().Format("20060102150405")
	seat := "A1"
	userID := "user123"

	// Test LockSeat
	t.Run("LockSeat", func(t *testing.T) {
		err := repo.LockSeat(ctx, eventID, seat, userID, 30*time.Second)
		if err != nil {
			t.Fatalf("Failed to lock seat: %v", err)
		}
	})

	// Test IsSeatLocked
	t.Run("IsSeatLocked", func(t *testing.T) {
		lockedUserID, err := repo.IsSeatLocked(ctx, eventID, seat)
		if err != nil {
			t.Fatalf("Failed to check seat lock: %v", err)
		}
		if lockedUserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, lockedUserID)
		}
	})

	// Test UnlockSeat
	t.Run("UnlockSeat", func(t *testing.T) {
		err := repo.UnlockSeat(ctx, eventID, seat)
		if err != nil {
			t.Fatalf("Failed to unlock seat: %v", err)
		}
	})

	// Verify seat is unlocked
	t.Run("IsSeatLocked after unlock", func(t *testing.T) {
		lockedUserID, err := repo.IsSeatLocked(ctx, eventID, seat)
		if err != nil {
			t.Fatalf("Failed to check seat lock: %v", err)
		}
		if lockedUserID != "" {
			t.Errorf("Expected empty user ID after unlock, got %s", lockedUserID)
		}
	})
}
