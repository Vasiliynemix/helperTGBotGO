package errorsMsg

import (
	"bot/internal/bot/lexicon"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type ErrorMsg struct {
	log     *logging.Logger
	bot     *tgbotapi.BotAPI
	lexicon *lexicon.Lexicon
}

func NewErrorsMsg(log *logging.Logger, bot *tgbotapi.BotAPI) *ErrorMsg {
	return &ErrorMsg{
		log:     log,
		bot:     bot,
		lexicon: lexicon.NewLexicon(),
	}
}

func (e *ErrorMsg) MsgErrorInternal(telegramID int64) {
	msgSend := tgbotapi.NewMessage(telegramID, e.lexicon.Err.OnInternalErr)

	_, err := e.bot.Send(msgSend)
	if err != nil {
		e.log.Error("Failed to send message", zap.Error(err))
	}

	e.log.Error("Сообщение об ошибке сервера отправлено пользователю", zap.Int64("telegramID", telegramID))
}
