package keybords

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

func (k *ReplyKeyboards) StartMenuReplyMP() tgbotapi.ReplyKeyboardMarkup {
	menuKB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.Profile),
		),
	)
	menuKB.ResizeKeyboard = true
	menuKB.OneTimeKeyboard = true

	return menuKB
}

func (k *InlineKeyboards) StartMenuInlineMP() tgbotapi.InlineKeyboardMarkup {
	menuKB := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(k.lexicon.KB.Inline.Profile, k.lexicon.KB.Inline.CallData.ProfileData),
		))
	return menuKB
}
