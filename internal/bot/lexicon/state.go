package lexicon

const (
	UserStateMsg     = "user_state"
	RegisterStateMsg = "register_state"

	RegisterNameKeyMsg  = "isRegisteredName"
	RegisterPhoneKeyMsg = "isRegisteredPhone"
)

type State struct {
	UserState     UserState
	RegisterState RegisterState
}

type UserState struct {
	ID string
}

type RegisterState struct {
	ID       string
	NameKey  string
	PhoneKey string
}

func NewStateMsg() State {
	return State{
		UserState: UserState{
			ID: UserStateMsg,
		},
		RegisterState: RegisterState{
			ID:       RegisterStateMsg,
			NameKey:  RegisterNameKeyMsg,
			PhoneKey: RegisterPhoneKeyMsg,
		},
	}
}
