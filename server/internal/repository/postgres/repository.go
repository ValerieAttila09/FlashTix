package postgres

import (
	"context"
	"time"

	"github.com/flashtix/server/internal/domain"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) domain.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	var events []*domain.Event
	err := r.db.WithContext(ctx).Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Event{}, "id = ?", id).Error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) domain.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) GetByID(ctx context.Context, id string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByEventID(ctx context.Context, eventID string) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).Find(&tickets, "event_id = ?", eventID).Error
	return tickets, err
}

func (r *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ticketRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Ticket{}, "id = ?", id).Error
}

func (r *ticketRepository) ReserveSeat(ctx context.Context, eventID, seat string, userID string, duration time.Duration) error {
	reservedUntil := time.Now().Add(duration)
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ? AND seat = ? AND status = ?", eventID, seat, "available").
		Updates(map[string]interface{}{
			"user_id":        userID,
			"status":         "reserved",
			"reserved_until": &reservedUntil,
		}).Error
}

func (r *ticketRepository) ReleaseSeat(ctx context.Context, eventID, seat string) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ? AND seat = ?", eventID, seat).
		Updates(map[string]interface{}{
			"user_id":        nil,
			"status":         "available",
			"reserved_until": nil,
		}).Error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}
