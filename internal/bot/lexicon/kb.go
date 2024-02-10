package lexicon

const (
	ProfileMsg     = "👤 Профиль"
	ProfileDataMsg = "profile"

	BalanceMsg     = "💰 Баланс"
	BalanceDataMsg = "balance"

	BackMsg           = "⬅️ Назад"
	BackToMenuDataMsg = "back_to_menu"

	AddBalanceMsg    = "➕ Пополнить"
	RemoveBalanceMsg = "➖ Вывести"
)

type KBMsg struct {
	Reply  ReplyKBMsg
	Inline InlineKBMsg
}

type ReplyKBMsg struct {
	Profile       string
	Balance       string
	AddBalance    string
	RemoveBalance string
	Back          string
}

type InlineKBMsg struct {
	Profile  string
	Balance  string
	Back     string
	CallData CallBackData
}

type CallBackData struct {
	Profile    string
	BackToMenu string
	Balance    string
}

func NewReplyKBMsg() ReplyKBMsg {
	return ReplyKBMsg{
		Profile:       ProfileMsg,
		Balance:       BalanceMsg,
		AddBalance:    AddBalanceMsg,
		RemoveBalance: RemoveBalanceMsg,
		Back:          BackMsg,
	}
}

func NewInlineKBMsg() InlineKBMsg {
	return InlineKBMsg{
		Profile:  ProfileMsg,
		Balance:  BalanceMsg,
		Back:     BackMsg,
		CallData: NewCallBackDataMsg(),
	}
}

func NewCallBackDataMsg() CallBackData {
	return CallBackData{
		Profile:    ProfileDataMsg,
		BackToMenu: BackToMenuDataMsg,
		Balance:    BalanceDataMsg,
	}
}

func NewKBMsg() KBMsg {
	return KBMsg{
		Reply:  NewReplyKBMsg(),
		Inline: NewInlineKBMsg(),
	}
}
