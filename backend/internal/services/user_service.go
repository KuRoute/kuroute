package services

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := domain.NewUserResponse(user)
	return &response, nil
}

func (s *UserService) CreateUser(req *domain.RegisterUserRequest, callerRole domain.UserRole) (*domain.UserResponse, error) {
	if callerRole != domain.UserRoleAdmin {
		return nil, errors.New("unauthorized: cannot create user with different role")
	}

	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}
	
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		HubID: req.HubID,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Role:  req.Role,
		PasswordHash: string(passwordHash),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := domain.NewUserResponse(user)
	return &response, nil
}

func (s *UserService) GetUserByEmail(email string, callerRole domain.UserRole) (*domain.UserResponse, error) {
	if callerRole != domain.UserRoleAdmin {
		return nil, errors.New("unauthorized: cannot access user with different role")
	}

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := domain.NewUserResponse(user)
	return &response, nil
}

func (s *UserService) UpdateUser(req *domain.UpdateUserRequest, callerRole domain.UserRole) error {
	if callerRole != domain.UserRoleAdmin {
		return errors.New("unauthorized: cannot update user with different role")
	}

	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	err := s.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUser(id uuid.UUID, callerRole domain.UserRole) error {
	if callerRole != domain.UserRoleAdmin {
		return errors.New("unauthorized: cannot update user with different role")
	}

	err := s.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}