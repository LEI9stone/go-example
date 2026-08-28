package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"example.com/acg-go-demo/internal/dto"
	"example.com/acg-go-demo/internal/model"
	"example.com/acg-go-demo/internal/repository"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Get(ctx context.Context, id uint64) (*dto.UserResponse, error)
	UpdateNickname(ctx context.Context, id uint64, nickname string) error
	Delete(ctx context.Context, id uint64) error
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrAccountTaken = errors.New("account already exists")
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	_, err := s.userRepo.GetByAccount(ctx, req.Account)

	switch {
	case err == nil:
		// 查到用户，说明账号已被占用
		return nil, ErrAccountTaken
	case errors.Is(err, gorm.ErrRecordNotFound):
		// 未查到用户，继续后续流程
	default:
		// 数据库连接错误，查询超时等其他错误
		log.Printf("get user by account failed, account=%q, err=%v", req.Account, err)
		return nil, fmt.Errorf("check account: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("hash password failed, err=%v", err)
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Account:      req.Account,
		PasswordHash: string(passwordHash),
		Nickname:     req.Nickname,
		State:        1,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrAccountTaken
		}
		log.Printf("create user failed, account=%q, err=%v", req.Account, err)
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Account:  user.Account,
		Nickname: user.Nickname,
	}, nil
}

func (s *userService) Delete(ctx context.Context, id uint64) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		log.Printf("delete user failed, id=%d, err=%v", id, err)
		return fmt.Errorf("delete user error: %w", err)
	}
	return nil
}

func (s *userService) Get(ctx context.Context, id uint64) (*dto.UserResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid user id: %d", id)
	}
	user, err := s.userRepo.GetByID(ctx, id)
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		log.Printf("user not found, id=%d, err=%v", id, err)
		return nil, ErrUserNotFound
	default:
		// 数据库连接错误，查询超时等其他错误
		log.Printf("get user by id failed, id=%d, err=%v", id, err)
		return nil, fmt.Errorf("get user by id : %w", err)
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Account:  user.Account,
		Nickname: user.Nickname,
	}, nil
}
