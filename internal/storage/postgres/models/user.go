package models

type User struct {
	Base
	UserName   string `gorm:"unique"`
	TelegramID int64  `gorm:"unique"`

	Name         string
	Phone        string
	IsRegistered bool

	CreatedAt int64 `gorm:"not null"`
	UpdatedAt int64
}
