package middlewares

import (
	"bot/internal/storage/postgres/models"
	"bot/pkg/logging"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

type Middlewares struct {
	log           *logging.Logger
	userProvider  UserProvider
	stateProvider StateProvider
}

type UserProvider interface {
	SetUser(user *models.User) (*models.User, error)
	GetUser(telegramID int64) (*models.User, error)
}

type StateProvider interface {
	GetStateAll(telegramID int64) (map[string]interface{}, error)
}

func InitMiddlewares(
	log *logging.Logger,
	userProvider UserProvider,
	stateProvider StateProvider,
) *Middlewares {
	return &Middlewares{
		log:           log,
		userProvider:  userProvider,
		stateProvider: stateProvider,
	}
}

func (m *Middlewares) UpdateInfoMv(msg tgbotapi.Update, timeout float64) {
	switch {
	case msg.Message != nil:
		m.log.Info(
			"Update info",
			zap.String("timeout", fmt.Sprintf("%f MS", timeout)),
			zap.Int("chat_id", int(msg.Message.Chat.ID)),
			zap.String("username", msg.Message.From.UserName),
			zap.String("text", msg.Message.Text),
		)
	case msg.CallbackQuery != nil:
		m.log.Info(
			"Update info",
			zap.String("timeout", fmt.Sprintf("%f MS", timeout)),
			zap.Int("chat_id", int(msg.CallbackQuery.Message.Chat.ID)),
			zap.String("username", msg.CallbackQuery.From.UserName),
			zap.String("data", msg.CallbackQuery.Data),
		)
	}
}

func (m *Middlewares) GetUserMv(msg tgbotapi.Update) *models.User {
	var telegramID int64
	var userName string

	switch {
	case msg.Message != nil:
		telegramID = msg.Message.Chat.ID
		userName = msg.Message.From.UserName

	case msg.CallbackQuery != nil:
		telegramID = msg.CallbackQuery.Message.Chat.ID
		userName = msg.CallbackQuery.From.UserName
	}

	userShow, err := m.userProvider.GetUser(telegramID)
	if err == nil {
		return userShow
	}

	userAdd := &models.User{
		TelegramID: telegramID,
		UserName:   userName,
		CreatedAt:  time.Now().UTC().Unix(),
	}

	userShow, err = m.userProvider.SetUser(userAdd)
	if err != nil {
		return nil
	}
	return userShow
}

func (m *Middlewares) GetStateMv(msg tgbotapi.Update) map[string]interface{} {
	var telegramID int64

	switch {
	case msg.Message != nil:
		telegramID = msg.Message.Chat.ID

	case msg.CallbackQuery != nil:
		telegramID = msg.CallbackQuery.Message.Chat.ID
	}

	state, err := m.stateProvider.GetStateAll(telegramID)
	if err != nil {
		return nil
	}

	return state
}
