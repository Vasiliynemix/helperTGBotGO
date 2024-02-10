package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (k *ReplyKeyboards) AdminPanelReplyMP(oneTimeKB bool) tgbotapi.ReplyKeyboardMarkup {
	menuKB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(k.lexicon.KB.Reply.User),
		),
	)
	menuKB.ResizeKeyboard = true
	menuKB.OneTimeKeyboard = oneTimeKB
	return menuKB
}
