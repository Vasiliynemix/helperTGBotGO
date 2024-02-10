package user

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (r *HandlerUser) CheckBackToMenuReply(msg *tgbotapi.Message, state interface{}, key string) bool {
	value := state.(map[string]interface{})[key].(string)
	ok := value == r.lexicon.State.MenuLevelState.MenuValue

	return msg != nil && msg.Text == r.lexicon.KB.Reply.Back && ok
}

func (r *HandlerUser) CheckBackToProfileReply(msg *tgbotapi.Message, state interface{}, key string) bool {
	value := state.(map[string]interface{})[key].(string)
	ok := value == r.lexicon.State.MenuLevelState.ProfileValue

	return msg != nil && msg.Text == r.lexicon.KB.Reply.Back && ok
}

func (r *HandlerUser) BackToMenuReply(msg *tgbotapi.Message) {
	msgSend := tgbotapi.NewMessage(msg.Chat.ID, r.lexicon.Msg.OnMenu)
	msgSend.ReplyMarkup = r.kb.Reply.StartMenuReplyMP(false)

	err := r.stateProvider.UpdateStateMap(
		msg.Chat.ID,
		r.lexicon.State.MenuLevelState.ID,
		map[string]interface{}{
			r.lexicon.State.MenuLevelState.CurrentMenuKey: r.lexicon.State.MenuLevelState.MenuValue,
			r.lexicon.State.MenuLevelState.BackMenuKey:    nil,
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

func (r *HandlerUser) BackToProfileReply(msg *tgbotapi.Message) {
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
