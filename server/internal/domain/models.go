package domain

import (
	"context"
	"time"
)

// Event represents an event entity
type Event struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Venue       string    `json:"venue"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Ticket represents a ticket entity
type Ticket struct {
	ID            string     `json:"id" gorm:"primaryKey"`
	EventID       string     `json:"event_id"`
	UserID        string     `json:"user_id"`
	Seat          string     `json:"seat"`
	Status        string     `json:"status"` // available, reserved, sold
	Price         float64    `json:"price"`
	ReservedUntil *time.Time `json:"reserved_until"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// User represents a user entity
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EventRepository interface
type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	GetByID(ctx context.Context, id string) (*Event, error)
	GetAll(ctx context.Context) ([]*Event, error)
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id string) error
}

// TicketRepository interface
type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	GetByID(ctx context.Context, id string) (*Ticket, error)
	GetByEventID(ctx context.Context, eventID string) ([]*Ticket, error)
	Update(ctx context.Context, ticket *Ticket) error
	Delete(ctx context.Context, id string) error
	ReserveSeat(ctx context.Context, eventID, seat string, userID string, duration time.Duration) error
	ReleaseSeat(ctx context.Context, eventID, seat string) error
}

// UserRepository interface
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
