package lexicon

type KBMsg struct {
	reply  ReplyKBMsg
	inline InlineKBMsg
}

type ReplyKBMsg struct {
}

type InlineKBMsg struct {
}

func NewReplyKBMsg() KBMsg {
	return KBMsg{
		reply:  ReplyKBMsg{},
		inline: InlineKBMsg{},
	}
}
