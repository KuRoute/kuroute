package services

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/KuRoute/kuroute/backend/package/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) RegisterUser(req *domain.RegisterUserRequest) (*domain.UserResponse, error) {
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		HubID:        req.HubID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Role:         req.Role,
		PasswordHash: string(hash),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := domain.NewUserResponse(user)
	return &response, nil
}

func (s *AuthService) Login(email, password string) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Role, user.HubID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Role, user.HubID)
	if err != nil {
		return nil, err
	}

	expiresIn := int(jwt.GetAccessExpiryDuration().Seconds())

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         domain.NewUserResponse(user),
	}, nil
}

func (s *AuthService) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
	if err != nil {
		return nil, errors.New("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = string(hash)
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Role, user.HubID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Role, user.HubID)
	if err != nil {
		return nil, err
	}

	expiresIn := int(jwt.GetAccessExpiryDuration().Seconds())

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         domain.NewUserResponse(user),
	}, nil
}

func (s *AuthService) RefreshToken(oldRefreshToken string) (*domain.TokenRefreshResponse, error) {
	// Stateless validation: verify the refresh token cryptographically
	claims, err := jwt.ValidateTokenAuth(oldRefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Generate a new token pair using claims from the old refresh token
	newAccessToken, err := jwt.GenerateAccessToken(claims.UserID, claims.Role, claims.HubID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(claims.UserID, claims.Role, claims.HubID)
	if err != nil {
		return nil, err
	}

	expiresIn := int(jwt.GetAccessExpiryDuration().Seconds())

	return &domain.TokenRefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	// Stateless: no DB operation needed.
	// The mobile app will delete the token from its SecureStore.
	// Optionally validate the token to decide if it's well-formed, but return success regardless.
	if refreshToken == "" {
		return errors.New("refresh token is required")
	}
	return nil
}
