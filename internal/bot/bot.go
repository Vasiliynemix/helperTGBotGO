package bot

import (
	"bot/internal/bot/handlers"
	"bot/internal/bot/middlewares"
	"bot/internal/config"
	"bot/internal/storage"
	"bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	log     *logging.Logger
	storage *storage.Storage
	bot     *tgbotapi.BotAPI
}

func NewBot(cfg *config.Config, log *logging.Logger, storage *storage.Storage) *Bot {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatal("failed to connect telegram bot", zap.Error(err))
	}

	log.Info("Connected to telegram bot", zap.String("botName", bot.Self.UserName))

	bot.Debug = cfg.Debug
	return &Bot{
		log:     log,
		storage: storage,
		bot:     bot,
	}
}

func (b *Bot) Run() {
	u := b.setupUpdateConfig()
	updates := b.bot.GetUpdatesChan(u)

	handler := handlers.NewHandlers(
		b.log,
		b.bot,
		updates,
		middlewares.InitMiddlewares(b.log, b.storage.Postgres.User, b.storage.Redis),
		b.storage,
	)

	go handler.CheckUpdates()
}

func (b *Bot) setupUpdateConfig() tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return u
}
