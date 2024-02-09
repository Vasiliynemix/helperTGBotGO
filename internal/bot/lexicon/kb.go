package lexicon

const (
	ProfileMsg     = "👤 Профиль"
	ProfileDataMsg = "profile"
)

type KBMsg struct {
	Reply  ReplyKBMsg
	Inline InlineKBMsg
}

type ReplyKBMsg struct {
	Profile string
}

type InlineKBMsg struct {
	Profile  string
	CallData CallBackData
}

type CallBackData struct {
	ProfileData string
}

func NewCallBackDataMsg() CallBackData {
	return CallBackData{
		ProfileData: ProfileDataMsg,
	}
}

func NewKBMsg() KBMsg {
	return KBMsg{
		Reply: ReplyKBMsg{
			Profile: ProfileMsg,
		},
		Inline: InlineKBMsg{
			Profile:  ProfileMsg,
			CallData: NewCallBackDataMsg(),
		},
	}
}
