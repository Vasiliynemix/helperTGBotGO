package start

import (
	"bot/internal/bot/handlers/errorsMsg"
	"bot/internal/bot/lexicon"
	"bot/internal/storage/postgres/models"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type HandlerStart struct {
	log           *logging.Logger
	bot           *tgbotapi.BotAPI
	lexicon       *lexicon.Lexicon
	errMsg        *errorsMsg.ErrorMsg
	stateProvider StateProvider
	userProvider  UserProvider
}

type UserProvider interface {
	UpdateUser(telegramID int64, user *models.User) error
}

type StateProvider interface {
	ClearAllStates(telegramID int64) error
	ClearState(telegramID int64, stateName string) error
	GetStateData(telegramID int64, stateName string) (map[string]interface{}, error)
	UpdateState(telegramID int64, stateName string, fieldName string, fieldValue interface{}) error
}

func NewHandlerStart(
	log *logging.Logger,
	bot *tgbotapi.BotAPI,
	errMsg *errorsMsg.ErrorMsg,
	stateProvider StateProvider,
	userProvider UserProvider,
) *HandlerStart {
	log.Info("Start handler created")

	return &HandlerStart{
		log:           log,
		bot:           bot,
		lexicon:       lexicon.NewLexicon(),
		errMsg:        errMsg,
		stateProvider: stateProvider,
		userProvider:  userProvider,
	}
}

func (r *HandlerStart) CheckStart(msg *tgbotapi.Message) bool {
	return msg != nil && msg.Command() == r.lexicon.Cmd.Start
}

func (r *HandlerStart) StartRegister(msg *tgbotapi.Message, user *models.User) {
	_ = r.stateProvider.ClearAllStates(msg.Chat.ID)

	msgSend := tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.OnRegisterStartCommand)

	err := r.userProvider.UpdateUser(msg.Chat.ID, user)
	if err != nil {
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	_, err = r.bot.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
		return
	}

	err = r.stateProvider.UpdateState(msg.Chat.ID, r.lexicon.State.RegisterState.ID, r.lexicon.State.RegisterState.NameKey, true)
	if err != nil {
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
	}
}

func (r *HandlerStart) Start(msg *tgbotapi.Message) {
	_ = r.stateProvider.ClearAllStates(msg.Chat.ID)

	msgSend := tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.OnStartCommand)

	_, err := r.bot.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
		return
	}
}
