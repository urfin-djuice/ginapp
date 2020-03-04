package account

import (
	"errors"
	"oko/pkg/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var AccStatusStr = map[uint]string{
	AccStatusUnconfirmed: "unconfirmed",
	AccStatusActive:      "active",
	AccStatusBaned:       "baned",
	AccStatusService:     "service",
}

type Account struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	PasswordHash string        `json:"password_hash"`
	Status       uint          `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    *time.Time    `json:"deleted_at"`
	Roles        []AccountRole `json:"roles" gorm:"foreignkey:AccountID;association_foreignkey:ID;PRELOAD:true"`
}

func (Account) TableName() string {
	return "account"
}

func (a *Account) SetPassword(password string) error {
	if password == "" {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	a.PasswordHash = string(passwordHash)
	return nil
}

func (a *Account) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(a.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func FindAccount(condition interface{}) (Account, error) {
	var model Account
	err := db.GetDB().Where(condition).First(&model).Error
	return model, err
}

func (a *Account) Update(data interface{}) error {
	err := db.GetDB().Model(a).Update(data).Error
	return err
}

func (a Account) GetStrStatus() string {
	return AccStatusStr[a.Status]
}

func (a Account) SendSignUpConfirmed() error {
	// TODO: Send message
	return nil
}
