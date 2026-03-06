package user

import (
	"Mini-Project-Game-Vault-API/service/user"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user user.User) error {
	return r.DB.WithContext(ctx).Create(&user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (user.User, error) {
	var user user.User

	err := r.DB.WithContext(ctx).First(&user, "id= ?", id).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	var user user.User

	err := r.DB.WithContext(ctx).First(&user, "email = ?", email).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	var user user.User

	err := r.DB.WithContext(ctx).First(&user, "username = ?", username).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) UpdatEmailVerification(ctx context.Context, user user.User) error {
	return r.DB.WithContext(ctx).Updates(&user).Error
}

func (r *UserRepository) UpdateDeposit(ctx context.Context, userID string, newBalance int) error {
	var user user.User

	return r.DB.WithContext(ctx).Model(&user).Where("id = ?", userID).Update("deposit_amount", newBalance).Error
}
