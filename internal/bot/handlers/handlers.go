package handlers

import (
	"bot/internal/bot/handlers/errorsMsg"
	"bot/internal/bot/handlers/register"
	"bot/internal/bot/handlers/start"
	"bot/internal/bot/keybords"
	"bot/internal/bot/lexicon"
	"bot/internal/bot/middlewares"
	"bot/internal/storage"
	"bot/internal/storage/postgres/models"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type Handler struct {
	updates         tgbotapi.UpdatesChannel
	lexicon         *lexicon.Lexicon
	mv              *middlewares.Middlewares
	bot             *tgbotapi.BotAPI
	storage         *storage.Storage
	errorsMsgSend   *errorsMsg.ErrorMsg
	startHandler    *start.HandlerStart
	registerHandler *register.HandlerRegister
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
	kb := keybords.NewKeyboards(l)
	startHandler := start.NewHandlerStart(
		log, bot, kb, l, errorsMsgSend,
		storage.Redis, storage.Postgres.User,
	)
	registerHandler := register.NewHandlerRegister(
		log, bot, kb, l, errorsMsgSend,
		storage.Redis, storage.Postgres.User,
	)

	return &Handler{
		updates:         updates,
		lexicon:         l,
		mv:              mv,
		bot:             bot,
		storage:         storage,
		errorsMsgSend:   errorsMsgSend,
		startHandler:    startHandler,
		registerHandler: registerHandler,
	}
}

func (h *Handler) CheckUpdates() {
	for update := range h.updates {
		timeNow := time.Now()
		state := h.mv.GetStateMv(update)
		if state == nil {
			h.errorsMsgSend.MsgErrorInternal(update.Message.Chat.ID)
			continue
		}
		user := h.mv.GetUserMv(update)
		if user == nil {
			h.errorsMsgSend.MsgErrorInternal(update.Message.Chat.ID)
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
			h.startHandler.StartRegister(msg, user)
		} else {
			h.startHandler.Start(msg)
		}
		return true

	case h.registerHandler.CheckRegisterName(
		msg,
		(state)[h.lexicon.State.RegisterState.ID],
		h.lexicon.State.RegisterState.NameKey,
	):
		h.registerHandler.RegisterName(msg)
		return true

	case h.registerHandler.CheckRegisterPhone(
		msg,
		(state)[h.lexicon.State.RegisterState.ID],
		h.lexicon.State.RegisterState.PhoneKey,
	):
		h.registerHandler.RegisterPhone(msg, user)
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
