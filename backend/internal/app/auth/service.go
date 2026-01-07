package auth

import (
	"context"
	"errors"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/internal/app/repository"
	"github.com/Beretta350/gochat/pkg/logger"
)

// DefaultChatUserEmail is the email of the user that all new users should have a direct chat with
const DefaultChatUserEmail = "beretta.gabrielpp@gmail.com"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// Service handles authentication operations
type Service struct {
	userRepo   repository.UserRepository
	convRepo   repository.ConversationRepository
	jwtService *JWTService
}

// NewService creates a new auth service (Fx provider)
func NewService(userRepo repository.UserRepository, convRepo repository.ConversationRepository, jwtService *JWTService) *Service {
	logger.Info("Auth service initialized")
	return &Service{
		userRepo:   userRepo,
		convRepo:   convRepo,
		jwtService: jwtService,
	}
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=128,strongpassword"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User   *model.UserResponse `json:"user"`
	Tokens *TokenPair          `json:"tokens"`
}

// Register creates a new user and returns tokens
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Create user
	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		IsActive: true,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	// Generate tokens
	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	logger.Infof("User registered: %s", user.Email)

	// Auto-create direct chat with Gabriel if:
	// 1. The registering user is NOT Gabriel
	// 2. Gabriel's user exists in the system
	if req.Email != DefaultChatUserEmail {
		s.createDefaultDirectChat(ctx, user.ID)
	}

	return &AuthResponse{
		User:   user.ToResponse(),
		Tokens: tokens,
	}, nil
}

// createDefaultDirectChat creates a direct conversation between the new user and Gabriel
func (s *Service) createDefaultDirectChat(ctx context.Context, newUserID string) {
	// Check if Gabriel's user exists
	gabriel, err := s.userRepo.GetByEmail(ctx, DefaultChatUserEmail)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			logger.Info("Default chat user (Gabriel) not found, skipping auto-chat creation")
			return
		}
		logger.Errorf("Error checking default chat user: %v", err)
		return
	}

	// Check if a direct conversation already exists
	_, err = s.convRepo.FindDirectConversation(ctx, newUserID, gabriel.ID)
	if err == nil {
		// Conversation already exists
		logger.Infof("Direct conversation with Gabriel already exists for user %s", newUserID)
		return
	}

	// Create direct conversation
	conv := &model.Conversation{
		Type: model.ConversationTypeDirect,
	}
	participantIDs := []string{newUserID, gabriel.ID}

	if err := s.convRepo.Create(ctx, conv, participantIDs); err != nil {
		logger.Errorf("Failed to create default direct chat: %v", err)
		return
	}

	logger.Infof("Auto-created direct chat between new user %s and Gabriel", newUserID)
}

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	logger.Infof("User logged in: %s", user.Email)

	return &AuthResponse{
		User:   user.ToResponse(),
		Tokens: tokens,
	}, nil
}

// Refresh generates a new access token from a refresh token
func (s *Service) Refresh(ctx context.Context, req *RefreshRequest) (*TokenPair, error) {
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Verify user still exists and is active
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserNotFound
	}

	// Generate new token pair
	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	logger.Infof("Token refreshed for: %s", user.Email)

	return tokens, nil
}

// GetUserFromToken extracts user info from a valid access token
func (s *Service) GetUserFromToken(tokenString string) (*Claims, error) {
	return s.jwtService.ValidateAccessToken(tokenString)
}
