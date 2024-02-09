package lexicon

type Lexicon struct {
	Cmd   CmdMessage
	Msg   MessageSend
	State State
	KB    KBMsg
	Err   ErrMsg
}

func NewLexicon() *Lexicon {
	return &Lexicon{
		Cmd:   NewCmdMsg(),
		Msg:   NewSendMsg(),
		State: NewStateMsg(),
		KB:    NewKBMsg(),
		Err:   NewErrMsg(),
	}
}
