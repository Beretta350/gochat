package auth

import (
	"context"
	"errors"

	"github.com/Beretta350/gochat/internal/app/model"
	"github.com/Beretta350/gochat/internal/app/repository"
	"github.com/Beretta350/gochat/pkg/logger"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// Service handles authentication operations
type Service struct {
	userRepo   repository.UserRepository
	jwtService *JWTService
}

// NewService creates a new auth service (Fx provider)
func NewService(userRepo repository.UserRepository, jwtService *JWTService) *Service {
	logger.Info("Auth service initialized")
	return &Service{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=8"`
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

	return &AuthResponse{
		User:   user.ToResponse(),
		Tokens: tokens,
	}, nil
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
