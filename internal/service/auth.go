package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"example.com/acg-go-demo/internal/dto"
	"example.com/acg-go-demo/internal/repository"
	"example.com/acg-go-demo/internal/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, rawToken string) error
}

type authService struct {
	userRepo repository.UserRepository
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByAccount(ctx, req.Account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("user not found. account=%q, err=%v", req.Account, err)
			return nil, ErrInvalidCredentials
		} else {
			log.Printf("get user by account failed, account=%q, err=%v", req.Account, err)
			return nil, fmt.Errorf("get user by account: %w", err)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("password verification failed, account=%q", req.Account)
		return nil, ErrInvalidCredentials
	}

	rawToken, err := token.Generate()
	if err != nil {
		log.Printf("generate auth token failed, account=%q, err=%v", req.Account, err)
		return nil, fmt.Errorf("generate auth token: %w", err)
	}

	if err := s.userRepo.UpdateAuthTokenHash(ctx, user.ID, token.Hash(rawToken)); err != nil {
		log.Printf("save auth token hash failed, user_id=%d, err=%v", user.ID, err)
		return nil, fmt.Errorf("save auth token: %w", err)
	}

	return &dto.LoginResponse{
		Token: rawToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Account:  user.Account,
			Nickname: user.Nickname,
		},
	}, nil

}

func (s *authService) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}

	user, err := s.userRepo.GetByAuthTokenHash(ctx, token.Hash(rawToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get user by auth token: %w", err)
	}

	if err := s.userRepo.ClearAuthTokenHash(ctx, user.ID); err != nil {
		return fmt.Errorf("clear auth token: %w", err)
	}

	return nil
}
