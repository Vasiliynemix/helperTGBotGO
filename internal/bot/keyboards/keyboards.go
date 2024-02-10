package keyboards

import (
	"bot/internal/bot/lexicon"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboards struct {
	Inline InlineKeyboards
	Reply  ReplyKeyboards
}

type InlineKeyboards struct {
	lexicon *lexicon.Lexicon
}

type ReplyKeyboards struct {
	lexicon *lexicon.Lexicon
}

func NewKeyboards(lexicon *lexicon.Lexicon) *Keyboards {
	return &Keyboards{
		Inline: InlineKeyboards{
			lexicon: lexicon,
		},
		Reply: ReplyKeyboards{
			lexicon: lexicon,
		},
	}
}

func (k *ReplyKeyboards) StartMenuReplyMP(oneTimeKB bool) tgbotapi.ReplyKeyboardMarkup {
	menuKB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.Profile),
		),
	)
	menuKB.ResizeKeyboard = true
	menuKB.OneTimeKeyboard = oneTimeKB
	return menuKB
}

func (k *ReplyKeyboards) ProfileReplyMP(oneTimeKB bool) tgbotapi.ReplyKeyboardMarkup {
	menuKB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.Balance),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.Back),
		),
	)
	menuKB.ResizeKeyboard = true
	menuKB.OneTimeKeyboard = oneTimeKB
	return menuKB
}

func (k *ReplyKeyboards) BalanceReplyMP(oneTimeKB bool) tgbotapi.ReplyKeyboardMarkup {
	menuKB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.AddBalance),
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.RemoveBalance),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.Back),
		),
	)
	menuKB.ResizeKeyboard = true
	menuKB.OneTimeKeyboard = oneTimeKB
	return menuKB
}
