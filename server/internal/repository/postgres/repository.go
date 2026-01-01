package postgres

import (
	"context"
	"time"

	"github.com/flashtix/server/db"
	"github.com/flashtix/server/internal/domain"
)

type eventRepository struct {
	client *db.PrismaClient
}

func NewEventRepository(client *db.PrismaClient) domain.EventRepository {
	return &eventRepository{client: client}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	_, err := r.client.Event.CreateOne(
		db.Event.Name.Set(event.Name),
		db.Event.Description.Set(event.Description),
		db.Event.Date.Set(event.Date),
		db.Event.Venue.Set(event.Venue),
		db.Event.Capacity.Set(event.Capacity),
	).Exec(ctx)
	return err
}

func (r *eventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	event, err := r.client.Event.FindUnique(
		db.Event.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.Event{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		Venue:       event.Venue,
		Capacity:    event.Capacity,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}, nil
}

func (r *eventRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	events, err := r.client.Event.FindMany().Exec(ctx)
	if err != nil {
		return nil, err
	}

	var result []*domain.Event
	for _, event := range events {
		result = append(result, &domain.Event{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Date:        event.Date,
			Venue:       event.Venue,
			Capacity:    event.Capacity,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		})
	}
	return result, nil
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	_, err := r.client.Event.FindUnique(
		db.Event.ID.Equals(event.ID),
	).Update(
		db.Event.Name.Set(event.Name),
		db.Event.Description.Set(event.Description),
		db.Event.Date.Set(event.Date),
		db.Event.Venue.Set(event.Venue),
		db.Event.Capacity.Set(event.Capacity),
	).Exec(ctx)
	return err
}

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Event.FindUnique(
		db.Event.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

type ticketRepository struct {
	client *db.PrismaClient
}

func NewTicketRepository(client *db.PrismaClient) domain.TicketRepository {
	return &ticketRepository{client: client}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	status := db.TicketStatusAvailable
	switch ticket.Status {
	case "reserved":
		status = db.TicketStatusReserved
	case "sold":
		status = db.TicketStatusSold
	}

	params := []db.TicketSetParam{
		db.Ticket.ID.Set(ticket.ID),
		db.Ticket.EventID.Set(ticket.EventID),
		db.Ticket.Status.Set(status),
		db.Ticket.Price.Set(ticket.Price),
	}

	if ticket.UserID != "" {
		params = append(params, db.Ticket.UserID.SetOptional(&ticket.UserID))
	}
	if ticket.ReservedUntil != nil {
		params = append(params, db.Ticket.ReservedUntil.SetOptional(ticket.ReservedUntil))
	}

	_, err := r.client.Ticket.CreateOne(
		db.Ticket.Seat.Set(ticket.Seat),
		db.Ticket.Event.Link(db.Event.ID.Equals(ticket.EventID)),
		params...,
	).Exec(ctx)
	return err
}

func (r *ticketRepository) GetByID(ctx context.Context, id string) (*domain.Ticket, error) {
	ticket, err := r.client.Ticket.FindUnique(
		db.Ticket.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	status := "available"
	switch ticket.Status {
	case db.TicketStatusReserved:
		status = "reserved"
	case db.TicketStatusSold:
		status = "sold"
	}

	userID, _ := ticket.UserID()
	reservedUntil, reservedUntilOk := ticket.ReservedUntil()

	var reservedUntilPtr *time.Time
	if reservedUntilOk {
		t := time.Time(reservedUntil)
		reservedUntilPtr = &t
	}

	return &domain.Ticket{
		ID:            ticket.ID,
		EventID:       ticket.EventID,
		UserID:        userID,
		Seat:          ticket.Seat,
		Status:        status,
		Price:         ticket.Price,
		ReservedUntil: reservedUntilPtr,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}, nil
}

func (r *ticketRepository) GetByEventID(ctx context.Context, eventID string) ([]*domain.Ticket, error) {
	tickets, err := r.client.Ticket.FindMany(
		db.Ticket.EventID.Equals(eventID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	var result []*domain.Ticket
	for _, ticket := range tickets {
		status := "available"
		switch ticket.Status {
		case db.TicketStatusReserved:
			status = "reserved"
		case db.TicketStatusSold:
			status = "sold"
		}

		userID, _ := ticket.UserID()
		reservedUntil, reservedUntilOk := ticket.ReservedUntil()

		var reservedUntilPtr *time.Time
		if reservedUntilOk {
			t := time.Time(reservedUntil)
			reservedUntilPtr = &t
		}

		result = append(result, &domain.Ticket{
			ID:            ticket.ID,
			EventID:       ticket.EventID,
			UserID:        userID,
			Seat:          ticket.Seat,
			Status:        status,
			Price:         ticket.Price,
			ReservedUntil: reservedUntilPtr,
			CreatedAt:     ticket.CreatedAt,
			UpdatedAt:     ticket.UpdatedAt,
		})
	}
	return result, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	status := db.TicketStatusAvailable
	switch ticket.Status {
	case "reserved":
		status = db.TicketStatusReserved
	case "sold":
		status = db.TicketStatusSold
	}

	params := []db.TicketSetParam{
		db.Ticket.Status.Set(status),
		db.Ticket.Price.Set(ticket.Price),
	}

	if ticket.UserID != "" {
		params = append(params, db.Ticket.UserID.SetOptional(&ticket.UserID))
	}
	if ticket.ReservedUntil != nil {
		params = append(params, db.Ticket.ReservedUntil.SetOptional(ticket.ReservedUntil))
	}

	_, err := r.client.Ticket.FindUnique(
		db.Ticket.ID.Equals(ticket.ID),
	).Update(params...).Exec(ctx)
	return err
}

func (r *ticketRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Ticket.FindUnique(
		db.Ticket.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func (r *ticketRepository) ReserveSeat(ctx context.Context, eventID, seat string, userID string, duration time.Duration) error {
	reservedUntil := time.Now().Add(duration)
	_, err := r.client.Ticket.FindMany(
		db.Ticket.EventID.Equals(eventID),
		db.Ticket.Seat.Equals(seat),
		db.Ticket.Status.Equals(db.TicketStatusAvailable),
	).Update(
		db.Ticket.UserID.SetOptional(&userID),
		db.Ticket.Status.Set(db.TicketStatusReserved),
		db.Ticket.ReservedUntil.SetOptional(&reservedUntil),
	).Exec(ctx)
	return err
}

func (r *ticketRepository) ReleaseSeat(ctx context.Context, eventID, seat string) error {
	_, err := r.client.Ticket.FindMany(
		db.Ticket.EventID.Equals(eventID),
		db.Ticket.Seat.Equals(seat),
	).Update(
		db.Ticket.UserID.SetOptional(nil),
		db.Ticket.Status.Set(db.TicketStatusAvailable),
		db.Ticket.ReservedUntil.SetOptional(nil),
	).Exec(ctx)
	return err
}

type userRepository struct {
	client *db.PrismaClient
}

func NewUserRepository(client *db.PrismaClient) domain.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.client.User.CreateOne(
		db.User.Email.Set(user.Email),
		db.User.Name.Set(user.Name),
	).Exec(ctx)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.client.User.FindUnique(
		db.User.ID.Equals(user.ID),
	).Update(
		db.User.Email.Set(user.Email),
		db.User.Name.Set(user.Name),
	).Exec(ctx)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}
