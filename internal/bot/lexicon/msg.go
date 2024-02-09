package lexicon

import "fmt"

const (
	OnStartCmdMsg = "Привет, вот твое меню. Что хочешь сделать?"

	OnStartCmdRegisterMsg = "Привет, давай сперва зарегестрируемся! Введи свое имя."
	OnRegisterPhoneMsg    = "Отлично %s! Введи свой номер телефона."
	OnEndRegisterMsg      = "Регистрация прошла успешно!"
)

type MessageSend struct {
	OnRegisterStartCommand string
	OnStartCommand         string
	OnEndRegister          string
}

func NewSendMsg() MessageSend {
	return MessageSend{
		OnRegisterStartCommand: OnStartCmdRegisterMsg,
		OnStartCommand:         OnStartCmdMsg,
		OnEndRegister:          OnEndRegisterMsg,
	}
}

func (m *MessageSend) CreateMsgOnRegisterPhone(name string) string {
	return fmt.Sprintf(OnRegisterPhoneMsg, name)
}
