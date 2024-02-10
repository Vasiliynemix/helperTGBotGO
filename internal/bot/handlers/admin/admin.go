package admin

import (
	"bot/internal/bot/handlers/errorsMsg"
	"bot/internal/bot/keyboards"
	"bot/internal/bot/lexicon"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type HandlerAdmin struct {
	log           *logging.Logger
	bot           *tgbotapi.BotAPI
	kb            *keyboards.Keyboards
	lexicon       *lexicon.Lexicon
	errMsg        *errorsMsg.ErrorMsg
	stateProvider StateProvider
}

type StateProvider interface {
	UpdateStateMap(telegramID int64, stateName string, stateData map[string]interface{}) error
	ClearState(telegramID int64, stateName string) error
}

func NewHandlerAdmin(
	log *logging.Logger,
	bot *tgbotapi.BotAPI,
	kb *keyboards.Keyboards,
	lexicon *lexicon.Lexicon,
	errMsg *errorsMsg.ErrorMsg,
	stateProvider StateProvider,
) *HandlerAdmin {
	log.Info("admin handler created")
	return &HandlerAdmin{
		log:           log,
		bot:           bot,
		kb:            kb,
		lexicon:       lexicon,
		errMsg:        errMsg,
		stateProvider: stateProvider,
	}
}

func (h *HandlerAdmin) CheckAdminPanelReply(msg *tgbotapi.Message) bool {
	return msg != nil && msg.Text == h.lexicon.KB.Reply.Admin
}

func (h *HandlerAdmin) AdminPanelReply(msg *tgbotapi.Message) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, h.lexicon.Msg.OnAdminPanel)
	msgSend.ReplyMarkup = h.kb.Reply.AdminPanelReplyMP(false)

	err := h.stateProvider.ClearState(msg.Chat.ID, h.lexicon.State.MenuLevelState.ID)
	if err != nil {
		h.log.Error("failed to clear state", zap.Error(err))
		h.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	_, err = h.bot.Send(msgSend)
	if err != nil {
		h.log.Error("failed to send message", zap.Error(err))
	}
}

func (h *HandlerAdmin) CheckUserPanelReply(msg *tgbotapi.Message) bool {
	return msg != nil && msg.Text == h.lexicon.KB.Reply.User
}

func (h *HandlerAdmin) UserPanelReply(msg *tgbotapi.Message) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, h.lexicon.Msg.OnUserPanel)
	msgSend.ReplyMarkup = h.kb.Reply.StartMenuReplyMP(false, true)

	err := h.stateProvider.ClearState(msg.Chat.ID, h.lexicon.State.MenuLevelState.ID)
	if err != nil {
		h.log.Error("failed to clear state", zap.Error(err))
		h.errMsg.MsgErrorInternal(msg.Chat.ID)
		return
	}

	_, err = h.bot.Send(msgSend)
	if err != nil {
		h.log.Error("failed to send message", zap.Error(err))
	}
}
