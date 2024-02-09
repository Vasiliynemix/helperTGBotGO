package lexicon

type Lexicon struct {
	Cmd   CmdMessage
	Msg   MessageSend
	Call  CallBackData
	State State
	KB    KBMsg
	Err   ErrMsg
}

func NewLexicon() *Lexicon {
	return &Lexicon{
		Cmd:   NewCmdMsg(),
		Msg:   NewSendMsg(),
		Call:  NewCallBackDataMsg(),
		State: NewStateMsg(),
		KB:    NewReplyKBMsg(),
		Err:   NewErrMsg(),
	}
}
