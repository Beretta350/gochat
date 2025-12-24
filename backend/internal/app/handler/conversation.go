package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/internal/app/repository"
	"github.com/Beretta350/gochat/pkg/logger"
)

// ConversationHandler handles conversation endpoints
type ConversationHandler struct {
	convRepo repository.ConversationRepository
	userRepo repository.UserRepository
	msgRepo  repository.MessageRepository
}

// NewConversationHandler creates a new conversation handler (Fx provider)
func NewConversationHandler(
	convRepo repository.ConversationRepository,
	userRepo repository.UserRepository,
	msgRepo repository.MessageRepository,
) *ConversationHandler {
	logger.Info("Conversation handler initialized")
	return &ConversationHandler{
		convRepo: convRepo,
		userRepo: userRepo,
		msgRepo:  msgRepo,
	}
}

// toParticipantResponses converts participants to flattened responses
func toParticipantResponses(participants []model.Participant) []*model.ParticipantResponse {
	result := make([]*model.ParticipantResponse, 0, len(participants))
	for _, p := range participants {
		if resp := p.ToResponse(); resp != nil {
			result = append(result, resp)
		}
	}
	return result
}

// CreateDirectRequest represents a request to create a direct conversation
type CreateDirectRequest struct {
	ParticipantID string `json:"participant_id" validate:"required"`
}

// CreateGroupRequest represents a request to create a group conversation
type CreateGroupRequest struct {
	Name           string   `json:"name" validate:"required"`
	ParticipantIDs []string `json:"participant_ids" validate:"required,min=1"`
}

// Create creates a new conversation (direct or group)
// POST /api/v1/conversations
func (h *ConversationHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// Try to parse as direct conversation first
	var directReq CreateDirectRequest
	if err := c.BodyParser(&directReq); err == nil && directReq.ParticipantID != "" {
		return h.createDirect(c, userID, directReq.ParticipantID)
	}

	// Try to parse as group conversation
	var groupReq CreateGroupRequest
	if err := c.BodyParser(&groupReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if groupReq.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Group name is required")
	}

	return h.createGroup(c, userID, &groupReq)
}

func (h *ConversationHandler) createDirect(c *fiber.Ctx, userID, participantID string) error {
	// Check if participant exists
	_, err := h.userRepo.GetByID(c.Context(), participantID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Participant not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify participant")
	}

	// Check if direct conversation already exists
	existing, err := h.convRepo.FindDirectConversation(c.Context(), userID, participantID)
	if err == nil && existing != nil {
		// Return existing conversation
		participants, _ := h.convRepo.GetParticipants(c.Context(), existing.ID)
		return c.JSON(fiber.Map{
			"conversation": existing,
			"participants": toParticipantResponses(participants),
			"is_new":       false,
		})
	}

	// Create new direct conversation
	conv := &model.Conversation{
		Type: model.ConversationTypeDirect,
	}

	participantIDs := []string{userID, participantID}
	if err := h.convRepo.Create(c.Context(), conv, participantIDs); err != nil {
		logger.Errorf("Failed to create conversation: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create conversation")
	}

	participants, _ := h.convRepo.GetParticipants(c.Context(), conv.ID)

	logger.Infof("Direct conversation created: %s", conv.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"conversation": conv,
		"participants": toParticipantResponses(participants),
		"is_new":       true,
	})
}

func (h *ConversationHandler) createGroup(c *fiber.Ctx, userID string, req *CreateGroupRequest) error {
	// Add creator to participants if not included
	participantIDs := req.ParticipantIDs
	hasCreator := false
	for _, id := range participantIDs {
		if id == userID {
			hasCreator = true
			break
		}
	}
	if !hasCreator {
		participantIDs = append([]string{userID}, participantIDs...)
	}

	// Verify all participants exist
	for _, id := range participantIDs {
		_, err := h.userRepo.GetByID(c.Context(), id)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return fiber.NewError(fiber.StatusNotFound, "Participant not found: "+id)
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify participants")
		}
	}

	// Create group conversation
	conv := &model.Conversation{
		Type:      model.ConversationTypeGroup,
		Name:      &req.Name,
		CreatedBy: &userID,
	}

	if err := h.convRepo.Create(c.Context(), conv, participantIDs); err != nil {
		logger.Errorf("Failed to create group: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create group")
	}

	participants, _ := h.convRepo.GetParticipants(c.Context(), conv.ID)

	logger.Infof("Group conversation created: %s (%s)", conv.ID, req.Name)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"conversation": conv,
		"participants": toParticipantResponses(participants),
	})
}

// List returns all conversations for the authenticated user
// GET /api/v1/conversations
func (h *ConversationHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	conversations, err := h.convRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		logger.Errorf("Failed to list conversations: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list conversations")
	}

	// Enrich with participants
	result := make([]fiber.Map, 0, len(conversations))
	for _, conv := range conversations {
		participants, _ := h.convRepo.GetParticipants(c.Context(), conv.ID)
		result = append(result, fiber.Map{
			"conversation": conv,
			"participants": toParticipantResponses(participants),
		})
	}

	return c.JSON(fiber.Map{
		"conversations": result,
		"count":         len(result),
	})
}

// Get returns a specific conversation with messages
// GET /api/v1/conversations/:id
func (h *ConversationHandler) Get(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	convID := c.Params("id")

	// Get conversation
	conv, err := h.convRepo.GetByID(c.Context(), convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Conversation not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get conversation")
	}

	// Verify user is participant
	participants, err := h.convRepo.GetParticipants(c.Context(), convID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get participants")
	}

	isParticipant := false
	for _, p := range participants {
		if p.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return fiber.NewError(fiber.StatusForbidden, "You are not a participant of this conversation")
	}

	return c.JSON(fiber.Map{
		"conversation": conv,
		"participants": toParticipantResponses(participants),
	})
}

// GetMessages returns messages for a conversation with pagination
// GET /api/v1/conversations/:id/messages
func (h *ConversationHandler) GetMessages(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	convID := c.Params("id")

	// Verify user is participant
	participants, err := h.convRepo.GetParticipants(c.Context(), convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Conversation not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify access")
	}

	isParticipant := false
	for _, p := range participants {
		if p.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return fiber.NewError(fiber.StatusForbidden, "You are not a participant of this conversation")
	}

	// Parse cursor (timestamp)
	var cursor *time.Time
	if cursorStr := c.Query("cursor"); cursorStr != "" {
		t, err := time.Parse(time.RFC3339Nano, cursorStr)
		if err == nil {
			cursor = &t
		}
	}

	// Parse limit
	limit := c.QueryInt("limit", 50)
	if limit > 100 {
		limit = 100
	}

	// Get messages
	page, err := h.msgRepo.GetByConversation(c.Context(), convID, cursor, limit)
	if err != nil {
		logger.Errorf("Failed to get messages: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get messages")
	}

	return c.JSON(page)
}
