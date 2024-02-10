package register

import (
	"bot/internal/bot/handlers/errorsMsg"
	"bot/internal/bot/keyboards"
	"bot/internal/bot/lexicon"
	"bot/internal/storage/postgres/models"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type HandlerRegister struct {
	log           *logging.Logger
	bot           *tgbotapi.BotAPI
	kb            *keyboards.Keyboards
	lexicon       *lexicon.Lexicon
	errMsg        *errorsMsg.ErrorMsg
	stateProvider StateProvider
	userUpdater   UserUpdater
}

type StateProvider interface {
	ClearAllStates(telegramID int64) error
	ClearState(telegramID int64, stateName string) error
	GetStateData(telegramID int64, stateName string) (map[string]interface{}, error)
	UpdateState(telegramID int64, stateName string, fieldName string, fieldValue interface{}) error
}

type UserUpdater interface {
	UpdateUser(telegramID int64, user *models.User) error
}

func NewHandlerRegister(
	log *logging.Logger,
	bot *tgbotapi.BotAPI,
	kb *keyboards.Keyboards,
	lexicon *lexicon.Lexicon,
	errMsg *errorsMsg.ErrorMsg,
	stateProvider StateProvider,
	userUpdater UserUpdater,
) *HandlerRegister {
	log.Info("register handler created")
	return &HandlerRegister{
		log:           log,
		bot:           bot,
		kb:            kb,
		lexicon:       lexicon,
		errMsg:        errMsg,
		stateProvider: stateProvider,
		userUpdater:   userUpdater,
	}
}

func (r *HandlerRegister) CheckRegisterName(msg *tgbotapi.Message, state interface{}, key string) bool {
	if msg.Text == "" || state == nil {
		return false
	}
	value, ok := state.(map[string]interface{})[key].(bool)
	if !ok {
		r.log.Warn("state isRegisteredName not boolean", zap.Any("state", state))
		return false
	}

	ok = msg != nil && value
	if ok {
		r.log.Debug("CheckRegisterName", zap.Any("stateRegister", state))
	}
	return ok
}

func (r *HandlerRegister) RegisterName(msg *tgbotapi.Message) {
	stateID := r.lexicon.State.RegisterState.ID
	err := r.stateProvider.UpdateState(msg.Chat.ID, stateID, r.lexicon.State.RegisterState.NameKey, false)
	err = r.stateProvider.UpdateState(msg.Chat.ID, stateID, r.lexicon.State.RegisterState.PhoneKey, true)
	err = r.stateProvider.UpdateState(msg.Chat.ID, stateID, "name", msg.Text)
	if err != nil {
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	msgSend := tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.CreateMsgOnRegisterPhone(msg.Text))

	_, err = r.bot.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}
}

func (r *HandlerRegister) CheckRegisterPhone(msg *tgbotapi.Message, state interface{}, key string) bool {
	if msg.Text == "" || state == nil {
		return false
	}
	value, ok := state.(map[string]interface{})[key].(bool)
	if !ok {
		r.log.Warn("state isRegisteredPhone not bool", zap.Any("state", state))
		return false
	}

	ok = msg != nil && value
	if ok {
		r.log.Debug("CheckRegisterPhone", zap.Any("stateRegisterPhone", state))
	}
	return ok
}

func (r *HandlerRegister) RegisterPhone(msg *tgbotapi.Message, user *models.User) {
	data, err := r.stateProvider.GetStateData(msg.Chat.ID, r.lexicon.State.RegisterState.ID)
	if err != nil {
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	name := data["name"].(string)
	phone := msg.Text

	user.Name = name
	user.Phone = phone
	user.IsRegistered = true
	err = r.userUpdater.UpdateUser(msg.Chat.ID, user)
	if err != nil {
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	err = r.stateProvider.ClearState(msg.Chat.ID, r.lexicon.State.RegisterState.ID)
	if err != nil {
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	msgSend := tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.OnEndRegister)

	_, err = r.bot.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}

	msgSend = tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.OnStartCommand)
	msgSend.ReplyMarkup = r.kb.Reply.StartMenuReplyMP(false)

	_, err = r.bot.Send(msgSend)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
	}

	r.log.Info("RegisterEnd", zap.Any("name", name), zap.Any("phone", phone))
}
