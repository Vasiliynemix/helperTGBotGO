package lexicon

const (
	OnCmdStartMsg = "start"
)

type CmdMessage struct {
	Start string
}

func NewCmdMsg() CmdMessage {
	return CmdMessage{
		Start: OnCmdStartMsg,
	}
}
