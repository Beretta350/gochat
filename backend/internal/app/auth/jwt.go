package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/logger"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// TokenType represents the type of JWT token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   string    `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Type     TokenType `json:"type"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
}

// JWTService handles JWT operations
type JWTService struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTService creates a new JWT service (Fx provider)
func NewJWTService(cfg *config.Config) *JWTService {
	logger.Info("JWT service initialized")
	return &JWTService{
		secret:        []byte(cfg.JWT.Secret),
		accessExpiry:  cfg.JWT.AccessExpiry,
		refreshExpiry: cfg.JWT.RefreshExpiry,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(userID, email, username string) (*TokenPair, error) {
	accessToken, err := s.generateToken(userID, email, username, AccessToken, s.accessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(userID, email, username, RefreshToken, s.refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
	}, nil
}

// GenerateAccessToken generates only an access token (for refresh)
func (s *JWTService) GenerateAccessToken(userID, email, username string) (string, error) {
	return s.generateToken(userID, email, username, AccessToken, s.accessExpiry)
}

func (s *JWTService) generateToken(userID, email, username string, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken validates a token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != AccessToken {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != RefreshToken {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
