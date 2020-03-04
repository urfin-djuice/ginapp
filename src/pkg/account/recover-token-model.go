package account

import (
	"fmt"
	"oko/pkg/cfg"
	"oko/pkg/db"
	"oko/pkg/rndstr"
	"time"
)

type RecoverToken struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Token     string    `json:"token"`
	IsUsed    bool      `json:"is_used"`
	ExpireAt  time.Time `json:"expire_at"`
	UsedAt    time.Time `json:"used_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Account   Account   `json:"account" gorm:"foreignkey:AccountID;PRELOAD:true"`
}

func (RecoverToken) TableName() string {
	return "recover_token"
}

func GenerateRecoverToken(accID int) string {
	return fmt.Sprintf("rt%d-%s", accID, rndstr.RandString(cfg.App.AuthTokenLength))
}

func GetByRecoverToken(token string) (*RecoverToken, error) {
	rt := new(RecoverToken)
	if err := db.GetDB().
		Preload("Account").
		Where("token=? and expire_at >= now() and not is_used", token).
		First(rt).Error; err != nil {
		return nil, err
	}
	return rt, nil
}

func GetRecoverToken(token string) (*RecoverToken, error) {
	rt := new(RecoverToken)
	if err := db.GetDB().
		Where("token=?", token).
		First(rt).Error; err != nil {
		return nil, err
	}
	return rt, nil
}

func (rt RecoverToken) SendChangeToken(email string) error {
	return rt.sendToken(email, RecoverTmpl, ChangePasswordSubject)
}

func (rt RecoverToken) SendRecoverToken(email string) error {
	return rt.sendToken(email, RecoverTmpl, RecoverSubject)
}

func (rt RecoverToken) sendToken(email, tmpl, subject string) error {
	// TODO: Send message
	return nil
}

//nolint //TODO unused
func (rt RecoverToken) getLink() string {
	return cfg.App.FrontURL + RecoverLink + "/" + rt.Token
}
