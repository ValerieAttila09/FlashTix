package handlers

import (
	"net/http"

	"github.com/flashtix/server/internal/domain"
	"github.com/flashtix/server/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketService *services.TicketService
}

func NewTicketHandler(ticketService *services.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) ReserveSeat(c *gin.Context) {
	var req struct {
		EventID string `json:"event_id" binding:"required"`
		Seat    string `json:"seat" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id") // from auth middleware
	err := h.ticketService.ReserveSeat(c.Request.Context(), req.EventID, req.Seat, userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat reserved successfully"})
}

func (h *TicketHandler) ConfirmPurchase(c *gin.Context) {
	var req struct {
		EventID string `json:"event_id" binding:"required"`
		Seat    string `json:"seat" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	err := h.ticketService.ConfirmPurchase(c.Request.Context(), req.EventID, req.Seat, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase confirmed"})
}

type EventHandler struct {
	eventRepo domain.EventRepository
}

func NewEventHandler(eventRepo domain.EventRepository) *EventHandler {
	return &EventHandler{eventRepo: eventRepo}
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	events, err := h.eventRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event domain.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = uuid.New().String()
	if err := h.eventRepo.Create(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}
