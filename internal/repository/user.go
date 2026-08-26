package repository

import (
	"context"

	"example.com/acg-go-demo/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByAccount(ctx context.Context, account string) (*model.User, error)
	UpdateNickname(ctx context.Context, id uint64, nickname string) error
	Delete(ctx context.Context, id uint64) error
	UpdateAuthTokenHash(ctx context.Context, id uint64, hash string) error
	ClearAuthTokenHash(ctx context.Context, id uint64) error
	GetByAuthTokenHash(ctx context.Context, hash string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("account = ?", account).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateNickname(ctx context.Context, id uint64, nickname string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("nickname", nickname).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) UpdateAuthTokenHash(ctx context.Context, id uint64, hash string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("auth_token_hash", hash).Error
}

func (r *userRepository) ClearAuthTokenHash(ctx context.Context, id uint64) error {
	return r.UpdateAuthTokenHash(ctx, id, "")
}

func (r *userRepository) GetByAuthTokenHash(ctx context.Context, hash string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("auth_token_hash = ?", hash).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
