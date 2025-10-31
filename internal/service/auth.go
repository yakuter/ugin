package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yakuter/ugin/internal/domain"
	"github.com/yakuter/ugin/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
)

type authService struct {
	userRepo              repository.UserRepository
	jwtSecret             string
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
	logger                Logger
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, cfg *AuthConfig, logger Logger) AuthService {
	return &authService{
		userRepo:             userRepo,
		jwtSecret:            cfg.JWTSecret,
		accessTokenDuration:  cfg.AccessTokenDuration,
		refreshTokenDuration: cfg.RefreshTokenDuration,
		logger:               logger,
	}
}

func (s *authService) SignIn(ctx context.Context, creds *domain.Credentials) (*domain.TokenDetails, error) {
	if creds == nil || creds.Email == "" || creds.MasterPassword == "" {
		return nil, repository.ErrInvalidInput
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, creds.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.logger.Info("user not found during sign in", "email", creds.Email)
			return nil, ErrInvalidCredentials
		}
		s.logger.Error("failed to get user during sign in", "email", creds.Email, "error", err)
		return nil, fmt.Errorf("sign in: %w", err)
	}

	// TODO: Use proper password hashing (bcrypt)
	// For now, simple comparison (NOT PRODUCTION READY)
	if user.MasterPassword != creds.MasterPassword {
		s.logger.Warn("invalid password attempt", "email", creds.Email)
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokenDetails, err := s.createTokens(user.Email, user.ID)
	if err != nil {
		s.logger.Error("failed to create tokens", "email", creds.Email, "error", err)
		return nil, fmt.Errorf("create tokens: %w", err)
	}

	s.logger.Info("user signed in successfully", "email", creds.Email)
	return tokenDetails, nil
}

func (s *authService) SignUp(ctx context.Context, creds *domain.Credentials) error {
	if creds == nil || creds.Email == "" || creds.MasterPassword == "" {
		return repository.ErrInvalidInput
	}

	// TODO: Add password validation (min length, complexity, etc.)
	// TODO: Hash password with bcrypt

	user := &domain.User{
		Email:          creds.Email,
		MasterPassword: creds.MasterPassword, // TODO: Hash this!
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			s.logger.Info("user already exists", "email", creds.Email)
			return err
		}
		s.logger.Error("failed to create user", "email", creds.Email, "error", err)
		return fmt.Errorf("sign up: %w", err)
	}

	s.logger.Info("user signed up successfully", "email", creds.Email)
	return nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenDetails, error) {
	if refreshToken == "" {
		return nil, repository.ErrInvalidInput
	}

	// Validate refresh token
	token, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Create new tokens
	tokenDetails, err := s.createTokens(email, uint(userID))
	if err != nil {
		s.logger.Error("failed to refresh tokens", "email", email, "error", err)
		return nil, fmt.Errorf("refresh token: %w", err)
	}

	s.logger.Info("token refreshed successfully", "email", email)
	return tokenDetails, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*domain.TokenClaims, error) {
	if tokenString == "" {
		return nil, repository.ErrInvalidInput
	}

	token, err := s.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	email, _ := claims["email"].(string)
	userUUID, _ := claims["user_uuid"].(string)
	uuid, _ := claims["uuid"].(string)

	return &domain.TokenClaims{
		Email:    email,
		UserUUID: userUUID,
		UUID:     uuid,
	}, nil
}

func (s *authService) createTokens(email string, userID uint) (*domain.TokenDetails, error) {
	td := &domain.TokenDetails{}

	now := time.Now()
	td.ATExpiresAt = now.Add(s.accessTokenDuration)
	td.RTExpiresAt = now.Add(s.refreshTokenDuration)

	// Create access token
	atClaims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"exp":     td.ATExpiresAt.Unix(),
		"iat":     now.Unix(),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	td.AccessToken = accessToken

	// Create refresh token
	rtClaims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"exp":     td.RTExpiresAt.Unix(),
		"iat":     now.Unix(),
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	td.RefreshToken = refreshToken

	// Generate transmission key (for additional security)
	td.TransmissionKey = generateSecureKey(16)

	return td, nil
}

func (s *authService) parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	return token, nil
}

// generateSecureKey generates a random secure key
func generateSecureKey(length int) string {
	// TODO: Use crypto/rand for production
	// This is a placeholder
	return fmt.Sprintf("%016x", time.Now().UnixNano())
}

