package account

import (
	"fmt"
	"oko/pkg/cfg"
	"oko/pkg/db"
	"oko/pkg/rndstr"
	"time"
)

type SignUpToken struct {
	ID        int        `json:"id"`
	AccountID int        `json:"account_id"`
	Token     string     `json:"token"`
	IsUsed    bool       `json:"is_used"`
	ExpireAt  time.Time  `json:"expire_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Account   Account    `json:"account" gorm:"foreignkey:AccountID;PRELOAD:true"`
}

func (SignUpToken) TableName() string {
	return "sign_up_token"
}

func GenerateSignUpToken(accID int) string {
	return fmt.Sprintf("su%d-%s", accID, rndstr.RandString(cfg.App.AuthTokenLength))
}

func GetBySignUpToken(token string) (*SignUpToken, error) {
	var sut SignUpToken
	if err := db.GetDB().
		Preload("Account").
		Where("token=? and expire_at >= now() and not is_used", token).
		First(&sut).Error; err != nil {
		return nil, err
	}
	return &sut, nil
}

func (sut SignUpToken) SendToken(email string) error {
	//TODO: Send message
	return nil
}

//nolint //TODO unused
func (sut SignUpToken) getLink() string {
	return cfg.App.FrontURL + SignUpLink + "/" + sut.Token
}
