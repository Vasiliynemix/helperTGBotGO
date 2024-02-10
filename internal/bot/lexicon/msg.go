package lexicon

import "fmt"

const (
	OnStartCmdMsg = "Меню"

	OnStartCmdRegisterMsg = "Привет, давай сперва зарегестрируемся! Введи свое имя."
	OnRegisterPhoneMsg    = "Отлично %s! Введи свой номер телефона."
	OnEndRegisterMsg      = "Регистрация прошла успешно!"

	OnProfileMsg    = ProfileMsg
	OnMenuMsg       = OnStartCmdMsg
	OnBalanceMsg    = BalanceMsg
	OnAddBalanceMsg = AddBalanceMsg

	OnAdminPanelMsg = AdminPanelMsg
	OnUserPanelMsg  = UserPanelMsg
)

type MessageSend struct {
	OnRegisterStartCommand string
	OnStartCommand         string
	OnEndRegister          string
	OnProfile              string
	OnMenu                 string
	OnBalance              string
	OnAddBalance           string
	OnAdminPanel           string
	OnUserPanel            string
}

func NewSendMsg() MessageSend {
	return MessageSend{
		OnRegisterStartCommand: OnStartCmdRegisterMsg,
		OnStartCommand:         OnStartCmdMsg,
		OnEndRegister:          OnEndRegisterMsg,
		OnProfile:              OnProfileMsg,
		OnMenu:                 OnMenuMsg,
		OnBalance:              OnBalanceMsg,
		OnAddBalance:           OnAddBalanceMsg,
		OnAdminPanel:           OnAdminPanelMsg,
		OnUserPanel:            OnUserPanelMsg,
	}
}

func (m *MessageSend) CreateMsgOnRegisterPhone(name string) string {
	return fmt.Sprintf(OnRegisterPhoneMsg, name)
}
