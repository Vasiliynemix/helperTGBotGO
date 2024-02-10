package handlers

import (
	"bot/internal/bot/handlers/errorsMsg"
	"bot/internal/bot/handlers/register"
	"bot/internal/bot/handlers/start"
	userH "bot/internal/bot/handlers/user"
	"bot/internal/bot/keyboards"
	"bot/internal/bot/lexicon"
	"bot/internal/bot/middlewares"
	"bot/internal/storage"
	"bot/internal/storage/postgres/models"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type Handler struct {
	log             *logging.Logger
	updates         tgbotapi.UpdatesChannel
	lexicon         *lexicon.Lexicon
	mv              *middlewares.Middlewares
	bot             *tgbotapi.BotAPI
	storage         *storage.Storage
	errorsMsgSend   *errorsMsg.ErrorMsg
	startHandler    *start.HandlerStart
	registerHandler *register.HandlerRegister
	userHandler     *userH.HandlerUser
}

func NewHandlers(
	log *logging.Logger,
	bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel,
	mv *middlewares.Middlewares,
	storage *storage.Storage,
) *Handler {
	errorsMsgSend := errorsMsg.NewErrorsMsg(log, bot)
	l := lexicon.NewLexicon()
	kb := keyboards.NewKeyboards(l)
	startHandler := start.NewHandlerStart(
		log, bot, kb, l, errorsMsgSend,
		storage.Redis, storage.Postgres.User,
	)
	registerHandler := register.NewHandlerRegister(
		log, bot, kb, l, errorsMsgSend,
		storage.Redis, storage.Postgres.User,
	)

	userHandler := userH.NewHandlerUser(
		log, bot, kb, l, errorsMsgSend,
		storage.Redis,
	)

	return &Handler{
		log:             log,
		updates:         updates,
		lexicon:         l,
		mv:              mv,
		bot:             bot,
		storage:         storage,
		errorsMsgSend:   errorsMsgSend,
		startHandler:    startHandler,
		registerHandler: registerHandler,
		userHandler:     userHandler,
	}
}

func (h *Handler) CheckUpdates() {
	for update := range h.updates {
		timeNow := time.Now()
		state := h.mv.GetStateMv(update)
		if state == nil {
			go h.errorsMsgSend.MsgErrorInternal(update.Message.Chat.ID)
			continue
		}
		user := h.mv.GetUserMv(update)
		if user == nil {
			go h.errorsMsgSend.MsgErrorInternal(update.Message.Chat.ID)
		}

		switch {
		case update.Message != nil:
			ok := h.checkMessageUpdates(update.Message, user, state)
			if !ok {
				continue
			}

		case update.CallbackQuery != nil:
			ok := h.checkCallbackUpdates(update.CallbackQuery, user, state)
			if !ok {
				continue
			}

		default:
			continue
		}
		timeEnd := time.Since(timeNow)
		executionTimeMilliseconds := float64(timeEnd.Nanoseconds()) / 1e6
		h.mv.UpdateInfoMv(update, executionTimeMilliseconds)
	}
}

func (h *Handler) checkMessageUpdates(msg *tgbotapi.Message, user *models.User, state map[string]interface{}) bool {
	switch {
	case h.startHandler.CheckStart(msg):
		if user == nil || !user.IsRegistered {
			go h.startHandler.StartRegister(msg, user)
		} else {
			go h.startHandler.Start(msg)
		}
		return true

	case h.checkRegisterUpdates(msg, user, state):
		return true

	case h.userHandler.CheckProfileReply(msg):
		go h.userHandler.ProfileReply(msg)
		return true

	case h.userHandler.CheckBalanceReply(msg):
		go h.userHandler.BalanceReply(msg)
		return true

	case h.checkBackUpdates(msg, user, state):
		return true

	default:
		return false
	}
}

func (h *Handler) checkRegisterUpdates(msg *tgbotapi.Message, user *models.User, state map[string]interface{}) bool {
	switch {
	case h.registerHandler.CheckRegisterName(
		msg,
		(state)[h.lexicon.State.RegisterState.ID],
		h.lexicon.State.RegisterState.NameKey,
	):
		go h.registerHandler.RegisterName(msg)
		return true

	case h.registerHandler.CheckRegisterPhone(
		msg,
		(state)[h.lexicon.State.RegisterState.ID],
		h.lexicon.State.RegisterState.PhoneKey,
	):
		go h.registerHandler.RegisterPhone(msg, user)
		return true

	default:
		return false
	}
}

func (h *Handler) checkBackUpdates(msg *tgbotapi.Message, user *models.User, state map[string]interface{}) bool {
	switch {
	case h.userHandler.CheckBackToMenuReply(
		msg,
		(state)[h.lexicon.State.MenuLevelState.ID],
		h.lexicon.State.MenuLevelState.BackMenuKey,
	):
		go h.userHandler.BackToMenuReply(msg)
		return true

	case h.userHandler.CheckBackToProfileReply(
		msg,
		(state)[h.lexicon.State.MenuLevelState.ID],
		h.lexicon.State.MenuLevelState.BackMenuKey,
	):
		go h.userHandler.BackToProfileReply(msg)
		return true

	default:
		return false
	}
}

func (h *Handler) checkCallbackUpdates(callback *tgbotapi.CallbackQuery, user *models.User, state map[string]interface{}) bool {
	switch {

	default:
		return false
	}
}
