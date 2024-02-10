package lexicon

const (
	UserStateMsg = "user_state"

	RegisterStateMsg = "register_state"

	RegisterNameKeyMsg  = "isRegisteredName"
	RegisterPhoneKeyMsg = "isRegisteredPhone"

	MenuLevelStateMsg = "menu_level_state"

	BackMenuKeyMsg    = "back_to"
	CurrentMenuKeyMsg = "current_menu"

	MenuValue    = "main_menu"
	ProfileValue = "profile"
	BalanceValue = "balance"
)

type State struct {
	UserState      UserState
	RegisterState  RegisterState
	MenuLevelState MenuLevelState
}

type UserState struct {
	ID string
}

type RegisterState struct {
	ID       string
	NameKey  string
	PhoneKey string
}

type MenuLevelState struct {
	ID             string
	CurrentMenuKey string
	BackMenuKey    string
	MenuValue      string
	ProfileValue   string
	BalanceValue   string
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
		MenuLevelState: MenuLevelState{
			ID:             MenuLevelStateMsg,
			CurrentMenuKey: CurrentMenuKeyMsg,
			BackMenuKey:    BackMenuKeyMsg,
			MenuValue:      MenuValue,
			ProfileValue:   ProfileValue,
			BalanceValue:   BalanceValue,
		},
	}
}
