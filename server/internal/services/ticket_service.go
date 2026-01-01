package services

import (
	"context"
	"errors"
	"time"

	"github.com/flashtix/server/internal/domain"
	"github.com/flashtix/server/internal/repository/redis"
)

type TicketService struct {
	ticketRepo   domain.TicketRepository
	eventRepo    domain.EventRepository
	seatLockRepo *redis.SeatLockRepository
	lockDuration time.Duration
}

func NewTicketService(ticketRepo domain.TicketRepository, eventRepo domain.EventRepository, seatLockRepo *redis.SeatLockRepository) *TicketService {
	return &TicketService{
		ticketRepo:   ticketRepo,
		eventRepo:    eventRepo,
		seatLockRepo: seatLockRepo,
		lockDuration: 10 * time.Minute, // 10 minutes lock
	}
}

func (s *TicketService) ReserveSeat(ctx context.Context, eventID, seat, userID string) error {
	// Check if seat is already locked
	lockedBy, err := s.seatLockRepo.IsSeatLocked(ctx, eventID, seat)
	if err != nil {
		return err
	}
	if lockedBy != "" && lockedBy != userID {
		return errors.New("seat is already reserved")
	}

	// Lock the seat in Redis
	err = s.seatLockRepo.LockSeat(ctx, eventID, seat, userID, s.lockDuration)
	if err != nil {
		return err
	}

	// Reserve in database
	err = s.ticketRepo.ReserveSeat(ctx, eventID, seat, userID, s.lockDuration)
	if err != nil {
		// Unlock if database update fails
		s.seatLockRepo.UnlockSeat(ctx, eventID, seat)
		return err
	}

	return nil
}

func (s *TicketService) ConfirmPurchase(ctx context.Context, eventID, seat, userID string) error {
	// Check if seat is locked by this user
	lockedBy, err := s.seatLockRepo.IsSeatLocked(ctx, eventID, seat)
	if err != nil {
		return err
	}
	if lockedBy != userID {
		return errors.New("seat not reserved by this user")
	}

	// Update ticket status to sold
	ticket, err := s.ticketRepo.GetByID(ctx, "dummy-id") // This needs proper ID, adjust as needed
	if err != nil {
		return err
	}
	ticket.Status = "sold"
	err = s.ticketRepo.Update(ctx, ticket)
	if err != nil {
		return err
	}

	// Unlock the seat
	return s.seatLockRepo.UnlockSeat(ctx, eventID, seat)
}

func (s *TicketService) ReleaseSeat(ctx context.Context, eventID, seat string) error {
	// Release from database
	err := s.ticketRepo.ReleaseSeat(ctx, eventID, seat)
	if err != nil {
		return err
	}

	// Unlock from Redis
	return s.seatLockRepo.UnlockSeat(ctx, eventID, seat)
}
