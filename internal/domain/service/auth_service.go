package service

import (
	"app/internal/domain"
	"app/internal/domain/model"
	"app/internal/infra/db"
	"app/internal/infra/service"
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *db.UserRepository
	tokenService service.TokenService
}

func NewAuthService(userRepo *db.UserRepository, tokenService service.TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) SignUp(c context.Context, username string, password string) (*model.User, error) {
	existingUser, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(passwordHash),
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *AuthService) SignIn(c context.Context, username string, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := s.tokenService.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
