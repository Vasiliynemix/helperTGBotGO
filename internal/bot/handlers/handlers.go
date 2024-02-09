package handlers

import (
	"bot/internal/bot/handlers/errorsMsg"
	"bot/internal/bot/handlers/register"
	"bot/internal/bot/handlers/start"
	"bot/internal/bot/lexicon"
	"bot/internal/bot/middlewares"
	"bot/internal/storage"
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
	startHandler := start.NewHandlerStart(log, bot, errorsMsgSend, storage.Redis, storage.Postgres.User)
	registerHandler := register.NewHandlerRegister(log, bot, errorsMsgSend, storage.Redis, storage.Postgres.User)

	return &Handler{
		updates:         updates,
		lexicon:         lexicon.NewLexicon(),
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
		case h.startHandler.CheckStart(update.Message):
			if user == nil || !user.IsRegistered {
				h.startHandler.StartRegister(update.Message, user)
			} else {
				h.startHandler.Start(update.Message)
			}

		case h.registerHandler.CheckRegisterName(
			update.Message,
			(state)[h.lexicon.State.RegisterState.ID],
			h.lexicon.State.RegisterState.NameKey,
		):
			h.registerHandler.RegisterName(update.Message)

		case h.registerHandler.CheckRegisterPhone(
			update.Message,
			(state)[h.lexicon.State.RegisterState.ID],
			h.lexicon.State.RegisterState.PhoneKey,
		):
			h.registerHandler.RegisterPhone(update.Message, user)

		default:
			continue
		}
		timeEnd := time.Since(timeNow)
		executionTimeMilliseconds := float64(timeEnd.Nanoseconds()) / 1e6
		h.mv.UpdateInfoMv(update, executionTimeMilliseconds)
	}
}
