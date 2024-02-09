package lexicon

const (
	OnInternalErrMsg = "Something went wrong, please try again later /start"
)

type ErrMsg struct {
	OnInternalErr string
}

func NewErrMsg() ErrMsg {
	return ErrMsg{
		OnInternalErr: OnInternalErrMsg,
	}
}
