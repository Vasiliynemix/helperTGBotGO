package lexicon

import (
	"bot/internal/bot/lexicon/commands"
	"fmt"
)

const (
	MessageOnStartCommand = "Привет, вот твое меню. Что хочешь сделать?"

	MessageOnStartRegisterCommand = "Привет, давай сперва зарегестрируемся! Введи свое имя."
	MessageOnRegisterPhone        = "Отлично %s! Введи свой номер телефона."
	MessageOnEndRegister          = "Регистрация прошла успешно!"

	MessageOnInternalError = "Something went wrong, please try again later /start"
)

type Lexicon struct {
	Cmd      CmdMessage
	Msg      MessageSend
	CallData CallBackData
	State    State
}

type State struct {
	UserState     string
	RegisterState string
}

type CmdMessage struct {
	Start string
}

type MessageSend struct {
	OnRegisterStartCommand string
	OnStartCommand         string
	OnEndRegister          string
	OnInternalError        string
}

type CallBackData struct {
}

func NewLexicon() *Lexicon {
	return &Lexicon{
		Cmd: CmdMessage{
			Start: commands.MsgCommandStart,
		},
		Msg: MessageSend{
			OnStartCommand:         MessageOnStartCommand,
			OnRegisterStartCommand: MessageOnStartRegisterCommand,
			OnEndRegister:          MessageOnEndRegister,
			OnInternalError:        MessageOnInternalError,
		},
		CallData: CallBackData{},
		State: State{
			UserState:     "user_state",
			RegisterState: "register_state",
		},
	}
}

func (m *MessageSend) CreateMsgOnRegisterPhone(name string) string {
	return fmt.Sprintf(MessageOnRegisterPhone, name)
}
