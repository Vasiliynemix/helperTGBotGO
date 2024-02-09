package repositories

import (
	"bot/internal/storage/postgres/models"
	"bot/pkg/logging"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

const (
	ErrDuplicateKeyStr = "duplicate key value violates unique constraint"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserIsExists = errors.New("user already exists")
)

type UserRepository struct {
	log *logging.Logger
	db  *gorm.DB
}

func NewUserRepository(log *logging.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{
		log: log,
		db:  db,
	}
}

func (r *UserRepository) SetUser(user *models.User) (*models.User, error) {
	result := r.db.Save(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), ErrDuplicateKeyStr) {
			r.log.Info("User already exists", zap.Error(result.Error))

			return nil, ErrUserIsExists
		}
		r.log.Error("Failed to add user", zap.Error(result.Error))

		return nil, result.Error
	}

	r.log.Info("User added", zap.Int64("telegram_id", user.TelegramID))

	return user, nil
}

func (r *UserRepository) GetUser(telegramID int64) (*models.User, error) {
	var userFind models.User
	mapFind := map[string]interface{}{
		"telegram_id": telegramID,
	}

	result := r.db.Find(&userFind, mapFind).First(&userFind)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.log.Warn("User not found", zap.Int64("telegram_id", telegramID))

			return nil, ErrUserNotFound
		}

		r.log.Error("Failed to get user", zap.Int64("telegram_id", telegramID), zap.Error(result.Error))

		return nil, result.Error
	}

	return &userFind, nil
}

func (r *UserRepository) UpdateUser(telegramID int64, user *models.User) error {
	result := r.db.Model(&models.User{}).Where("telegram_id = ?", telegramID).Updates(user)
	if result.Error != nil {
		r.log.Error("Failed to update user", zap.Int64("telegram_id", telegramID), zap.Error(result.Error))
		return result.Error
	}

	return nil
}
