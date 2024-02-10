package user

import (
	"bot/internal/bot/handlers/errorsMsg"
	"bot/internal/bot/keyboards"
	"bot/internal/bot/lexicon"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type HandlerUser struct {
	log           *logging.Logger
	bot           *tgbotapi.BotAPI
	kb            *keyboards.Keyboards
	lexicon       *lexicon.Lexicon
	errMsg        *errorsMsg.ErrorMsg
	stateProvider StateProvider
}

type StateProvider interface {
	UpdateState(telegramID int64, stateName string, fieldName string, fieldValue interface{}) error
	UpdateStateMap(telegramID int64, stateName string, stateData map[string]interface{}) error
}

func NewHandlerUser(
	log *logging.Logger,
	bot *tgbotapi.BotAPI,
	kb *keyboards.Keyboards,
	lexicon *lexicon.Lexicon,
	errMsg *errorsMsg.ErrorMsg,
	stateProvider StateProvider,
) *HandlerUser {
	log.Info("user handler created")
	return &HandlerUser{
		log:           log,
		bot:           bot,
		kb:            kb,
		lexicon:       lexicon,
		errMsg:        errMsg,
		stateProvider: stateProvider,
	}
}

func (r *HandlerUser) CheckProfileReply(msg *tgbotapi.Message) bool {
	return msg != nil && msg.Text == r.lexicon.KB.Reply.Profile
}

func (r *HandlerUser) ProfileReply(msg *tgbotapi.Message) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.OnProfile)
	msgSend.ReplyMarkup = r.kb.Reply.ProfileReplyMP(false)

	err := r.stateProvider.UpdateStateMap(
		msg.Chat.ID,
		r.lexicon.State.MenuLevelState.ID,
		map[string]interface{}{
			r.lexicon.State.MenuLevelState.CurrentMenuKey: r.lexicon.State.MenuLevelState.ProfileValue,
			r.lexicon.State.MenuLevelState.BackMenuKey:    r.lexicon.State.MenuLevelState.MenuValue,
		},
	)
	if err != nil {
		r.log.Error("can't update state", zap.Error(err))
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	_, err = r.bot.Send(msgSend)
	if err != nil {
		r.log.Error("can't send message", zap.Error(err))
	}
}

func (r *HandlerUser) CheckBalanceReply(msg *tgbotapi.Message) bool {
	return msg != nil && msg.Text == r.lexicon.KB.Reply.Balance
}

func (r *HandlerUser) BalanceReply(msg *tgbotapi.Message) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.OnBalance)
	msgSend.ReplyMarkup = r.kb.Reply.BalanceReplyMP(false)

	err := r.stateProvider.UpdateStateMap(
		msg.Chat.ID,
		r.lexicon.State.MenuLevelState.ID,
		map[string]interface{}{
			r.lexicon.State.MenuLevelState.CurrentMenuKey: r.lexicon.State.MenuLevelState.BalanceValue,
			r.lexicon.State.MenuLevelState.BackMenuKey:    r.lexicon.State.MenuLevelState.ProfileValue,
		},
	)
	if err != nil {
		r.log.Error("can't update state", zap.Error(err))
		r.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	_, err = r.bot.Send(msgSend)
	if err != nil {
		r.log.Error("can't send message", zap.Error(err))
	}
}
